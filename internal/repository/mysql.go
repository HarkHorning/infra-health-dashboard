package repository

import (
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Config holds the database configuration
type Config struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// DefaultConfig returns a Config with default values
// Note: Set DB_PORT=3307 if running Docker MySQL alongside local MySQL server
func DefaultConfig() Config {
	return Config{
		Host:            "localhost",
		Port:            3306,
		User:            "root",
		Password:        "devpassword",
		Database:        "infra_dashboard",
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
	}
}

// ConfigFromEnv loads database configuration from environment variables
func ConfigFromEnv() Config {
	cfg := DefaultConfig()

	if host := os.Getenv("DB_HOST"); host != "" {
		cfg.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Port = p
		}
	}
	if user := os.Getenv("DB_USER"); user != "" {
		cfg.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		cfg.Password = password
	}
	if database := os.Getenv("DB_NAME"); database != "" {
		cfg.Database = database
	}
	if maxOpen := os.Getenv("DB_MAX_OPEN_CONNS"); maxOpen != "" {
		if m, err := strconv.Atoi(maxOpen); err == nil {
			cfg.MaxOpenConns = m
		}
	}
	if maxIdle := os.Getenv("DB_MAX_IDLE_CONNS"); maxIdle != "" {
		if m, err := strconv.Atoi(maxIdle); err == nil {
			cfg.MaxIdleConns = m
		}
	}
	if lifetime := os.Getenv("DB_CONN_MAX_LIFETIME"); lifetime != "" {
		if d, err := time.ParseDuration(lifetime); err == nil {
			cfg.ConnMaxLifetime = d
		}
	}

	return cfg
}

// NewDB creates a new database connection pool with the given configuration
func NewDB(cfg Config) (*sqlx.DB, error) {
	// Debug: print connection info (remove in production)
	fmt.Printf("Connecting to MySQL: host=%s, port=%d, user=%s, database=%s\n",
		cfg.Host, cfg.Port, cfg.User, cfg.Database)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Verify the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// NewDBFromEnv creates a new database connection pool using environment variables
func NewDBFromEnv() (*sqlx.DB, error) {
	cfg := ConfigFromEnv()
	return NewDB(cfg)
}
