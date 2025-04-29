USE auth_system;

-- Eliminar tablas en orden inverso a su creación
DROP TABLE IF EXISTS password_reset_tokens;
DROP TABLE IF EXISTS user_oauth_identities;
DROP TABLE IF EXISTS oauth_providers;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS login_attempts;
DROP TABLE IF EXISTS users;

-- Eliminar eventos y procedimientos
DROP EVENT IF EXISTS daily_reset_token_cleanup;
DROP PROCEDURE IF EXISTS clean_expired_reset_tokens;