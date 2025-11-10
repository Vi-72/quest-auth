# Quest Auth - Testing Guide

## 🧪 Test Strategy

Quest Auth uses a comprehensive test pyramid with tests at multiple levels.

---

## Test Pyramid

```
       ╱╲
      ╱  ╲ E2E Tests
     ╱────╲
    ╱      ╲ Integration Tests
   ╱────────╲
  ╱          ╲ Contract Tests
 ╱────────────╲
╱              ╲ Unit (Domain) Tests
```

---

## 📁 Test Organization

```
tests/
├── domain/                    # Unit tests
│   ├── user_test.go
│   ├── email_test.go
│   └── phone_test.go
├── contracts/                 # Contract tests
│   ├── repository_contracts_test.go
│   └── event_publisher_contracts_test.go
└── integration/              # Integration & E2E tests
    └── tests/
        ├── auth_http_tests/
        ├── auth_grpc_tests/
        ├── auth_handler_tests/
        ├── auth_e2e_tests/
        └── repository_tests/
```

---

## 🏃 Running Tests

### All Tests
```bash
make test
```

### Unit Tests (Domain)
```bash
make test-unit
# or
go test ./tests/domain -v
```

### Contract Tests
```bash
make test-contracts
# or
go test ./tests/contracts -v
```

### Integration Tests
```bash
make test-integration
# or
go test -tags=integration ./tests/integration/... -v -p 1
```

### Specific Test Suites
```bash
# HTTP tests
go test -tags=integration ./tests/integration/tests/auth_http_tests -v

# gRPC tests
go test -tags=integration ./tests/integration/tests/auth_grpc_tests -v

# E2E tests
go test -tags=integration ./tests/integration/tests/auth_e2e_tests -v

# Repository tests
go test -tags=integration ./tests/integration/tests/repository_tests -v
```

---

## 🎯 Test Coverage

### Generate Coverage Report
```bash
make test-coverage
# or
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Target Coverage
- **Overall:** >70%
- **Domain Layer:** 100%
- **Application Layer:** >80%
- **Infrastructure Layer:** >60%

---

## 🧩 Testing Patterns

### Domain Tests
Pure business logic tests without external dependencies.

```go
func TestUser_VerifyPassword(t *testing.T) {
    hasher := bcryptadapter.NewHasher()
    user, _ := auth.NewUser(email, phone, "John", "password123", hasher, clock)
    
    assert.True(t, user.VerifyPassword("password123", hasher))
    assert.False(t, user.VerifyPassword("wrongpassword", hasher))
}
```

### Contract Tests
Verify interface implementations comply with contracts.

### Integration Tests
Full stack tests with PostgreSQL.

**Setup:**
- PostgreSQL test database (port 5433)
- Clean database before each test
- Test DI container with real dependencies

---

## 🔧 Test Configuration

### Environment Variables
```bash
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=auth_test
DB_SSLMODE=disable
JWT_SECRET_KEY=test-secret-key
```

### Test Database
```bash
# Start test PostgreSQL
docker compose -f docker-compose.yml up -d postgres

# Run migrations
go run ./cmd/app
```

---

## 🐛 Debugging Tests

### Verbose Output
```bash
go test -v ./tests/domain
```

### Run Single Test
```bash
go test -v ./tests/domain -run TestUser_VerifyPassword
```

### With Coverage
```bash
go test -v -cover ./tests/domain
```

---

**Last Updated:** November 10, 2025
