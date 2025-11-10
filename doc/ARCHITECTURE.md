# Quest Auth - Architecture Overview

## 🏗️ High-Level Architecture

Quest Auth is a backend authentication and authorization service built following **Clean Architecture** and **Domain-Driven Design (DDD)** principles.

**Core Principles:**
- Domain-centric design
- Dependency inversion
- Separation of concerns
- CQRS (Command Query Responsibility Segregation)
- Event-driven architecture

---

## 🎯 Architectural Layers

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

**Dependency Rule:** Dependencies point **inward** (outer layers depend on inner layers).

---

## 🎨 Design Patterns

### 1. Clean Architecture (Uncle Bob)
**Goal:** Separation of concerns, testability, independence from frameworks

**Layers:**
- **Domain** - Pure business logic (no external dependencies)
- **Application** - Use cases orchestration
- **Infrastructure** - External systems (database, JWT, bcrypt)
- **Presentation** - HTTP/gRPC handlers, API

**Benefits:**
- Testable (domain tests with no mocks)
- Flexible (swap implementations easily)
- Maintainable (clear boundaries)

---

### 2. Domain-Driven Design (DDD)
**Goal:** Model complex business logic in code

**Tactical Patterns:**
- **Aggregates:** User (enforce invariants)
- **Entities:** Objects with identity
- **Value Objects:** Email, Phone (immutable)
- **Domain Events:** UserRegistered, UserLoggedIn
- **Repositories:** Data access abstraction

**Strategic Patterns:**
- **Bounded Context:** Authentication domain
- **Ubiquitous Language:** User, Registration, Authentication

---

### 3. Hexagonal Architecture (Ports & Adapters)
**Goal:** Isolate core logic from external dependencies

```
         ┌─────────────────────┐
         │  External Systems   │
         │  (HTTP, Database)   │
         └──────────┬──────────┘
                    │
         ┌──────────▼──────────┐
         │      Adapters       │  ← Infrastructure
         │  (Implementation)   │
         └──────────┬──────────┘
                    │
         ┌──────────▼──────────┐
         │       Ports         │  ← Interfaces
         │   (Interfaces)      │
         └──────────┬──────────┘
                    │
         ┌──────────▼──────────┐
         │   Core Domain       │  ← Business Logic
         │  (Pure Logic)       │
         └─────────────────────┘
```

**Ports (Interfaces):**
- `UserRepository`
- `TransactionManager`
- `EventPublisher`
- `JWTService`
- `PasswordHasher`

**Adapters (Implementations):**
- PostgreSQL repositories
- bcrypt password hasher
- JWT service (HS256)
- HTTP handlers
- gRPC handlers

---

### 4. CQRS (Command Query Responsibility Segregation)
**Goal:** Separate read and write operations

**Commands (Write):**
- `RegisterUserCommand` → Modify state
- `LoginUserCommand` → Modify state
- Use transactions and events

**Queries (Read):**
- `AuthenticateByTokenQuery` → Read state
- No transactions needed, faster

**Benefits:**
- Optimized read/write operations
- Independent scaling
- Clear separation of concerns

---

### 5. Event-Driven Architecture
**Goal:** Communicate changes through events

**Event Flow:**
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

**Events:**
- `UserRegistered`
- `UserLoggedIn`
- `UserPhoneChanged`
- `UserNameChanged`
- `UserPasswordChanged`

**Benefits:**
- Audit trail
- System integration points
- Future: Message queue integration

---

## 🔐 Security Architecture

### Authentication Flow
```
Client Request + Credentials
    ↓
HTTP/gRPC Handler
    ↓
Use Case Handler
    ↓
User Repository
    ↓
Password Verification (bcrypt)
    ↓
JWT Token Generation (HS256)
    ↓
Token Returned to Client
```

### Token Validation Flow
```
Service Request + JWT Token
    ↓
gRPC AuthService.Authenticate
    ↓
JWT Token Validation
    ↓
Extract User Claims
    ↓
Return User Info
    ↓
Service Uses User ID
```

### Security Layers
1. **Transport:** HTTPS (recommended for production)
2. **Authentication:** JWT Bearer tokens
3. **Password Storage:** bcrypt hashing
4. **Input Validation:** Multi-layer (OpenAPI + Domain + Database)
5. **Error Handling:** No sensitive data in error messages

---

## 🗄️ Data Architecture

### Database Schema
```
┌──────────────┐
│   users      │
│              │
│ - id         │
│ - email      │
│ - phone      │
│ - name       │
│ - password   │
│ - created_at │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   events     │
│              │
│ - id         │
│ - event_type │
│ - agg_id     │
│ - data       │
│ - created_at │
└──────────────┘
```

**Relationships:**
- Events → User (aggregate_id references user.id)

**Constraints:**
- UNIQUE on email and phone
- NOT NULL on required fields
- UUID for all IDs
- Timestamps (created_at, updated_at)

---

## 🔄 Request Lifecycle

### Complete Request Flow

```
1. HTTP/gRPC Request arrives
   ↓
2. Router matches route
   ↓
3. OpenAPI Validation (HTTP only)
   - Validate request schema
   - Check required fields
   - Validate formats
   ↓
4. HTTP/gRPC Handler
   - Extract request data
   - Build command/query
   - Call use case handler
   ↓
5. Use Case Handler
   - Begin transaction (for commands)
   - Validate domain rules
   - Execute business logic
   - Save changes
   - Publish events
   - Commit transaction
   ↓
6. Response Mapping
   - Convert domain → API models
   - Generate JWT tokens (for auth operations)
   - Format as JSON/Protobuf
   ↓
7. HTTP/gRPC Response
   - Return to client
```

**Timing (approximate):**
- Validation: ~2ms
- Handler: ~1ms
- Use case: ~10-150ms (depends on operation)
- Total: ~15-155ms per request

---

## 🧩 Component Integration

### Dependency Injection (Composition Root)

**Pattern:** Factory-based dependency injection

```go
type CompositionRoot struct {
    configs        Config
    db             *gorm.DB
    
    // Dependencies
    txManager      ports.TransactionManager
    jwtService     ports.JWTService
    passwordHasher ports.PasswordHasher
    clock          ports.Clock
}
```

**Lifecycle:**
- **Singleton:** JWTService, PasswordHasher, Clock
- **Per-Request:** TransactionManager repositories
- **Stateless:** All components (thread-safe)

---

## 📊 Scalability Considerations

### Horizontal Scaling
✅ **Stateless Design:**
- No in-memory session storage
- JWT tokens contain all user context
- No shared state between instances
- Database is single source of truth

✅ **Load Balancing:**
- Round-robin across instances
- Health checks for readiness
- Graceful shutdown

### Vertical Scaling
- Database connection pooling
- Efficient SQL queries with indexes
- Transaction optimization

### Performance Optimizations
- **JWT Validation**: Signature-only verification (no DB lookup)
- **Password Hashing**: bcrypt with default cost
- **Connection Pooling**: Reuse DB connections
- **Query Optimization**: Indexed unique constraints

---

## 🔮 Future Architecture Evolution

### Phase 1: Current (v1.0.0)
- Monolithic API
- Single database
- Sync operations
- Basic event storage

### Phase 2: Enhanced (v2.0)
- Message queue for events (RabbitMQ/Kafka)
- Redis caching for token validation
- Rate limiting
- Metrics & observability

### Phase 3: Distributed (v3.0)
- Microservices (separate user management)
- Event sourcing
- CQRS with separate read models
- Multi-tenancy support

---

## 🎯 Architecture Decision Records (ADRs)

### ADR-001: Clean Architecture + DDD
**Decision:** Use Clean Architecture with DDD tactical patterns  
**Rationale:** Clear boundaries, testability, business-centric design  
**Status:** Accepted

### ADR-002: CQRS
**Decision:** Separate commands and queries  
**Rationale:** Different optimization strategies, clearer code  
**Status:** Accepted

### ADR-003: Event Storage in PostgreSQL
**Decision:** Store events in same database as users  
**Rationale:** Transactional consistency, simpler infrastructure  
**Status:** Accepted  
**Future:** May migrate to message broker

### ADR-004: JWT Authentication (HS256)
**Decision:** Use JWT tokens with HS256 signing  
**Rationale:** Stateless authentication, microservice integration  
**Status:** Accepted

### ADR-005: bcrypt for Password Hashing
**Decision:** Use bcrypt for password storage  
**Rationale:** Industry standard, secure, proven  
**Status:** Accepted

### ADR-006: gRPC Service
**Decision:** Provide gRPC API for token validation  
**Rationale:** High-performance, type-safe microservice integration  
**Status:** Accepted

### ADR-007: TransactionManager Pattern
**Decision:** Use closure-based transactions (ThreeDots Labs pattern)  
**Rationale:** Simpler than UnitOfWork, GORM manages lifecycle  
**Status:** Accepted (v1.0.0)

---

## 🔗 Microservice Interaction

### High-level interactions
```
┌────────────┐        HTTP/JSON         ┌──────────────────────┐
│   Client   │ ───────────────────────▶ │   Quest Auth API     │
└────────────┘                          │  (Register, Login)   │
                                        └──────────┬───────────┘
                                                   │
                                                   │ gRPC (JWT validation)
                                                   ▼
                                        ┌──────────────────────┐
│  Quest Manager │ ───────────────────▶│  AuthService (gRPC)  │
│ Quest Infra    │                     └──────────────────────┘
                                                   
                                                   │
                                                   │ SQL (Tx via TM)
                                                   ▼
                                        ┌──────────────────────┐
                                        │     PostgreSQL       │
                                        │  users, events       │
                                        └──────────────────────┘
```

### Transaction and event publishing flow
```
CommandHandler
  ↓
TransactionManager.RunInTransaction(ctx, fn)
  - GORM Transaction begins
  - Create repository instances with transaction
  - Execute business logic closure
  - Publish events synchronously in same transaction
  - Commit or rollback automatically

EventPublisher (Publish)
  - Writes to events table within transaction
  - Events are part of same transaction as domain changes
```

### Query path
- Queries use JWT validation without database access
- JWT signature verification only
- No transaction overhead for token validation

### Notes
- TransactionManager uses closure pattern (ThreeDots Labs style)
- GORM manages transaction lifecycle automatically
- All repositories within closure share the same transaction

---

## 📐 Quality Attributes

### Maintainability
- **Score:** ⭐⭐⭐⭐⭐
- Clear layer separation
- Comprehensive tests
- Good documentation

### Testability
- **Score:** ⭐⭐⭐⭐⭐
- Tests across all layers
- Domain tests without mocks
- Fast test execution

### Performance
- **Score:** ⭐⭐⭐⭐
- <50ms token validation
- <200ms registration
- <150ms login

### Security
- **Score:** ⭐⭐⭐⭐⭐
- JWT authentication
- bcrypt password hashing
- Input validation
- Error sanitization

### Scalability
- **Score:** ⭐⭐⭐⭐
- Stateless design
- Horizontal scaling ready
- Connection pooling

---

## 🔗 Related Documentation

For detailed information, see:
- [**Components**](COMPONENTS.md) - Detailed component breakdown
- [**API Documentation**](API.md) - API reference
- [**Domain Events**](DOMAIN_EVENTS.md) - Event system details
- [**Testing**](TESTING.md) - Testing strategies
- [**Development**](DEVELOPMENT.md) - Development guide

---

**Architecture Version:** 1.0.0  
**Last Updated:** November 10, 2025  
**Status:** Production Ready ✅

