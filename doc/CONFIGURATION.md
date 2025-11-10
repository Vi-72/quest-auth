# Quest Auth - Configuration Guide

## 🔧 Environment Variables

### HTTP Server
```bash
HTTP_PORT=8080                    # HTTP server port
```

### gRPC Server
```bash
GRPC_PORT=9090                    # gRPC server port
```

### Database
```bash
DB_HOST=localhost                 # PostgreSQL host
DB_PORT=5432                      # PostgreSQL port
DB_USER=postgres                  # Database user
DB_PASSWORD=password              # Database password
DB_NAME=auth                      # Database name
DB_SSLMODE=disable                # SSL mode (disable/require)
```

### JWT Configuration
```bash
JWT_SECRET_KEY=your-secret-key    # Secret key for JWT signing (CHANGE IN PRODUCTION!)
JWT_ACCESS_TOKEN_DURATION=15      # Access token duration (minutes)
JWT_REFRESH_TOKEN_DURATION=168    # Refresh token duration (hours, 7 days)
```

### Event Processing
```bash
EVENT_GOROUTINE_LIMIT=10          # Max concurrent event processing goroutines
```

---

## 📝 Configuration Files

### `config.example`
Template for environment configuration. Copy to `.env`:
```bash
cp config.example .env
```

---

## 🔐 Security Settings

### Production JWT Secret
**IMPORTANT:** Generate strong secret key for production:
```bash
openssl rand -base64 32
```

### Password Hashing
bcrypt cost is set to default (10) in code. Higher cost = more secure but slower.

---

## 🚀 Environment-Specific Configuration

### Development
```bash
HTTP_PORT=8080
GRPC_PORT=9090
DB_HOST=localhost
DB_PORT=5432
JWT_SECRET_KEY=dev-secret-key-change-in-production
JWT_ACCESS_TOKEN_DURATION=15
JWT_REFRESH_TOKEN_DURATION=168
```

### Production
```bash
HTTP_PORT=8080
GRPC_PORT=9090
DB_HOST=production-db-host
DB_PORT=5432
DB_SSLMODE=require
JWT_SECRET_KEY=<generated-strong-secret>
JWT_ACCESS_TOKEN_DURATION=15
JWT_REFRESH_TOKEN_DURATION=168
```

---

**Last Updated:** November 10, 2025
