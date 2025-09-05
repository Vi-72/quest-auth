package cmd

import (
	"log"
	"time"

	"quest-auth/internal/adapters/in/grpc"
	"quest-auth/internal/adapters/in/http"
	bcryptadapter "quest-auth/internal/adapters/out/bcrypt"
	"quest-auth/internal/adapters/out/jwt"
	"quest-auth/internal/adapters/out/postgres"
	"quest-auth/internal/adapters/out/postgres/eventrepo"
	timeadapter "quest-auth/internal/adapters/out/time"
	"quest-auth/internal/core/application/usecases/commands"
	"quest-auth/internal/core/application/usecases/queries"
	"quest-auth/internal/core/ports"
	"quest-auth/internal/generated/servers"

	"gorm.io/gorm"
)

type CompositionRoot struct {
	configs        Config
	db             *gorm.DB
	unitOfWork     ports.UnitOfWork
	eventPublisher ports.EventPublisher
	jwtService     ports.JWTService
	passwordHasher ports.PasswordHasher
	clock          ports.Clock
	closers        []Closer
}

func NewCompositionRoot(configs Config, db *gorm.DB) *CompositionRoot {
	// Create Unit of Work once during initialization
	unitOfWork, err := postgres.NewUnitOfWork(db)
	if err != nil {
		log.Fatalf("cannot create UnitOfWork: %v", err)
	}

	// Create EventPublisher with same Tracker as UoW for transactionality
	eventPublisher, err := eventrepo.NewRepository(unitOfWork.(ports.Tracker), configs.EventGoroutineLimit)
	if err != nil {
		log.Fatalf("cannot create EventPublisher: %v", err)
	}

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
		unitOfWork:     unitOfWork,
		eventPublisher: eventPublisher,
		jwtService:     jwtService,
		passwordHasher: passwordHasher,
		clock:          clock,
		closers:        []Closer{},
	}
}

// GetUnitOfWork returns the single UnitOfWork instance
func (cr *CompositionRoot) GetUnitOfWork() ports.UnitOfWork {
	return cr.unitOfWork
}

// EventPublisher returns EventPublisher
func (cr *CompositionRoot) EventPublisher() ports.EventPublisher {
	return cr.eventPublisher
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
	return commands.NewRegisterUserHandler(cr.GetUnitOfWork(), cr.EventPublisher(), cr.JWTService(), cr.PasswordHasher(), cr.Clock())
}

// NewLoginUserHandler creates a handler for user login
func (cr *CompositionRoot) NewLoginUserHandler() *commands.LoginUserHandler {
	return commands.NewLoginUserHandler(cr.GetUnitOfWork(), cr.EventPublisher(), cr.JWTService(), cr.PasswordHasher(), cr.Clock())
}

// HTTP Handlers

// NewAPIHandler creates OpenAPI handler
func (cr *CompositionRoot) NewAPIHandler() servers.StrictServerInterface {
	handlers, err := http.NewAPIHandler(
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
