package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	GitHub   GitHubConfig
	Crypto   CryptoConfig
	Redis    RedisConfig
	OpenAI   OpenAIConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type GitHubConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type CryptoConfig struct {
	EncryptionKey string
	JWTSecret     string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	Enabled  bool
}

type OpenAIConfig struct {
	APIKey  string
	Enabled bool
}

func Load() (*Config, error) {
	// Load .env file if it exists (ignore error in production)
	_ = godotenv.Load()

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "resume_builder"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		GitHub: GitHubConfig{
			ClientID:     mustGetEnv("GITHUB_CLIENT_ID"),
			ClientSecret: mustGetEnv("GITHUB_CLIENT_SECRET"),
			RedirectURL:  getEnv("GITHUB_REDIRECT_URL", "http://localhost:8080/auth/callback"),
		},
		Crypto: CryptoConfig{
			EncryptionKey: mustGetEnv("ENCRYPTION_KEY"),
			JWTSecret:     mustGetEnv("JWT_SECRET"),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       0,
			Enabled:  getEnv("REDIS_ENABLED", "false") == "true",
		},
		OpenAI: OpenAIConfig{
			APIKey:  getEnv("OPENAI_API_KEY", ""),
			Enabled: getEnv("OPENAI_ENABLED", "false") == "true",
		},
	}

	return cfg, nil
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	return value
}
