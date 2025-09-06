package database

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
	"github.com/walterfan/lazy-rabbit-reminder/internal/models"
	"github.com/walterfan/lazy-rabbit-reminder/internal/service"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Type     string `mapstructure:"type"` // sqlite, postgres, mysql
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	SSLMode  string `mapstructure:"ssl_mode"`
	Charset  string `mapstructure:"charset"`
	FilePath string `mapstructure:"file_path"` // for SQLite
}

// InitDB initializes database connection based on configuration
func InitDB() error {
	config := loadDatabaseConfig()

	var err error
	DB, err = connectDatabase(config)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate database schema
	if err := DB.AutoMigrate(models.GetAllModels()...); err != nil {
		return fmt.Errorf("auto-migration failed: %w", err)
	}

	// Initialize data
	if err := InitData(); err != nil {
		return fmt.Errorf("failed to initialize data: %w", err)
	}

	log.Printf("Successfully connected to %s database", config.Type)
	return nil
}

// loadDatabaseConfig loads database configuration from environment variables and config
func loadDatabaseConfig() *DatabaseConfig {
	config := &DatabaseConfig{
		Type:     getEnvOrDefault("DB_TYPE", "sqlite"),
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     getEnvIntOrDefault("DB_PORT", 5432),
		Username: getEnvOrDefault("DB_USERNAME", ""),
		Password: getEnvOrDefault("DB_PASSWORD", ""),
		Database: getEnvOrDefault("DB_NAME", "lazy-rabbit-reminder"),
		SSLMode:  getEnvOrDefault("DB_SSL_MODE", "disable"),
		Charset:  getEnvOrDefault("DB_CHARSET", "utf8mb4"),
		FilePath: getEnvOrDefault("DB_FILE_PATH", "lazy-rabbit-reminder.db"),
	}

	// Override with viper config if available
	if viper.IsSet("database.type") {
		config.Type = viper.GetString("database.type")
	}
	if viper.IsSet("database.host") {
		config.Host = viper.GetString("database.host")
	}
	if viper.IsSet("database.port") {
		config.Port = viper.GetInt("database.port")
	}
	if viper.IsSet("database.username") {
		config.Username = viper.GetString("database.username")
	}
	if viper.IsSet("database.password") {
		config.Password = viper.GetString("database.password")
	}
	if viper.IsSet("database.database") {
		config.Database = viper.GetString("database.database")
	}
	if viper.IsSet("database.ssl_mode") {
		config.SSLMode = viper.GetString("database.ssl_mode")
	}
	if viper.IsSet("database.charset") {
		config.Charset = viper.GetString("database.charset")
	}
	if viper.IsSet("database.file_path") {
		config.FilePath = viper.GetString("database.file_path")
	}

	return config
}

// connectDatabase establishes database connection based on type
func connectDatabase(config *DatabaseConfig) (*gorm.DB, error) {
	switch config.Type {
	case "sqlite":
		return connectSQLite(config)
	case "postgres", "postgresql":
		return connectPostgreSQL(config)
	case "mysql":
		return connectMySQL(config)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", config.Type)
	}
}

// connectSQLite establishes SQLite connection
func connectSQLite(config *DatabaseConfig) (*gorm.DB, error) {
	dsn := config.FilePath
	if dsn == "" {
		dsn = "lazy-rabbit-reminder.db"
	}

	log.Printf("Connecting to SQLite database: %s", dsn)
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
}

// connectPostgreSQL establishes PostgreSQL connection
func connectPostgreSQL(config *DatabaseConfig) (*gorm.DB, error) {
	if config.Port == 0 {
		config.Port = 5432
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.Database, config.SSLMode)

	log.Printf("Connecting to PostgreSQL database: %s:%d/%s", config.Host, config.Port, config.Database)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// connectMySQL establishes MySQL connection
func connectMySQL(config *DatabaseConfig) (*gorm.DB, error) {
	if config.Port == 0 {
		config.Port = 3306
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.Username, config.Password, config.Host, config.Port, config.Database, config.Charset)

	log.Printf("Connecting to MySQL database: %s:%d/%s", config.Host, config.Port, config.Database)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

// InitData initializes database with default data (legacy function)
func InitData() error {
	var count int64

	// Initialize prompts
	DB.Model(&models.Prompt{}).Count(&count)
	if count == 0 {
		// Load prompts from config
		var prompts []models.Prompt
		if err := viper.UnmarshalKey("prompts", &prompts); err != nil {
			return fmt.Errorf("unable to decode prompts into struct: %w", err)
		}

		if len(prompts) > 0 {
			result := DB.Create(&prompts)
			if result.Error != nil {
				return fmt.Errorf("failed to insert initial prompt data: %w", result.Error)
			}
			log.Printf("Initialized database with %d prompts", len(prompts))
		}
	}

	// Call the comprehensive data initialization from service package
	return service.InitCompleteData(DB)
}

// Helper functions
func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvOrDefault(key, defaultValue string) string {
	return GetEnvOrDefault(key, defaultValue)
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// CloseDB closes the database connection
func CloseDB() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
