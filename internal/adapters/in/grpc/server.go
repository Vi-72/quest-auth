package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	authpb "quest-auth/api/grpc/proto"
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
	authpb.RegisterAuthServiceServer(grpcServer, authHandler)

	// Включаем reflection для удобства отладки
	reflection.Register(grpcServer)

	log.Printf("gRPC server running on :%s", port)

	// Запускаем сервер
	return grpcServer.Serve(lis)
}
