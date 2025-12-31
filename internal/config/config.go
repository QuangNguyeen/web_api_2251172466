package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config holds all configuration
type Config struct {
	Port               string
	DBHost             string
	DBPort             string
	DBName             string
	DBUser             string
	DBPassword         string
	JWTSecret          string
	JWTExpiration      time.Duration
	CORSAllowedOrigins string
}

// Load loads configuration from environment variables
func Load() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Parse JWT expiration
	jwtExp, err := time.ParseDuration(getEnv("JWT_EXPIRATION", "24h"))
	if err != nil {
		jwtExp = 24 * time.Hour
	}

	return &Config{
		Port:               getEnv("PORT", "8080"),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBName:             getEnv("DB_NAME", "library_db"),
		DBUser:             getEnv("DB_USER", ""),
		DBPassword:         getEnv("DB_PASSWORD", ""),
		JWTSecret:          getEnv("JWT_SECRET", "default-secret-key-change-in-production"),
		JWTExpiration:      jwtExp,
		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "*"),
	}
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// InitDB initializes database connection
func InitDB(cfg *Config) *gorm.DB {
	// Build DSN dynamically to support empty user/password
	dsn := fmt.Sprintf("host=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Ho_Chi_Minh",
		cfg.DBHost, cfg.DBName, cfg.DBPort)

	// Add user if provided
	if cfg.DBUser != "" {
		dsn = fmt.Sprintf("user=%s %s", cfg.DBUser, dsn)
	}

	// Add password if provided
	if cfg.DBPassword != "" {
		dsn = fmt.Sprintf("%s password=%s", dsn, cfg.DBPassword)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("✅ Database connected successfully!")
	return db
}
