package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mnadev/pantryos/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB wraps the GORM database connection
// This provides a clean interface for database operations and allows
// for easier testing and mocking
type DB struct {
	*gorm.DB
}

// Initialize creates a new database connection and runs migrations
// This function handles both development and production database setup
// For production, it connects to PostgreSQL. For testing, use SetupTestDB instead.
func Initialize() (*DB, error) {
	var db *gorm.DB
	var err error

	// Check if we should use ramsql for development
	// Note: Ramsql is only available for testing, not for production use
	useRamsql := getEnv("USE_RAMSQL", "true")

	if useRamsql == "true" {
		// Use ramsql in-memory database for development
		// Note: Driver registration is handled in tests
		// For production, use PostgreSQL instead
		log.Println("Ramsql is only available for testing. Please use PostgreSQL for production.")
		return nil, fmt.Errorf("ramsql is only available for testing. Set USE_RAMSQL=false to use PostgreSQL")
	} else {
		// Use PostgreSQL for production
		// Get database configuration from environment variables
		host := getEnv("DB_HOST", "localhost")
		port := getEnv("DB_PORT", "5432")
		user := getEnv("DB_USER", "postgres")
		password := getEnv("DB_PASSWORD", "password")
		dbname := getEnv("DB_NAME", "pantryos_db")
		sslmode := getEnv("DB_SSLMODE", "disable")

		// Create DSN (Data Source Name) for PostgreSQL connection
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			host, port, user, password, dbname, sslmode)

		// Open database connection with configured logger
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: getGormLogger(),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}

		// Get underlying sql.DB for connection pool configuration
		// This allows us to configure connection pooling for better performance
		sqlDB, err := db.DB()
		if err != nil {
			return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
		}

		// Configure connection pool settings
		// These settings help manage database connections efficiently
		sqlDB.SetMaxIdleConns(10)           // Maximum number of idle connections
		sqlDB.SetMaxOpenConns(100)          // Maximum number of open connections
		sqlDB.SetConnMaxLifetime(time.Hour) // Maximum lifetime of connections
	}

	// Run database migrations to ensure all tables exist
	// This creates or updates the database schema as needed
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database connection initialized successfully")
	return &DB{DB: db}, nil
}

// getGormLogger returns a configured GORM logger
// This provides structured logging for database operations
func getGormLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second, // Log queries that take longer than 1 second
			LogLevel:                  logger.Info, // Log all SQL queries
			IgnoreRecordNotFoundError: true,        // Don't log "record not found" errors
			Colorful:                  true,        // Use colors in console output
		},
	)
}

// runMigrations creates all database tables
// This ensures the database schema matches the current model definitions
// GORM's AutoMigrate will create tables, add columns, and update indexes as needed
func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Organization{},      // Multi-tenant organizations
		&models.Account{},           // Business locations/accounts within organizations
		&models.User{},              // Users with roles and permissions
		&models.Category{},          // Categories for organizing inventory and menu items
		&models.InventoryItem{},     // Inventory items with stock levels
		&models.MenuItem{},          // Menu items with categories and pricing
		&models.RecipeIngredient{},  // Recipe ingredients linking menu items to inventory
		&models.InventorySnapshot{}, // Historical inventory snapshots
		&models.Delivery{},          // Delivery records for inventory replenishment
		&models.AccountInvitation{}, // User invitations for account access
	)
}

// getEnv gets an environment variable with a fallback default value
// This provides a clean way to handle configuration with sensible defaults
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// Close closes the database connection
// This should be called when shutting down the application
func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
