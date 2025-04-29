package domain

import (
	"context"
	"time"
)

// UserRepository define las operaciones de persistencia para los usuarios
type UserRepository interface {
	// Operaciones básicas CRUD
	Save(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error

	// Operaciones específicas de autenticación
	UpdatePassword(ctx context.Context, userID, newPasswordHash string) error
	VerifyEmail(ctx context.Context, email string) error

	// Operaciones de recuperación de contraseña
	CreatePasswordResetToken(ctx context.Context, userID string, token string, expiresAt time.Time) error
	FindByPasswordResetToken(ctx context.Context, token string) (*User, error)
	ConsumePasswordResetToken(ctx context.Context, token string) error

	// Operaciones de auditoría
	LogLoginAttempt(ctx context.Context, userID, ipAddress, userAgent string, success bool) error
}

// OAuthRepository define las operaciones para OAuth
type OAuthRepository interface {
	FindOAuthUser(ctx context.Context, provider, providerUserID string) (*User, error)
	SaveOAuthUser(ctx context.Context, provider string, user *User, providerUserID string, profileData map[string]interface{}) error
	LinkOAuthToUser(ctx context.Context, userID, provider, providerUserID string, profileData map[string]interface{}) error
}

// TokenRepository define las operaciones para manejo de tokens
type TokenRepository interface {
	SaveRefreshToken(ctx context.Context, userID, token string, expiresAt time.Time) error
	FindRefreshToken(ctx context.Context, token string) (string, error) // Devuelve userID
	RevokeRefreshToken(ctx context.Context, token string) error
	RevokeAllRefreshTokensForUser(ctx context.Context, userID string) error
}

// SessionRepository define las operaciones para manejo de sesiones
type SessionRepository interface {
	CreateSession(ctx context.Context, session *Session) error
	FindSessionByID(ctx context.Context, sessionID string) (*Session, error)
	DeleteSession(ctx context.Context, sessionID string) error
	DeleteSessionsForUser(ctx context.Context, userID string) error
}

// Session representa una sesión de usuario
type Session struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

// AuditLogRepository define las operaciones para registro de actividades
type AuditLogRepository interface {
	LogEvent(ctx context.Context, event *AuditEvent) error
	GetEventsForUser(ctx context.Context, userID string, limit, offset int) ([]*AuditEvent, error)
}

// AuditEvent representa un evento de auditoría
type AuditEvent struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"`
	Entity    string    `json:"entity"`
	EntityID  string    `json:"entity_id"`
	Metadata  string    `json:"metadata"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

// UnitOfWork define el patrón Unit of Work para transacciones
type UnitOfWork interface {
	Begin(ctx context.Context) (UnitOfWork, error)
	Commit() error
	Rollback() error

	Users() UserRepository
	OAuth() OAuthRepository
	Tokens() TokenRepository
	Sessions() SessionRepository
	AuditLogs() AuditLogRepository
}

// RepositoryFactory es una fábrica abstracta para crear repositorios
type RepositoryFactory interface {
	NewUserRepository(ctx context.Context) UserRepository
	NewOAuthRepository(ctx context.Context) OAuthRepository
	NewTokenRepository(ctx context.Context) TokenRepository
	NewSessionRepository(ctx context.Context) SessionRepository
	NewAuditLogRepository(ctx context.Context) AuditLogRepository
	NewUnitOfWork(ctx context.Context) (UnitOfWork, error)
}
