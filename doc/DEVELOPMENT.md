# Quest Auth - Development Guide

## 🚀 Quick Start

### Prerequisites
- Go 1.23+
- PostgreSQL 14+
- Docker & Docker Compose
- Make

### Initial Setup
```bash
# Clone repository
git clone <repository>
cd quest-auth

# Copy environment file
cp config.example .env

# Start PostgreSQL
docker compose up -d postgres

# Run application
go run ./cmd/app
```

### Verify Installation
```bash
curl http://localhost:8080/health
```

---

## 📁 Project Structure

```
quest-auth/
├── api/                      # API specifications
│   ├── grpc/                # gRPC proto files
│   └── http/                # OpenAPI specs
├── cmd/                     # Application entry points
│   └── app/                # Main application
├── internal/               # Internal packages
│   ├── adapters/          # External integrations
│   │   ├── in/           # Inbound (HTTP, gRPC)
│   │   └── out/          # Outbound (DB, JWT)
│   ├── core/             # Core business logic
│   │   ├── application/  # Use cases
│   │   ├── domain/       # Domain models
│   │   └── ports/        # Interfaces
│   └── pkg/              # Shared packages
├── tests/                # Test suites
└── doc/                  # Documentation
```

---

## 🛠️ Development Workflow

### 1. Feature Development
```bash
# 1. Create feature branch
git checkout -b feature/my-feature

# 2. Implement changes following Clean Architecture

# 3. Run tests
make test

# 4. Run linter
make lint

# 5. Commit changes
git commit -m "Add feature: description"
```

### 2. Code Generation
```bash
# Generate OpenAPI code
make gen-api

# Generate gRPC code
cd api/grpc && make generate
```

### 3. Database Changes
```bash
# Update schema in migrations
# Run migrations
# Update repository methods
```

---

## 🧪 Testing During Development

```bash
# Run all tests
make test

# Run specific test
go test -v ./tests/domain -run TestUser

# Watch mode (with entr)
ls **/*.go | entr -c go test ./tests/domain
```

---

## 🔧 Common Tasks

### Add New Use Case
1. Define command/query in `internal/core/application/usecases/`
2. Implement handler
3. Wire in `cmd/composition_root.go`
4. Add HTTP/gRPC handler
5. Write tests

### Add New Domain Event
1. Define event in `internal/core/domain/model/auth/events.go`
2. Publish in aggregate method
3. Handle in event repository
4. Add to event catalog docs

### Update API
1. Modify OpenAPI spec: `api/http/auth/v1/openapi.yaml`
2. Regenerate code: `make gen-api`
3. Update handlers
4. Update tests

---

## 📊 Code Quality

### Linting
```bash
make lint
# or
golangci-lint run
```

### Formatting
```bash
make fmt
# or
goimports -w .
```

---

## 🐛 Debugging

### Run with Debugger (Delve)
```bash
dlv debug ./cmd/app
```

### View Logs
```bash
# Application logs
tail -f app.log

# Database logs
docker compose logs -f postgres
```

---

## 📚 Resources

- [Architecture](ARCHITECTURE.md)
- [Testing](TESTING.md)
- [API Documentation](API.md)

**Last Updated:** November 10, 2025
