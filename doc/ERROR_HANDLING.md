# Quest Auth - Error Handling

## 🚨 Error Response Format

All errors follow **RFC 7807 Problem Details** format.

### Problem Details Structure
```json
{
  "type": "about:blank",
  "title": "Error Title",
  "status": 400,
  "detail": "Detailed error message"
}
```

---

## 📋 Error Categories

### 1. Validation Errors (400)

**Email Format Invalid:**
```json
{
  "type": "about:blank",
  "title": "Validation Error",
  "status": 400,
  "detail": "email: invalid email format"
}
```

**Phone Format Invalid:**
```json
{
  "type": "about:blank",
  "title": "Validation Error",
  "status": 400,
  "detail": "phone: invalid phone format"
}
```

**Email Already Exists:**
```json
{
  "type": "about:blank",
  "title": "Validation Error",
  "status": 400,
  "detail": "email already exists"
}
```

**Phone Already Exists:**
```json
{
  "type": "about:blank",
  "title": "Validation Error",
  "status": 400,
  "detail": "phone already exists"
}
```

### 2. Authentication Errors (401)

**Invalid Credentials:**
```json
{
  "type": "about:blank",
  "title": "Authentication Error",
  "status": 401,
  "detail": "invalid email or password"
}
```

**Invalid Token:**
```json
{
  "type": "about:blank",
  "title": "Authentication Error",
  "status": 401,
  "detail": "invalid or expired token"
}
```

### 3. Not Found Errors (404)

**User Not Found:**
```json
{
  "type": "about:blank",
  "title": "Not Found",
  "status": 404,
  "detail": "user not found"
}
```

### 4. Server Errors (500)

**Internal Server Error:**
```json
{
  "type": "about:blank",
  "title": "Internal Server Error",
  "status": 500,
  "detail": "An unexpected error occurred"
}
```

---

## 🔌 gRPC Error Codes

### Mapping

| Domain Error | gRPC Code | HTTP Status |
|--------------|-----------|-------------|
| DomainValidationError | INVALID_ARGUMENT | 400 |
| JWTValidationError | UNAUTHENTICATED | 401 |
| NotFoundError | NOT_FOUND | 404 |
| InfrastructureError | INTERNAL | 500 |

---

## 🛠️ Handling Errors in Clients

### HTTP Client (JavaScript)
```javascript
try {
  const response = await fetch('/api/v1/auth/register', {
    method: 'POST',
    body: JSON.stringify(userData)
  });
  
  if (!response.ok) {
    const problem = await response.json();
    console.error(`${problem.title}: ${problem.detail}`);
  }
} catch (error) {
  console.error('Network error:', error);
}
```

### gRPC Client (Go)
```go
resp, err := client.Authenticate(ctx, &authv1.AuthenticateRequest{
    JwtToken: token,
})
if err != nil {
    if s, ok := status.FromError(err); ok {
        switch s.Code() {
        case codes.Unauthenticated:
            // Handle invalid token
        case codes.InvalidArgument:
            // Handle validation error
        default:
            // Handle other errors
        }
    }
}
```

---

**Last Updated:** November 10, 2025
