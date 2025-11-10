# Quest Auth - API Documentation

## 📚 API Overview

Quest Auth provides two APIs:
1. **REST API** (HTTP) - User-facing authentication operations
2. **gRPC API** - Token validation for microservices

---

## 🔐 Authentication

### HTTP API
No authentication required for registration and login endpoints (they return tokens).

### gRPC API
Requires JWT token in request for validation.

---

## 📡 REST API Endpoints

### Base URL
```
http://localhost:8080
```

### Health Check

**GET /health**

Returns service health status.

**Response 200:**
```json
{
  "status": "ok"
}
```

---

### User Registration

**POST /api/v1/auth/register**

Register a new user and receive JWT tokens.

**Request:**
```json
{
  "email": "user@example.com",
  "phone": "+1234567890",
  "name": "John Doe",
  "password": "securepassword123"
}
```

**Response 201:**
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "phone": "+1234567890",
    "name": "John Doe"
  },
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 900
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "phone": "+1234567890",
    "name": "John Doe",
    "password": "securepassword123"
  }'
```

---

### User Login

**POST /api/v1/auth/login**

Authenticate with email/password and receive JWT tokens.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response 200:**
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "phone": "+1234567890",
    "name": "John Doe"
  },
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 900
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword123"
  }'
```

---

## 🔌 gRPC API

### AuthService

**Package:** `auth.v1`

**Service:** `AuthService`

### Authenticate

Validate JWT token and return user information.

**Method:** `Authenticate`

**Request:**
```protobuf
message AuthenticateRequest {
  string jwt_token = 1;
}
```

**Response:**
```protobuf
message AuthenticateResponse {
  User user = 1;
}

message User {
  string id = 1;
  string email = 2;
  string name = 3;
  string phone = 4;
  string created_at = 5;
}
```

**Proto File:** `api/grpc/proto/auth/v1/auth.proto`

---

## ❌ Error Responses

All errors follow RFC 7807 Problem Details format.

**Error Response:**
```json
{
  "type": "about:blank",
  "title": "Validation Error",
  "status": 400,
  "detail": "email already exists"
}
```

### HTTP Status Codes
- `200` - Success
- `201` - Created
- `400` - Bad Request (validation error)
- `401` - Unauthorized (invalid credentials)
- `409` - Conflict (email/phone already exists)
- `500` - Internal Server Error

### gRPC Status Codes
- `OK` - Success
- `INVALID_ARGUMENT` - Invalid request
- `UNAUTHENTICATED` - Invalid/expired token
- `NOT_FOUND` - User not found
- `INTERNAL` - Server error

---

## 🔑 JWT Token Format

### Claims
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "name": "John Doe",
  "phone": "+1234567890",
  "created_at": 1699632000,
  "exp": 1699632900,
  "iat": 1699632000
}
```

### Token Types
- **Access Token**: Short-lived (15 minutes), used for API requests
- **Refresh Token**: Long-lived (7 days), used to obtain new access tokens

### Token Usage
```bash
curl -X GET https://api.example.com/resource \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

## 📖 OpenAPI Specification

Full OpenAPI 3.0 specification available at:
```
api/http/auth/v1/openapi.yaml
```

Generate client code:
```bash
oapi-codegen -config config.yaml api/http/auth/v1/openapi.yaml
```

---

## 📦 gRPC Code Generation

Generate client code:
```bash
cd api/grpc
buf generate
```

Or use pre-generated SDK:
```go
import authv1 "github.com/Vi-72/quest-auth/api/grpc/sdk/go/auth/v1"
```

---

**API Version:** 1.0.0  
**Last Updated:** November 10, 2025
