package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"

	"quest-auth/cmd"
	grpcAdapter "quest-auth/internal/adapters/in/grpc"
)

func main() {
	_ = godotenv.Load(".env")
	configs := getConfigs()

	connectionString, err := cmd.MakeConnectionString(
		configs.DbHost,
		configs.DbPort,
		configs.DbUser,
		configs.DbPassword,
		configs.DbName,
		configs.DbSslMode)
	if err != nil {
		log.Fatal(err.Error())
	}

	cmd.CreateDbIfNotExists(configs.DbHost,
		configs.DbPort,
		configs.DbUser,
		configs.DbPassword,
		configs.DbName,
		configs.DbSslMode)
	gormDb := cmd.MustGormOpen(connectionString)
	cmd.MustAutoMigrate(gormDb)

	compositionRoot := cmd.NewCompositionRoot(
		configs,
		gormDb,
	)
	defer compositionRoot.CloseAll()

	// Создаем WaitGroup для ожидания обоих серверов
	var wg sync.WaitGroup
	wg.Add(2)

	// Запускаем HTTP сервер в горутине
	go func() {
		defer wg.Done()
		router := cmd.NewRouter(compositionRoot)
		log.Printf("HTTP server running on :%s", configs.HttpPort)
		err := http.ListenAndServe(":"+configs.HttpPort, router)
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
		HttpPort:                getEnv("HTTP_PORT"),
		GrpcPort:                getEnv("GRPC_PORT"),
		DbHost:                  getEnv("DB_HOST"),
		DbPort:                  getEnv("DB_PORT"),
		DbUser:                  getEnv("DB_USER"),
		DbPassword:              getEnv("DB_PASSWORD"),
		DbName:                  getEnv("DB_NAME"),
		DbSslMode:               getEnv("DB_SSLMODE"),
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
