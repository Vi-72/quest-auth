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
	"github.com/Vi-72/quest-auth/internal/adapters/out/postgres/eventrepo"
	timeadapter "github.com/Vi-72/quest-auth/internal/adapters/out/time"
	"github.com/Vi-72/quest-auth/internal/core/application/usecases/commands"
	"github.com/Vi-72/quest-auth/internal/core/application/usecases/queries"
	"github.com/Vi-72/quest-auth/internal/core/ports"

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
	return commands.NewRegisterUserHandler(
		cr.GetUnitOfWork(),
		cr.EventPublisher(),
		cr.JWTService(),
		cr.PasswordHasher(),
		cr.Clock(),
	)
}

// NewLoginUserHandler creates a handler for user login
func (cr *CompositionRoot) NewLoginUserHandler() *commands.LoginUserHandler {
	return commands.NewLoginUserHandler(
		cr.GetUnitOfWork(),
		cr.EventPublisher(),
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
