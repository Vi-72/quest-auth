package grpc

import (
	"log"
	"net"
	authv1 "quest-auth/api/grpc/proto/auth/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// StartServer запускает gRPC сервер
func StartServer(port string, authHandler *AuthHandler) error {
	// Создаем listener
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// Создаем gRPC сервер
	grpcServer := grpc.NewServer()

	// Регистрируем сервис аутентификации
	authv1.RegisterAuthServiceServer(grpcServer, authHandler)

	// Включаем reflection для удобства отладки
	reflection.Register(grpcServer)

	log.Printf("gRPC server running on :%s", port)

	// Запускаем сервер
	return grpcServer.Serve(lis)
}
