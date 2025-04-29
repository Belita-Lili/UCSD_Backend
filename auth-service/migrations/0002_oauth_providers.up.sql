USE auth_system;

-- Tabla para proveedores OAuth (Google, Facebook, etc.)
CREATE TABLE IF NOT EXISTS oauth_providers (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    client_id VARCHAR(255) NOT NULL,
    client_secret VARCHAR(255) NOT NULL,
    redirect_url VARCHAR(255) NOT NULL,
    enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabla para asociar usuarios con identidades OAuth
CREATE TABLE IF NOT EXISTS user_oauth_identities (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    provider_id VARCHAR(36) NOT NULL,
    provider_user_id VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    profile_data JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (provider_id) REFERENCES oauth_providers(id) ON DELETE CASCADE,
    UNIQUE KEY uk_provider_user (provider_id, provider_user_id),
    INDEX idx_user_id (user_id),
    INDEX idx_provider (provider_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insertar proveedores OAuth básicos
INSERT INTO oauth_providers (id, name, client_id, client_secret, redirect_url)
VALUES 
    (UUID(), 'google', 'your-google-client-id', 'your-google-client-secret', 'https://yourapp.com/auth/google/callback'),
    (UUID(), 'facebook', 'your-facebook-app-id', 'your-facebook-app-secret', 'https://yourapp.com/auth/facebook/callback')
ON DUPLICATE KEY UPDATE updated_at = CURRENT_TIMESTAMP;