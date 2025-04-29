# Build stage
# FROM golang:1.21-alpine AS builder

WORKDIR /app

# Dependencias necesarias para compilación
RUN apk add --no-cache git gcc musl-dev

# Copiar e instalar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/auth-system ./cmd/server/

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Instalar dependencias de runtime
RUN apk add --no-cache ca-certificates tzdata

# Copiar binario y assets
COPY --from=builder /app/auth-system /app/auth-system
COPY --from=builder /app/migrations /app/migrations

# Configuración de salud y usuarios
RUN adduser -D -g '' appuser && \
    chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

# Comando de inicio
CMD ["/app/auth-system"]