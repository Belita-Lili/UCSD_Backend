USE auth_system;

-- Tabla para tokens de recuperación de contraseña
CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    token VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_token (token),
    INDEX idx_expires_at (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Procedimiento almacenado para limpiar tokens expirados
DELIMITER //
CREATE PROCEDURE clean_expired_reset_tokens()
BEGIN
    DELETE FROM password_reset_tokens WHERE expires_at < NOW();
END //
DELIMITER ;

-- Evento para ejecutar limpieza diaria
CREATE EVENT IF NOT EXISTS daily_reset_token_cleanup
ON SCHEDULE EVERY 1 DAY
STARTS CURRENT_TIMESTAMP
DO
    CALL clean_expired_reset_tokens();