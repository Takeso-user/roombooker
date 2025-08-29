package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	Graph    GraphConfig
	App      AppConfig
}

type ServerConfig struct {
	Port int
}

type DatabaseConfig struct {
	Driver string
	DSN    string
}

type AuthConfig struct {
	JWTSecret        string
	OIDCIssuer       string
	OIDCClientID     string
	OIDCClientSecret string
	OIDCRedirectURL  string
}

type GraphConfig struct {
	ClientID     string
	ClientSecret string
	TenantID     string
}

type AppConfig struct {
	BaseURL  string
	OfficeTZ string
}

func Load() (*Config, error) {
	godotenv.Load(".env")

	viper.SetDefault("SERVER_PORT", 8080)
	viper.SetDefault("DATABASE_DRIVER", "sqlite3")
	viper.SetDefault("DATABASE_DSN", "file:roombooker.db?cache=shared&_fk=1")
	viper.SetDefault("JWT_SECRET", "your-secret-key")
	viper.SetDefault("APP_BASE_URL", "http://localhost:8080")
	viper.SetDefault("OFFICE_TZ", "America/New_York")

	viper.AutomaticEnv()

	cfg := &Config{
		Server: ServerConfig{
			Port: viper.GetInt("SERVER_PORT"),
		},
		Database: DatabaseConfig{
			Driver: viper.GetString("DATABASE_DRIVER"),
			DSN:    viper.GetString("DATABASE_DSN"),
		},
		Auth: AuthConfig{
			JWTSecret:        viper.GetString("JWT_SECRET"),
			OIDCIssuer:       viper.GetString("OIDC_ISSUER"),
			OIDCClientID:     viper.GetString("OIDC_CLIENT_ID"),
			OIDCClientSecret: viper.GetString("OIDC_CLIENT_SECRET"),
			OIDCRedirectURL:  viper.GetString("OIDC_REDIRECT_URL"),
		},
		Graph: GraphConfig{
			ClientID:     viper.GetString("GRAPH_CLIENT_ID"),
			ClientSecret: viper.GetString("GRAPH_CLIENT_SECRET"),
			TenantID:     viper.GetString("GRAPH_TENANT_ID"),
		},
		App: AppConfig{
			BaseURL:  viper.GetString("APP_BASE_URL"),
			OfficeTZ: viper.GetString("OFFICE_TZ"),
		},
	}

	return cfg, nil
}
