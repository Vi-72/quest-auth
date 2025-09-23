package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"

	"quest-auth/cmd"
	grpcAdapter "quest-auth/internal/adapters/in/grpc"
)

func main() {
	_ = godotenv.Load(".env")
	configs := getConfigs()

	connectionString, err := cmd.MakeConnectionString(
		configs.DBHost,
		configs.DBPort,
		configs.DBUser,
		configs.DBPassword,
		configs.DBName,
		configs.DBSslMode)
	if err != nil {
		log.Fatal(err.Error())
	}

	cmd.CreateDBIfNotExists(configs.DBHost,
		configs.DBPort,
		configs.DBUser,
		configs.DBPassword,
		configs.DBName,
		configs.DBSslMode)
	gormDB := cmd.MustGormOpen(connectionString)
	cmd.MustAutoMigrate(gormDB)

	compositionRoot := cmd.NewCompositionRoot(
		configs,
		gormDB,
	)
	defer compositionRoot.CloseAll()

	// Создаем WaitGroup для ожидания обоих серверов
	const numServers = 2 // HTTP и gRPC серверы
	var wg sync.WaitGroup
	wg.Add(numServers)

	// Запускаем HTTP сервер в горутине
	go func() {
		defer wg.Done()
		router := cmd.NewRouter(compositionRoot)
		log.Printf("HTTP server running on :%s", configs.HTTPPort)

		// Создаем HTTP сервер с таймаутами для безопасности
		const (
			readTimeout  = 15 * time.Second
			writeTimeout = 15 * time.Second
			idleTimeout  = 60 * time.Second
		)
		server := &http.Server{
			Addr:         ":" + configs.HTTPPort,
			Handler:      router,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			IdleTimeout:  idleTimeout,
		}

		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("failed to start HTTP server: %v", err)
		}
	}()

	// Запускаем gRPC сервер в горутине
	go func() {
		defer wg.Done()
		authHandler := compositionRoot.NewGRPCAuthHandler()
		err := grpcAdapter.StartServer(configs.GrpcPort, authHandler)
		if err != nil {
			log.Fatalf("failed to start gRPC server: %v", err)
		}
	}()

	// Ожидаем завершения обоих серверов
	wg.Wait()
}

func getConfigs() cmd.Config {
	return cmd.Config{
		HTTPPort:                getEnv("HTTP_PORT"),
		GrpcPort:                getEnv("GRPC_PORT"),
		DBHost:                  getEnv("DB_HOST"),
		DBPort:                  getEnv("DB_PORT"),
		DBUser:                  getEnv("DB_USER"),
		DBPassword:              getEnv("DB_PASSWORD"),
		DBName:                  getEnv("DB_NAME"),
		DBSslMode:               getEnv("DB_SSLMODE"),
		EventGoroutineLimit:     getEnvInt("EVENT_GOROUTINE_LIMIT"),
		JWTSecretKey:            getEnv("JWT_SECRET_KEY"),
		JWTAccessTokenDuration:  getEnvInt("JWT_ACCESS_TOKEN_DURATION"),
		JWTRefreshTokenDuration: getEnvInt("JWT_REFRESH_TOKEN_DURATION"),
	}
}

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Missing env var: %s", key)
	}
	return val
}

func getEnvInt(key string) int {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Missing env var: %s", key)
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("Invalid integer value for env var %s: %s", key, val)
	}
	return intVal
}
