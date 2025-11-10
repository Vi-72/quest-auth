# Quest Auth - Technical Context

## 🛠️ Technology Stack

### Core Technologies
- **Language**: Go 1.23+
- **Database**: PostgreSQL 14+ 
- **API**: REST (HTTP) and gRPC
- **ORM**: GORM v2 for database operations
- **Password Hashing**: bcrypt
- **JWT**: HS256 signing algorithm

### Key Dependencies
```go
// Core Framework
github.com/go-chi/chi/v5          // HTTP router
github.com/gofrs/uuid/v5          // UUID generation
github.com/oapi-codegen/oapi-codegen // OpenAPI code generation

// Database & ORM
gorm.io/gorm                      // ORM
gorm.io/driver/postgres           // PostgreSQL driver

// gRPC & Protocol Buffers
google.golang.org/grpc            // gRPC server
google.golang.org/protobuf        // Protocol buffers

// Security
github.com/golang-jwt/jwt/v5      // JWT tokens
golang.org/x/crypto/bcrypt        // Password hashing

// Validation & Error Handling
github.com/getkin/kin-openapi/openapi3 // OpenAPI validation
github.com/pkg/errors             // Error wrapping

// Testing
github.com/stretchr/testify       // Test assertions
```

## 🏗️ Development Setup

### Prerequisites
```bash
# Required Software
Go 1.23+
PostgreSQL 14+
Docker & Docker Compose
Make

# Optional (for development)
golangci-lint
oapi-codegen
buf (for gRPC)
```

### Environment Configuration
```bash
# HTTP Server
HTTP_PORT=8080

# gRPC Server
GRPC_PORT=9090

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=auth
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET_KEY=your-secret-key-here
JWT_ACCESS_TOKEN_DURATION=15    # minutes
JWT_REFRESH_TOKEN_DURATION=168  # hours (7 days)

# Event Processing
EVENT_GOROUTINE_LIMIT=10
```

### Quick Start
```bash
# 1. Clone and setup
git clone <repository>
cd quest-auth
cp config.example .env

# 2. Start dependencies
docker compose up -d postgres

# 3. Run application
go run ./cmd/app

# 4. Verify
curl http://localhost:8080/health
```

## 🔧 Build & Development Tools

### Make Commands
```bash
# Development
make run              # Run application
make build            # Build binary
make clean            # Clean build artifacts

# Testing
make test             # Run all tests
make test-unit        # Unit tests only
make test-contracts   # Contract tests
make test-integration # Integration tests
make test-coverage    # Coverage report

# Code Quality
make lint             # Run golangci-lint
make fmt              # Format code
make generate         # Generate OpenAPI code

# Database
make db-migrate       # Run migrations
```

### Code Generation

#### OpenAPI Code Generation
```bash
# Generate HTTP server from OpenAPI spec
make gen-api
# or
oapi-codegen -config api/http/auth/v1/config.yaml api/http/auth/v1/openapi.yaml

# Generated Files
api/http/auth/v1/servers.gen.go  # HTTP server interface
```

#### gRPC Code Generation
```bash
# Generate gRPC server and client from proto files
cd api/grpc
make generate

# Generated Files
api/grpc/sdk/go/auth/v1/auth.pb.go        # Protocol buffer messages
api/grpc/sdk/go/auth/v1/auth_grpc.pb.go   # gRPC server and client
```

## 🗄️ Database Schema

### Core Tables
```sql
-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Events table for domain events
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type VARCHAR(255) NOT NULL,
    aggregate_id UUID NOT NULL,
    data JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### Performance Indexes
```sql
-- User indexes for fast lookups
CREATE UNIQUE INDEX idx_users_email ON users(email);
CREATE UNIQUE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_created_at ON users(created_at);

-- Event indexes for queries
CREATE INDEX idx_events_aggregate_id ON events(aggregate_id);
CREATE INDEX idx_events_event_type ON events(event_type);
CREATE INDEX idx_events_created_at ON events(created_at);
```

## 🔐 Security Configuration

### JWT Configuration
```go
type JWTConfig struct {
    SecretKey            string
    AccessTokenDuration  time.Duration  // 15 minutes default
    RefreshTokenDuration time.Duration  // 7 days default
    SigningMethod        string         // HS256
}

// Token Claims
type JWTClaims struct {
    UserID    string `json:"user_id"`
    Email     string `json:"email"`
    Name      string `json:"name"`
    Phone     string `json:"phone"`
    CreatedAt int64  `json:"created_at"`
    ExpiresAt int64  `json:"exp"`
    IssuedAt  int64  `json:"iat"`
}
```

### Password Hashing Configuration
```go
// bcrypt configuration
const (
    BcryptCost = bcrypt.DefaultCost  // Cost factor 10
)

// Password Requirements
// - Minimum length: 8 characters
// - No maximum length
// - All characters allowed
```

## 🧪 Testing Infrastructure

### Test Database Setup
```bash
# Docker Compose for testing
version: '3.8'
services:
  postgres-test:
    image: postgres:14
    environment:
      POSTGRES_DB: auth_test
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    ports:
      - "5433:5432"
    volumes:
      - postgres_test_data:/var/lib/postgresql/data

volumes:
  postgres_test_data:
```

### Test Configuration
```go
// Test environment variables
func getTestConfig() cmd.Config {
    return cmd.Config{
        HTTPPort: "8080",
        GrpcPort: "9090",
        DBHost:   "localhost",
        DBPort:   "5433",
        DBUser:   "postgres",
        DBPassword: "password",
        DBName:   "auth_test",
        DBSslMode: "disable",
        JWTSecretKey: "test-secret-key",
        JWTAccessTokenDuration: 1,  // 1 minute for tests
        JWTRefreshTokenDuration: 24, // 24 hours for tests
    }
}
```

### Test Categories
```bash
# Unit Tests (Domain Layer)
go test ./tests/domain -v

# Contract Tests (Interface Layer)
go test ./tests/contracts -v

# Integration Tests (Full Stack)
go test -tags=integration ./tests/integration/... -v -p 1

# HTTP Tests
go test -tags=integration ./tests/integration/tests/auth_http_tests -v

# gRPC Tests
go test -tags=integration ./tests/integration/tests/auth_grpc_tests -v

# E2E Tests
go test -tags=integration ./tests/integration/tests/auth_e2e_tests -v

# Repository Tests
go test -tags=integration ./tests/integration/tests/repository_tests -v
```

## 📊 Monitoring & Observability

### Health Checks
```go
// Health endpoint
func Health(w http.ResponseWriter, r *http.Request) {
    health := map[string]string{
        "status":    "healthy",
        "timestamp": time.Now().UTC().Format(time.RFC3339),
        "version":   "1.0.0",
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(health)
}
```

### Logging Configuration
```go
// Structured logging with slog
func NewLogger() *slog.Logger {
    return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))
}
```

### Error Handling
```go
// Problem Details (RFC 7807)
type ProblemDetails struct {
    Type     string `json:"type"`
    Title    string `json:"title"`
    Status   int    `json:"status"`
    Detail   string `json:"detail"`
    Instance string `json:"instance,omitempty"`
}

func WriteProblem(w http.ResponseWriter, problem ProblemDetails) {
    w.Header().Set("Content-Type", "application/problem+json")
    w.WriteHeader(problem.Status)
    json.NewEncoder(w).Encode(problem)
}
```

## 🚀 Deployment Configuration

### Docker Configuration
```dockerfile
# Multi-stage build
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o quest-auth ./cmd/app

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/quest-auth .
EXPOSE 8080 9090
CMD ["./quest-auth"]
```

### Docker Compose
```yaml
version: '3.8'
services:
  quest-auth:
    build: .
    ports:
      - "8080:8080"  # HTTP
      - "9090:9090"  # gRPC
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=secret
      - DB_NAME=auth
      - JWT_SECRET_KEY=${JWT_SECRET_KEY}
    depends_on:
      - postgres

  postgres:
    image: postgres:14
    environment:
      POSTGRES_DB: auth
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: secret
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

volumes:
  postgres_data:
```

## 🔧 Development Constraints

### Code Style
- **Go Version**: 1.23+ required
- **Linting**: golangci-lint with strict rules
- **Formatting**: goimports for import organization
- **Comments**: English language, godoc format

### Performance Requirements
- **Token Validation**: <50ms per request
- **User Registration**: <200ms per request
- **User Login**: <150ms per request
- **Memory Usage**: <50MB per instance
- **Database Connections**: Max 25 concurrent connections

### Security Requirements
- **Password Hashing**: bcrypt with default cost (10)
- **JWT Signing**: HS256 algorithm
- **Token Expiry**: Access (15 min), Refresh (7 days)
- **Input Validation**: Multi-layer (OpenAPI + Domain + Database)
- **Error Handling**: No sensitive data in error messages

## 🧩 Application Behaviors (Transaction Management)

- **TransactionManager**: Closure-based transaction management using `RunInTransaction` method
- **Command Handlers**: Use `TransactionManager.RunInTransaction` with closure pattern
- **Query Handlers**: Use bare repositories without transactions for read operations
- **Event Publisher**: Publishes events synchronously within the same transaction as domain changes

---

**This technical context provides the foundation for development, testing, and deployment of Quest Auth, ensuring consistent practices across the development team.**

