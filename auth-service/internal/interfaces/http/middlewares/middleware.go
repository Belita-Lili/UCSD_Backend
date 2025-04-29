package middlewares

import (
	"net/http"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/LiliBeta/auth-service/pkg/utils"

	"github.com/LiliBeta/auth-service/internal/domain"
	"github.com/LiliBeta/auth-service/pkg/logger"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware maneja la autenticación JWT
type AuthMiddleware struct {
	authService domain.AuthService
}

func NewAuthMiddleware(authService domain.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func (m *AuthMiddleware) ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// Formato esperado: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		token := parts[1]
		valid, err := m.authService.VerifyToken(c.Request.Context(), token)
		if err != nil || !valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Extraer claims y establecer en contexto
		claims, err := m.authService.ExtractClaims(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to extract token claims"})
			return
		}

		c.Set("userID", claims["sub"])
		c.Set("roles", claims["roles"])
		c.Next()
	}
}

// CORSMiddleware configura los headers CORS
func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowed := false

		for _, o := range allowedOrigins {
			if o == "*" || o == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers",
				"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods",
				"POST, OPTIONS, GET, PUT, PATCH, DELETE")

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
		}

		c.Next()
	}
}

// SecurityHeadersMiddleware añade headers de seguridad
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Headers de seguridad básicos
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Writer.Header().Set("Content-Security-Policy",
			"default-src 'self'; script-src 'self' 'unsafe-inline' cdn.example.com; style-src 'self' 'unsafe-inline'")

		c.Next()
	}
}

// RateLimiterMiddleware limita las peticiones por IP
func RateLimiterMiddleware(requestsPerMinute int) gin.HandlerFunc {
	type client struct {
		limiter  *utils.RateLimiter
		lastSeen time.Time
	}

	var (
		clients = make(map[string]*client)
		mu      sync.Mutex
	)

	// Limpiar clientes antiguos cada minuto
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := utils.GetClientIP(c)

		mu.Lock()
		cli, exists := clients[ip]
		if !exists {
			cli = &client{
				limiter: utils.NewRateLimiter(requestsPerMinute, time.Minute),
			}
			clients[ip] = cli
		}
		cli.lastSeen = time.Now()

		if !cli.limiter.Allow() {
			mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
			return
		}
		mu.Unlock()

		c.Next()
	}
}

// CSRFMiddleware protege contra ataques CSRF
func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" || c.Request.Method == "DELETE" {
			// Verificar token CSRF para métodos que modifican datos
			csrfToken := c.GetHeader("X-CSRF-Token")
			if csrfToken == "" {
				csrfToken = c.Query("csrf_token")
			}

			sessionToken, err := c.Cookie("csrf_token")
			if err != nil || csrfToken != sessionToken {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "Invalid CSRF token",
				})
				return
			}
		}

		// Para GET/HEAD, establecer cookie CSRF si no existe
		if _, err := c.Cookie("csrf_token"); err != nil {
			token := utils.GenerateRandomString(32)
			c.SetCookie("csrf_token", token, 3600, "/", "", false, true)
		}

		c.Next()
	}
}

// RequestLoggerMiddleware registra información de las peticiones
func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		fields := []logger.Field{
			logger.Int("status", c.Writer.Status()),
			logger.String("method", c.Request.Method),
			logger.String("path", path),
			logger.String("query", query),
			logger.String("ip", c.ClientIP()),
			logger.String("user-agent", c.Request.UserAgent()),
			logger.Duration("latency", latency),
		}

		if c.Writer.Status() >= http.StatusInternalServerError {
			logger.Error("Server error", fields...)
		} else {
			logger.Info("Request handled", fields...)
		}
	}
}

// RoleBasedAuthMiddleware verifica roles del usuario
func RoleBasedAuthMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoles, exists := c.Get("roles")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No roles information found"})
			return
		}

		roles, ok := userRoles.([]string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid roles format"})
			return
		}

		hasRole := false
		for _, requiredRole := range requiredRoles {
			for _, userRole := range roles {
				if userRole == requiredRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			return
		}

		c.Next()
	}
}

// RecoveryMiddleware maneja panics y errores no capturados
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Registrar el error
				logger.Error("Recovered from panic",
					logger.Any("error", err),
					logger.String("stack", string(debug.Stack())),
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
			}
		}()

		c.Next()
	}
}
