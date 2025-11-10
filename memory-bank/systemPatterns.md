# Quest Auth - System Patterns

## 🏗️ Architecture Overview

Quest Auth implements **Clean Architecture** with **Domain-Driven Design** principles, creating a maintainable, testable, and secure authentication system.

## 🎯 Core Architectural Patterns

### 1. Clean Architecture (Uncle Bob)

**Structure**: 4 concentric layers with dependency inversion

```
┌─────────────────────────────────────────────────────────┐
│                   Presentation Layer                     │
│              (HTTP, gRPC, Middleware)                    │
└────────────────────────┬────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────┐
│                  Application Layer                       │
│              (Use Cases, Commands, Queries)              │
└────────────────────────┬────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────┐
│                    Domain Layer                          │
│           (Business Logic, User Aggregate, Events)       │
└────────────────────────┬────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────┐
│                Infrastructure Layer                      │
│         (PostgreSQL, bcrypt, JWT, Repositories)          │
└─────────────────────────────────────────────────────────┘
```

**Key Principles**:
- Dependencies point **inward** (outer layers depend on inner layers)
- Domain layer has **no external dependencies**
- Business logic isolated from infrastructure concerns
- Testable without external systems

### 2. Domain-Driven Design (DDD)

#### Aggregates
- **User**: Central business entity with authentication logic
  - Registration logic
  - Password management
  - Authentication verification
  - Profile updates (name, phone, password changes)

#### Domain Events
```go
// User Events
UserRegistered{UserID, Email, Phone, Name, CreatedAt}
UserLoggedIn{UserID, LoginAt}
UserPhoneChanged{UserID, OldPhone, NewPhone, ChangedAt}
UserNameChanged{UserID, OldName, NewName, ChangedAt}
UserPasswordChanged{UserID, ChangedAt}
```

#### Value Objects
- **Email**: Immutable email with validation
- **Phone**: Immutable phone number with format validation
- **JWTToken**: Immutable token with validation

### 3. Hexagonal Architecture (Ports & Adapters)

**Ports (Interfaces)**:
```go
type UserRepository interface {
    Create(user *User) error
    GetByID(id uuid.UUID) (*User, error)
    GetByEmail(email Email) (*User, error)
    EmailExists(email Email) (bool, error)
    PhoneExists(phone Phone) (bool, error)
    Update(user *User) error
    Delete(id uuid.UUID) error
}

type TransactionManager interface {
    RunInTransaction(ctx context.Context, fn func(ctx context.Context, repos Repositories) error) error
}

type EventPublisher interface {
    Publish(ctx context.Context, events ...DomainEvent) error
}

type JWTService interface {
    GenerateTokenPair(id, email, name, phone string, createdAt time.Time) (TokenPair, error)
    ValidateToken(token string) (*JWTClaims, error)
}

type PasswordHasher interface {
    Hash(password string) (string, error)
    Compare(hashedPassword, password string) bool
}
```

**Adapters (Implementations)**:
- PostgreSQL repositories (userrepo, eventrepo)
- bcrypt password hasher
- JWT service (HS256)
- HTTP handlers
- gRPC handlers

### 4. CQRS (Command Query Responsibility Segregation)

#### Commands (Write Operations)
```go
// Command Handlers
RegisterUserCommand → RegisterUserHandler
LoginUserCommand → LoginUserHandler

// Characteristics:
- Modify state
- Use TransactionManager
- Publish domain events
- Return user data + JWT tokens
```

#### Queries (Read Operations)
```go
// Query Handlers
AuthenticateByTokenQuery → AuthenticateByTokenHandler

// Characteristics:
- Read-only operations
- No transactions needed
- Optimized for performance
- Return user claims from JWT
```

### 5. Event-Driven Architecture

#### Event Flow
```
Domain Operation
    ↓
Add Domain Event
    ↓
Publish Event (in transaction)
    ↓
Event Persisted
    ↓
Transaction Commits
    ↓
Events Cleared
```

#### Event Storage
- Events stored in PostgreSQL `events` table
- Transactional consistency with User aggregate
- Future: Message queue integration (RabbitMQ/Kafka)

## 🔧 Implementation Patterns

### 1. Composition Root Pattern (Dependency Injection)

**Dependency Injection Container**:
```go
type CompositionRoot struct {
    configs        Config
    db             *gorm.DB
    txManager      ports.TransactionManager
    jwtService     ports.JWTService
    passwordHasher ports.PasswordHasher
    clock          ports.Clock
}

// Handlers initialized with dependencies
func (cr *CompositionRoot) NewRegisterUserHandler() *commands.RegisterUserHandler {
    return commands.NewRegisterUserHandler(
        cr.TransactionManager(),
        cr.JWTService(),
        cr.PasswordHasher(),
        cr.Clock(),
    )
}
```

### 2. Transaction Management Pattern (ThreeDots Labs)

**Closure-based Transactions**:
```go
// TransactionManager handles database transactions
type TransactionManager interface {
    RunInTransaction(ctx context.Context, fn func(ctx context.Context, repos Repositories) error) error
}

// Repositories holds all repository instances for a transaction
type Repositories struct {
    User  UserRepository
    Event EventPublisher
}

// Command handler using TransactionManager
func (h *RegisterUserHandler) Handle(ctx context.Context, cmd RegisterUserCommand) (RegisterUserResult, error) {
    email, _ := kernel.NewEmail(cmd.Email)
    phone, _ := kernel.NewPhone(cmd.Phone)
    
    var createdUser auth.User
    
    err := h.txManager.RunInTransaction(ctx, func(ctx context.Context, repos ports.Repositories) error {
        // Check uniqueness
        emailExists, _ := repos.User.EmailExists(email)
        if emailExists {
            return errs.NewDomainValidationError("email", "email already exists")
        }
        
        // Create user
        user, _ := auth.NewUser(email, phone, cmd.Name, cmd.Password, h.passwordHasher, h.clock)
        
        // Save user
        if err := repos.User.Create(&user); err != nil {
            return err
        }
        
        // Publish events synchronously in same transaction
        if err := repos.Event.Publish(ctx, user.GetDomainEvents()...); err != nil {
            return err
        }
        
        createdUser = user
        return nil
    })
    
    // Generate JWT tokens outside transaction
    tokenPair, _ := h.jwtService.GenerateTokenPair(...)
    
    return RegisterUserResult{User: ..., AccessToken: ..., RefreshToken: ...}, err
}
```

**Key Benefits**:
- GORM manages transaction lifecycle automatically
- Explicit transaction boundaries via closure
- Repository instances scoped to transaction
- Simpler than UnitOfWork pattern

### 3. Repository Pattern

**Abstraction Layer**:
```go
type UserRepository interface {
    Create(user *User) error
    GetByID(id uuid.UUID) (*User, error)
    GetByEmail(email Email) (*User, error)
    EmailExists(email Email) (bool, error)
    PhoneExists(phone Email) (bool, error)
}

// PostgreSQL Implementation
type postgresUserRepository struct {
    db *gorm.DB
}

func (r *postgresUserRepository) Create(user *auth.User) error {
    dto := domainToDTO(user)
    return r.db.Create(&dto).Error
}
```

### 4. Middleware Pattern

**HTTP Middleware Chain**:
```go
func (c *CompositionRoot) Middlewares() []func(http.Handler) http.Handler {
    middlewares := []func(http.Handler) http.Handler{
        middleware.Recovery,
        middleware.Logger,
    }
    
    return middlewares
}
```

## 🗄️ Data Patterns

### 1. User Storage

**User Table**:
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

**Event Storage**:
```sql
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type VARCHAR(255) NOT NULL,
    aggregate_id UUID NOT NULL,
    data JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### 2. Password Security

**Bcrypt Hashing**:
```go
type PasswordHasher interface {
    Hash(password string) (string, error)
    Compare(hashedPassword, password string) bool
}

// Implementation
func (h *bcryptHasher) Hash(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

func (h *bcryptHasher) Compare(hashedPassword, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}
```

### 3. JWT Token Management

**Token Generation**:
```go
type TokenPair struct {
    AccessToken  string
    RefreshToken string
    TokenType    string
    ExpiresIn    int64
}

func (s *JWTService) GenerateTokenPair(id, email, name, phone string, createdAt time.Time) (TokenPair, error) {
    // Access token (short-lived)
    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id":    id,
        "email":      email,
        "name":       name,
        "phone":      phone,
        "created_at": createdAt.Unix(),
        "exp":        time.Now().Add(accessTokenDuration).Unix(),
        "iat":        time.Now().Unix(),
    })
    
    // Refresh token (long-lived)
    refreshToken := jwt.NewWithClaims(...)
    
    return TokenPair{...}, nil
}
```

## 🔐 Security Patterns

### 1. Password Management

**Security Flow**:
```
User Registration
    ↓
Password → bcrypt.Hash()
    ↓
Store hashed password
    ↓
Never store plain-text

User Login
    ↓
Retrieve hashed password
    ↓
bcrypt.Compare(stored, provided)
    ↓
Grant/Deny access
```

### 2. JWT Authentication

**Token Flow**:
```
User Login
    ↓
Generate JWT (HS256)
    ↓
Sign with secret key
    ↓
Return to user

Service Request
    ↓
Extract JWT from header/gRPC
    ↓
Validate signature
    ↓
Extract claims
    ↓
Use user_id for authorization
```

### 3. Input Validation

**Multi-Layer Validation**:
```go
// 1. OpenAPI Schema Validation (HTTP Layer)
validationMiddleware.Validate(r) // Format, ranges, required fields

// 2. Domain Validation (Business Layer)
email, err := kernel.NewEmail(rawEmail)
// Email format validation

phone, err := kernel.NewPhone(rawPhone)
// Phone format validation

user, err := auth.NewUser(email, phone, name, password, hasher, clock)
// Business rules validation

// 3. Resource Validation (Application Layer)
emailExists, err := userRepo.EmailExists(email)
// Uniqueness checks
```

## 🧪 Testing Patterns

### 1. Test Pyramid

**Unit Tests** (Domain Layer):
```go
func TestUser_VerifyPassword(t *testing.T) {
    hasher := bcryptadapter.NewHasher()
    clock := timeadapter.NewClock()
    user, _ := auth.NewUser(email, phone, "John", "password123", hasher, clock)
    
    assert.True(t, user.VerifyPassword("password123", hasher))
    assert.False(t, user.VerifyPassword("wrongpassword", hasher))
}
```

**Contract Tests** (Interface Layer):
```go
func TestUserRepositoryContract(t *testing.T) {
    mockRepo := mocks.NewMockUserRepository()
    // Verify interface methods work correctly
}
```

**Integration Tests** (Full Stack):
```go
func TestRegisterUserE2E(t *testing.T) {
    // Full HTTP request → database → response
    resp := httptest.NewRequest("POST", "/api/v1/auth/register", userData)
    assert.Equal(t, 201, resp.StatusCode)
}
```

### 2. Test Data Builders

**Domain Test Data**:
```go
func NewTestUser() *auth.User {
    email, _ := kernel.NewEmail("test@example.com")
    phone, _ := kernel.NewPhone("+1234567890")
    hasher := &mockHasher{}
    clock := &mockClock{}
    
    return auth.NewUser(email, phone, "Test User", "password123", hasher, clock)
}
```

## 🚀 Performance Patterns

### 1. Database Indexing

**Optimized Queries**:
```sql
-- Unique indexes for fast lookups
CREATE UNIQUE INDEX idx_users_email ON users(email);
CREATE UNIQUE INDEX idx_users_phone ON users(phone);

-- Index for event queries
CREATE INDEX idx_events_aggregate_id ON events(aggregate_id);
CREATE INDEX idx_events_created_at ON events(created_at);
```

### 2. Stateless Design

**JWT Benefits**:
- No session storage required
- Horizontal scaling without sticky sessions
- Fast token validation (signature verification only)
- All user context in token claims

---

**These patterns form the foundation of Quest Auth's architecture, ensuring maintainability, testability, security, and scalability while following industry best practices.**

