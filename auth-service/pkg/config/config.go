package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env      string `mapstructure:"ENV"`
	DB       DatabaseConfig
	HTTP     HTTPConfig
	Keycloak KeycloakConfig
	OAuth    OAuthConfig
	Email    EmailConfig
}

type DatabaseConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
	SSLMode  string `mapstructure:"DB_SSL_MODE"`
}

type HTTPConfig struct {
	Port               string        `mapstructure:"HTTP_PORT"`
	ReadTimeout        time.Duration `mapstructure:"HTTP_READ_TIMEOUT"`
	WriteTimeout       time.Duration `mapstructure:"HTTP_WRITE_TIMEOUT"`
	IdleTimeout        time.Duration `mapstructure:"HTTP_IDLE_TIMEOUT"`
	ShutdownTimeout    time.Duration `mapstructure:"HTTP_SHUTDOWN_TIMEOUT"`
	CORSAllowedOrigins []string      `mapstructure:"CORS_ALLOWED_ORIGINS"`
}

type KeycloakConfig struct {
	URL          string `mapstructure:"KEYCLOAK_URL"`
	Realm        string `mapstructure:"KEYCLOAK_REALM"`
	ClientID     string `mapstructure:"KEYCLOAK_CLIENT_ID"`
	ClientSecret string `mapstructure:"KEYCLOAK_CLIENT_SECRET"`
}

type OAuthConfig struct {
	Google   OAuthProviderConfig
	Facebook OAuthProviderConfig
}

type OAuthProviderConfig struct {
	ClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	ClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	RedirectURL  string `mapstructure:"GOOGLE_REDIRECT_URL"`
}

type EmailConfig struct {
	From     string `mapstructure:"EMAIL_FROM"`
	SMTPHost string `mapstructure:"SMTP_HOST"`
	SMTPPort int    `mapstructure:"SMTP_PORT"`
	Username string `mapstructure:"SMTP_USERNAME"`
	Password string `mapstructure:"SMTP_PASSWORD"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("HTTP_PORT", "8080")
	viper.SetDefault("ENV", "development")

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return
		}
		log.Println("No config file found, using environment variables")
	}

	err = viper.Unmarshal(&config)
	return
}
