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

type DB struct {
	*gorm.DB
}

func Initialize() (*DB, error) {
	var db *gorm.DB
	var err error

	useRamsql := getEnv("USE_RAMSQL", "true")

	if useRamsql == "true" {
		log.Println("Ramsql is only available for testing. Please use PostgreSQL for production.")
		return nil, fmt.Errorf("ramsql is only available for testing. Set USE_RAMSQL=false to use PostgreSQL")
	} else {
		host := getEnv("DB_HOST", "localhost")
		port := getEnv("DB_PORT", "5432")
		user := getEnv("DB_USER", "postgres")
		password := getEnv("DB_PASSWORD", "password")
		dbname := getEnv("DB_NAME", "pantryos_db")
		sslmode := getEnv("DB_SSLMODE", "disable")

		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			host, port, user, password, dbname, sslmode)

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: getGormLogger(),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
		}

		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database connection initialized successfully")
	return &DB{DB: db}, nil
}

func getGormLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
}

func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Organization{},
		&models.Account{},
		&models.User{},
		&models.Category{},
		&models.InventoryItem{},
		&models.MenuItem{},
		&models.RecipeIngredient{},
		&models.InventorySnapshot{},
		&models.Sale{},
		&models.SaleItem{},
		&models.Order{},
		&models.OrderItem{},
		&models.Delivery{},
		&models.AccountInvitation{},
		&models.EmailSchedule{},
	)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
