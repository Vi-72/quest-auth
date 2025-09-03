package cmd

import (
	"log"
	"time"

	"quest-auth/internal/adapters/in/http"
	"quest-auth/internal/adapters/out/jwt"
	"quest-auth/internal/adapters/out/postgres"
	"quest-auth/internal/adapters/out/postgres/eventrepo"
	"quest-auth/internal/core/application/usecases/auth"
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

	return &CompositionRoot{
		configs:        configs,
		db:             db,
		unitOfWork:     unitOfWork,
		eventPublisher: eventPublisher,
		jwtService:     jwtService,
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

// Auth Use Case Handlers

// NewRegisterUserHandler creates a handler for user registration
func (cr *CompositionRoot) NewRegisterUserHandler() *auth.RegisterUserHandler {
	return auth.NewRegisterUserHandler(cr.GetUnitOfWork(), cr.EventPublisher(), cr.JWTService())
}

// NewLoginUserHandler creates a handler for user login
func (cr *CompositionRoot) NewLoginUserHandler() *auth.LoginUserHandler {
	return auth.NewLoginUserHandler(cr.GetUnitOfWork(), cr.EventPublisher(), cr.JWTService())
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
