package tests

import (
	"context"
	"net/http"
	"os"
	"time"

	"quest-auth/cmd"
	bcryptadapter "quest-auth/internal/adapters/out/bcrypt"
	"quest-auth/internal/adapters/out/jwt"
	"quest-auth/internal/adapters/out/postgres"
	timeadapter "quest-auth/internal/adapters/out/time"
	"quest-auth/internal/core/application/usecases/commands"
	"quest-auth/internal/core/ports"
	stor "quest-auth/tests/integration/core/storage"

	"gorm.io/gorm"
)

// getTestConfig возвращает конфигурацию для тестов, используя те же env переменные что и приложение
func getTestConfig() cmd.Config {
	return cmd.Config{
		HttpPort:                getTestEnv("HTTP_PORT", "8080"),
		GrpcPort:                getTestEnv("GRPC_PORT", "9090"),
		DbHost:                  getTestEnv("DB_HOST", "localhost"),
		DbPort:                  getTestEnv("DB_PORT", "5433"),
		DbUser:                  getTestEnv("DB_USER", "postgres"),
		DbPassword:              getTestEnv("DB_PASSWORD", "password"),
		DbName:                  getTestEnv("DB_NAME", "auth_test"),
		DbSslMode:               getTestEnv("DB_SSLMODE", "disable"),
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
	DB         *gorm.DB
	CloseDB    func()
	UnitOfWork ports.UnitOfWork

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
	cmd.CreateDbIfNotExists(
		testConfig.DbHost,
		testConfig.DbPort,
		testConfig.DbUser,
		testConfig.DbPassword,
		testConfig.DbName,
		testConfig.DbSslMode,
	)

	// Формируем connection string используя ту же функцию что и в приложении
	databaseURL, err := cmd.MakeConnectionString(
		testConfig.DbHost,
		testConfig.DbPort,
		testConfig.DbUser,
		testConfig.DbPassword,
		testConfig.DbName,
		testConfig.DbSslMode,
	)
	suiteContainer.Require().NoError(err, "Failed to create database connection string")

	db, sqlDB, err := cmd.MustConnectDB(databaseURL)
	suiteContainer.Require().NoError(err, "Failed to connect to test database")

	// Создание Unit of Work
	unitOfWork, err := postgres.NewUnitOfWork(db)
	suiteContainer.Require().NoError(err, "Failed to create unit of work")

	// Получаем репозитории из UnitOfWork
	userRepo := unitOfWork.UserRepository()

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
	loginUserHandler := commands.NewLoginUserHandler(unitOfWork, eventPublisher, jwtService, passwordHasher, clock)
	registerUserHandler := commands.NewRegisterUserHandler(unitOfWork, eventPublisher, jwtService, passwordHasher, clock)

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
		UnitOfWork: unitOfWork,

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
