package tests

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mwangaben/validation"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db    *gorm.DB
	sqlDB *sql.DB
)

// getEnv gets environment variable or returns default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// TestMain sets up the database connection and runs tests
func TestMain(m *testing.M) {
	// Get database configuration from environment with defaults
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "root")
	dbName := getEnv("DB_NAME", "gotest")

	// Print connection info for debugging
	fmt.Printf("Connecting to database: %s@%s:%s/%s\n", dbUser, dbHost, dbPort, dbName)

	// Build DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err = db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB: %v", err)
	}

	// Ping database to verify connection
	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	fmt.Println("✅ Successfully connected to database")

	// Run migrations
	err = setupDatabase()
	if err != nil {
		log.Fatalf("Failed to setup database: %v", err)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	teardownDatabase()

	os.Exit(code)
}

// User model for testing - FIXED: Use gorm tags to specify VARCHAR length
type User struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"type:varchar(100);uniqueIndex:idx_users_name"`
	Email string `gorm:"type:varchar(100);uniqueIndex:idx_users_email"`
	Age   int
}

// Product model for testing - FIXED: Use gorm tags to specify VARCHAR length
type Product struct {
	ID    uint   `gorm:"primaryKey"`
	Code  string `gorm:"type:varchar(50);uniqueIndex:idx_products_code"`
	Name  string `gorm:"type:varchar(100)"`
	Price float64
}

func setupDatabase() error {
	// Drop existing tables first (clean slate)
	db.Exec("DROP TABLE IF EXISTS products")
	db.Exec("DROP TABLE IF EXISTS users")

	// Auto migrate tables with proper types
	err := db.AutoMigrate(&User{}, &Product{})
	if err != nil {
		return fmt.Errorf("failed to migrate: %v", err)
	}

	// Insert test data
	users := []User{
		{Name: "John Doe", Email: "john@example.com", Age: 30},
		{Name: "Jane Smith", Email: "jane@example.com", Age: 25},
		{Name: "Bob Johnson", Email: "bob@example.com", Age: 35},
	}

	for _, user := range users {
		result := db.Create(&user)
		if result.Error != nil {
			return fmt.Errorf("failed to insert test data: %v", result.Error)
		}
	}

	products := []Product{
		{Code: "P001", Name: "Laptop", Price: 999.99},
		{Code: "P002", Name: "Mouse", Price: 29.99},
	}

	for _, product := range products {
		result := db.Create(&product)
		if result.Error != nil {
			return fmt.Errorf("failed to insert test data: %v", result.Error)
		}
	}

	fmt.Println("✅ Test data inserted successfully")
	return nil
}

func teardownDatabase() {
	// Clean up tables
	db.Exec("DROP TABLE IF EXISTS users")
	db.Exec("DROP TABLE IF EXISTS products")

	if sqlDB != nil {
		sqlDB.Close()
	}
	fmt.Println("✅ Database cleaned up")
}
