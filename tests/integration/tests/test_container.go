package tests

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/Vi-72/quest-auth/cmd"
	bcryptadapter "github.com/Vi-72/quest-auth/internal/adapters/out/bcrypt"
	"github.com/Vi-72/quest-auth/internal/adapters/out/jwt"
	"github.com/Vi-72/quest-auth/internal/adapters/out/postgres"
	"github.com/Vi-72/quest-auth/internal/adapters/out/postgres/userrepo"
	timeadapter "github.com/Vi-72/quest-auth/internal/adapters/out/time"
	"github.com/Vi-72/quest-auth/internal/core/application/usecases/commands"
	"github.com/Vi-72/quest-auth/internal/core/ports"
	stor "github.com/Vi-72/quest-auth/tests/integration/core/storage"

	"gorm.io/gorm"
)

// getTestConfig возвращает конфигурацию для тестов, используя те же env переменные что и приложение
func getTestConfig() cmd.Config {
	return cmd.Config{
		HTTPPort:                getTestEnv("HTTP_PORT", "8080"),
		GrpcPort:                getTestEnv("GRPC_PORT", "9090"),
		DBHost:                  getTestEnv("DB_HOST", "localhost"),
		DBPort:                  getTestEnv("DB_PORT", "5433"),
		DBUser:                  getTestEnv("DB_USER", "postgres"),
		DBPassword:              getTestEnv("DB_PASSWORD", "password"),
		DBName:                  getTestEnv("DB_NAME", "auth_test"),
		DBSslMode:               getTestEnv("DB_SSLMODE", "disable"),
		EventGoroutineLimit:     10,
		JWTSecretKey:            getTestEnv("JWT_SECRET_KEY", "test-secret-key-for-testing-only"),
		JWTAccessTokenDuration:  1,  // 1 minute for tests
		JWTRefreshTokenDuration: 24, // 24 hours for tests
	}
}

// getTestEnv получает environment переменную или возвращает default значение
func getTestEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// TestDIContainer содержит все зависимости для интеграционных тестов
type TestDIContainer struct {
	SuiteDIContainer
	DB                 *gorm.DB
	CloseDB            func()
	TransactionManager ports.TransactionManager

	// Repositories
	UserRepository ports.UserRepository
	EventPublisher ports.EventPublisher
	JWTService     ports.JWTService

	// Use Case Handlers
	LoginUserHandler    *commands.LoginUserHandler
	RegisterUserHandler *commands.RegisterUserHandler

	// HTTP Router for API testing
	HTTPRouter http.Handler

	// Test storages
	EventStorage *stor.EventStorage
}

// NewTestDIContainer создает новый TestDIContainer для тестов
func NewTestDIContainer(suiteContainer SuiteDIContainer) TestDIContainer {
	// Используем тот же подход что и в основном приложении
	testConfig := getTestConfig()

	// Создаем базу данных если ее нет (для локальной разработки и CI)
	cmd.CreateDBIfNotExists(
		testConfig.DBHost,
		testConfig.DBPort,
		testConfig.DBUser,
		testConfig.DBPassword,
		testConfig.DBName,
		testConfig.DBSslMode,
	)

	// Формируем connection string используя ту же функцию что и в приложении
	databaseURL, err := cmd.MakeConnectionString(
		testConfig.DBHost,
		testConfig.DBPort,
		testConfig.DBUser,
		testConfig.DBPassword,
		testConfig.DBName,
		testConfig.DBSslMode,
	)
	suiteContainer.Require().NoError(err, "Failed to create database connection string")

	db, sqlDB, err := cmd.MustConnectDB(databaseURL)
	suiteContainer.Require().NoError(err, "Failed to connect to test database")

	// Создаем transaction manager
	txManager := postgres.NewTransactionManager(db)

	// Репозитории из общей базы (для запросов вне транзакции)
	userRepo := userrepo.NewRepository(db)

	// Создание EventPublisher (используем NullEventPublisher для тестов)
	eventPublisher := &ports.NullEventPublisher{}

	// Создание JWT Service для тестов
	jwtService := jwt.NewService(
		testConfig.JWTSecretKey,
		time.Duration(testConfig.JWTAccessTokenDuration)*time.Minute,
		time.Duration(testConfig.JWTRefreshTokenDuration)*time.Hour,
	)

	// Password hasher and clock
	passwordHasher := bcryptadapter.NewHasher()
	clock := timeadapter.NewClock()

	// Создание обработчиков use cases
	loginUserHandler := commands.NewLoginUserHandler(txManager, jwtService, passwordHasher, clock)
	registerUserHandler := commands.NewRegisterUserHandler(txManager, jwtService, passwordHasher, clock)

	// Create HTTP Router for API testing
	compositionRoot := cmd.NewCompositionRoot(testConfig, db)
	httpRouter := cmd.NewRouter(compositionRoot)

	// Event storage helper
	eventStorage := stor.NewEventStorage(db)

	return TestDIContainer{
		SuiteDIContainer: suiteContainer,
		DB:               db,
		CloseDB: func() {
			err := sqlDB.Close()
			if err != nil {
				return
			}
		},
		TransactionManager: txManager,

		UserRepository: userRepo,
		EventPublisher: eventPublisher,
		JWTService:     jwtService,

		LoginUserHandler:    loginUserHandler,
		RegisterUserHandler: registerUserHandler,

		HTTPRouter:   httpRouter,
		EventStorage: eventStorage,
	}
}

// TearDownTest очищает ресурсы после теста
func (c *TestDIContainer) TearDownTest() {
	if c.CloseDB != nil {
		c.CloseDB()
	}
}

// CleanupDatabase очищает тестовую базу данных
func (c *TestDIContainer) CleanupDatabase() error {
	// Очищаем таблицы в правильном порядке (учитывая внешние ключи)
	if err := c.DB.Exec("TRUNCATE TABLE events CASCADE").Error; err != nil {
		return err
	}
	if err := c.DB.Exec("TRUNCATE TABLE users CASCADE").Error; err != nil {
		return err
	}
	return nil
}

// WaitForEventProcessing actively waits until the expected number of events is stored.
// If expectedCount is 0, the method waits until the number of events stops changing.
// Waiting is canceled by a context with timeout to avoid hanging.
func (c *TestDIContainer) WaitForEventProcessing(expectedCount int64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	var lastCount int64 = -1

	for {
		select {
		case <-ctx.Done():
			c.Require().Fail("timeout waiting for events")
			return
		case <-ticker.C:
			// For auth microservice, we don't have event storage yet
			// This is a placeholder for future event processing
			count := int64(0)

			if expectedCount > 0 {
				if count >= expectedCount {
					return
				}
			} else {
				if lastCount != -1 && count == lastCount {
					return
				}
				lastCount = count
			}
		}
	}
}
