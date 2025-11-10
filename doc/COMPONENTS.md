# Quest Auth - Components

## 🏗️ Component Overview

Quest Auth follows Clean Architecture with clear separation of concerns across 4 layers.

---

## 1. Domain Layer (`internal/core/domain`)

### User Aggregate (`model/auth/user.go`)
Central business entity with authentication logic.

**Responsibilities:**
- User registration
- Password verification
- Domain event generation
- Profile updates

**Key Methods:**
- `NewUser()` - Create new user with validation
- `VerifyPassword()` - Compare passwords with bcrypt
- `ChangeName()` - Update user name
- `ChangePhone()` - Update phone number
- `ChangePassword()` - Update password
- `MarkLoggedIn()` - Record login event

### Value Objects (`model/kernel/`)
- `Email` - Validated email address
- `Phone` - Validated phone number with international format
- `JWTToken` - Immutable token value

### Domain Events (`model/auth/events.go`)
- `UserRegistered`
- `UserLoggedIn`
- `UserPhoneChanged`
- `UserNameChanged`
- `UserPasswordChanged`

---

## 2. Application Layer (`internal/core/application`)

### Commands (`usecases/commands/`)

**RegisterUserHandler:**
- Creates new user
- Validates uniqueness (email, phone)
- Publishes UserRegistered event
- Generates JWT tokens

**LoginUserHandler:**
- Validates credentials
- Publishes UserLoggedIn event
- Generates JWT tokens

### Queries (`usecases/queries/`)

**AuthenticateByTokenHandler:**
- Validates JWT token
- Extracts user claims
- Returns user information

---

## 3. Infrastructure Layer (`internal/adapters/out`)

### Repositories (`postgres/`)
- **UserRepository** - User CRUD operations
- **EventRepository** - Event persistence

### Services
- **JWTService** (`jwt/`) - Token generation and validation
- **PasswordHasher** (`bcrypt/`) - Password hashing
- **Clock** (`time/`) - Time operations

### Transaction Management
- **TransactionManager** (`postgres/`) - Closure-based transactions

---

## 4. Presentation Layer (`internal/adapters/in`)

### HTTP Handlers (`http/`)
- `RegisterHandler` - POST /api/v1/auth/register
- `LoginHandler` - POST /api/v1/auth/login

### gRPC Handlers (`grpc/`)
- `AuthHandler` - AuthService.Authenticate

---

## 🔗 Component Dependencies

```
HTTP/gRPC Handlers
    ↓
Use Case Handlers (Commands/Queries)
    ↓
Domain Models (User Aggregate)
    ↓
Repositories (User, Event)
    ↓
PostgreSQL
```

---

## 📦 Ports (Interfaces)

Located in `internal/core/ports/`:

- `UserRepository` - User data access
- `EventPublisher` - Event publishing
- `TransactionManager` - Transaction coordination
- `JWTService` - Token operations
- `PasswordHasher` - Password operations
- `Clock` - Time operations

---

**Last Updated:** November 10, 2025
