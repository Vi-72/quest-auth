package cmd

import (
	"log"
	"time"

	openapihttp "github.com/Vi-72/quest-auth/api/http/auth/v1"

	"github.com/Vi-72/quest-auth/internal/adapters/in/grpc"
	adapterhttp "github.com/Vi-72/quest-auth/internal/adapters/in/http"
	bcryptadapter "github.com/Vi-72/quest-auth/internal/adapters/out/bcrypt"
	"github.com/Vi-72/quest-auth/internal/adapters/out/jwt"
	"github.com/Vi-72/quest-auth/internal/adapters/out/postgres"
	timeadapter "github.com/Vi-72/quest-auth/internal/adapters/out/time"
	"github.com/Vi-72/quest-auth/internal/core/application/usecases/commands"
	"github.com/Vi-72/quest-auth/internal/core/application/usecases/queries"
	"github.com/Vi-72/quest-auth/internal/core/ports"

	"gorm.io/gorm"
)

type CompositionRoot struct {
	configs        Config
	db             *gorm.DB
	txManager      ports.TransactionManager
	jwtService     ports.JWTService
	passwordHasher ports.PasswordHasher
	clock          ports.Clock
	closers        []Closer
}

func NewCompositionRoot(configs Config, db *gorm.DB) *CompositionRoot {
	txManager := postgres.NewTransactionManager(db)

	// Create JWT Service
	jwtService := jwt.NewService(
		configs.JWTSecretKey,
		time.Duration(configs.JWTAccessTokenDuration)*time.Minute,
		time.Duration(configs.JWTRefreshTokenDuration)*time.Hour,
	)

	// Create PasswordHasher and Clock
	passwordHasher := bcryptadapter.NewHasher()
	clock := timeadapter.NewClock()

	return &CompositionRoot{
		configs:        configs,
		db:             db,
		txManager:      txManager,
		jwtService:     jwtService,
		passwordHasher: passwordHasher,
		clock:          clock,
		closers:        []Closer{},
	}
}

// TransactionManager returns the transaction manager
func (cr *CompositionRoot) TransactionManager() ports.TransactionManager {
	return cr.txManager
}

// JWTService returns JWT service
func (cr *CompositionRoot) JWTService() ports.JWTService {
	return cr.jwtService
}

// PasswordHasher returns password hasher
func (cr *CompositionRoot) PasswordHasher() ports.PasswordHasher {
	return cr.passwordHasher
}

// Clock returns system clock
func (cr *CompositionRoot) Clock() ports.Clock {
	return cr.clock
}

// Auth Use Case Handlers

// NewRegisterUserHandler creates a handler for user registration
func (cr *CompositionRoot) NewRegisterUserHandler() *commands.RegisterUserHandler {
	return commands.NewRegisterUserHandler(
		cr.TransactionManager(),
		cr.JWTService(),
		cr.PasswordHasher(),
		cr.Clock(),
	)
}

// NewLoginUserHandler creates a handler for user login
func (cr *CompositionRoot) NewLoginUserHandler() *commands.LoginUserHandler {
	return commands.NewLoginUserHandler(
		cr.TransactionManager(),
		cr.JWTService(),
		cr.PasswordHasher(),
		cr.Clock(),
	)
}

// HTTP Handlers

// NewAPIHandler creates OpenAPI handler
func (cr *CompositionRoot) NewAPIHandler() openapihttp.StrictServerInterface {
	handlers, err := adapterhttp.NewAPIHandler(
		cr.NewRegisterUserHandler(),
		cr.NewLoginUserHandler(),
	)
	if err != nil {
		log.Fatalf("Error initializing HTTP Server: %v", err)
	}
	return handlers
}

// NewGRPCAuthHandler creates gRPC auth handler
func (cr *CompositionRoot) NewGRPCAuthHandler() *grpc.AuthHandler {
	authenticateByToken := queries.NewAuthenticateByTokenHandler(cr.JWTService())
	return grpc.NewAuthHandler(authenticateByToken)
}
