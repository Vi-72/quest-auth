package utils

import (
	"fmt"
	"log"
	"os"
	"testing"

	"quest-auth/internal/adapters/out/postgres/eventrepo"
	"quest-auth/internal/adapters/out/postgres/userrepo"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestDatabase представляет тестовую базу данных
type TestDatabase struct {
	DB   *gorm.DB
	Name string
}

// NewTestDatabase создает новую тестовую базу данных
func NewTestDatabase(t *testing.T) *TestDatabase {
	dbName := fmt.Sprintf("quest_auth_test_%s", t.Name())

	// Подключение к PostgreSQL для создания тестовой БД
	mainDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "username"),
		getEnv("DB_PASSWORD", "secret"),
		"postgres", // Подключаемся к системной БД
		getEnv("DB_SSLMODE", "disable"),
	)

	mainDB, err := gorm.Open(postgres.Open(mainDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	sqlDB, err := mainDB.DB()
	if err != nil {
		t.Fatalf("Failed to get SQL DB: %v", err)
	}

	// Создаем тестовую БД
	_, err = sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	sqlDB.Close()

	// Подключаемся к тестовой БД
	testDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "username"),
		getEnv("DB_PASSWORD", "secret"),
		dbName,
		getEnv("DB_SSLMODE", "disable"),
	)

	testDB, err := gorm.Open(postgres.Open(testDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Выполняем миграции
	err = testDB.AutoMigrate(&userrepo.UserDTO{}, &eventrepo.EventDTO{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return &TestDatabase{
		DB:   testDB,
		Name: dbName,
	}
}

// Cleanup очищает тестовую базу данных
func (td *TestDatabase) Cleanup(t *testing.T) {
	sqlDB, err := td.DB.DB()
	if err != nil {
		log.Printf("Failed to get SQL DB for cleanup: %v", err)
		return
	}
	sqlDB.Close()

	// Подключаемся к основной БД для удаления тестовой
	mainDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "username"),
		getEnv("DB_PASSWORD", "secret"),
		"postgres",
		getEnv("DB_SSLMODE", "disable"),
	)

	mainDB, err := gorm.Open(postgres.Open(mainDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Printf("Failed to connect for cleanup: %v", err)
		return
	}

	sqlDB, err = mainDB.DB()
	if err != nil {
		log.Printf("Failed to get SQL DB for cleanup: %v", err)
		return
	}
	defer sqlDB.Close()

	// Закрываем соединения с тестовой БД и удаляем её
	_, err = sqlDB.Exec(fmt.Sprintf(`
		SELECT pg_terminate_backend(pid) 
		FROM pg_stat_activity 
		WHERE datname = '%s' AND pid <> pg_backend_pid()
	`, td.Name))
	if err != nil {
		log.Printf("Failed to terminate connections: %v", err)
	}

	_, err = sqlDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", td.Name))
	if err != nil {
		log.Printf("Failed to drop test database: %v", err)
	}
}

// ClearAllTables очищает все таблицы в тестовой БД
func (td *TestDatabase) ClearAllTables() {
	td.DB.Exec("TRUNCATE TABLE users, events RESTART IDENTITY CASCADE")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
