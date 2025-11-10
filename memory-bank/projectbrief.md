# Quest Auth - Project Brief

## 🎯 Project Overview

**Quest Auth** is a backend authentication and authorization service providing secure user management, JWT token generation, and identity verification for the quest application ecosystem.

## 🎪 Core Purpose

Quest Auth serves as the centralized authentication hub for the quest-based application ecosystem, enabling:

- **User Lifecycle Management**: Complete user registration and authentication
- **JWT Token Services**: Secure token generation and validation
- **Identity Verification**: gRPC service for token validation by other microservices
- **Event-Driven Architecture**: Real-time tracking of authentication events

## 🏗️ Architectural Foundation

### Design Philosophy
- **Clean Architecture**: Clear separation of concerns across layers
- **Domain-Driven Design**: Business logic centered around user authentication domain
- **Event-Driven**: State changes communicated through domain events
- **Security-First**: Bcrypt password hashing, JWT tokens with industry standards

### Core Principles
1. **Domain-Centric**: Authentication business logic isolated in domain layer
2. **Dependency Inversion**: Outer layers depend on inner layers
3. **CQRS**: Separation of command and query operations
4. **Hexagonal Architecture**: Ports and adapters for external integrations

## 🎯 Key Requirements

### Functional Requirements
- ✅ **User Registration**: Email, phone, name, password (bcrypt hashed)
- ✅ **User Authentication**: Login with email/password
- ✅ **JWT Token Management**: Access and refresh tokens (HS256)
- ✅ **Token Validation**: gRPC service for other microservices
- ✅ **Event Tracking**: Domain events for audit and integration

### Non-Functional Requirements
- **Security**: Bcrypt password hashing, JWT tokens, input validation
- **Performance**: <50ms response time for token validation
- **Scalability**: Horizontal scaling capability
- **Reliability**: Transactional consistency, graceful error handling
- **Testability**: >70% code coverage, comprehensive test suite

## 🎪 Domain Model

### Core Aggregates
1. **User**: Central business entity with authentication capabilities

### Key Value Objects
1. **Email**: Validated email address with format validation
2. **Phone**: Validated phone number with international format
3. **JWTToken**: Immutable token value object

### Key Business Rules
- Email must be unique across all users
- Phone must be unique across all users
- Passwords must meet complexity requirements
- JWT tokens contain user claims (id, email, name, phone, created_at)
- User events tracked for audit trail

## 🔧 Technical Constraints

### Technology Stack
- **Language**: Go 1.23+
- **Database**: PostgreSQL
- **Authentication**: JWT (HS256), bcrypt for passwords
- **APIs**: REST (HTTP) and gRPC
- **Architecture**: Clean Architecture + DDD

### Integration Requirements
- **HTTP API**: REST endpoints for user-facing authentication
- **gRPC API**: Token validation service for other microservices
- **PostgreSQL**: Primary data store with user and event storage

## 🎯 Success Criteria

### Quality Metrics
- **Code Coverage**: >70%
- **Test Suite**: Comprehensive tests across all layers
- **Performance**: <50ms average token validation time
- **Reliability**: Zero data loss, transactional consistency

### Business Value
- **Developer Experience**: Clean, testable, maintainable codebase
- **System Integration**: gRPC API for seamless microservice integration
- **Security**: Production-ready authentication with JWT and bcrypt
- **Scalability**: Foundation for future multi-tenant authentication

## 🚀 Project Scope

### In Scope
- User registration with email, phone, name, password
- User login with email/password
- JWT token generation (access and refresh tokens)
- Token validation via gRPC
- Domain event system
- Comprehensive test coverage
- OpenAPI documentation

### Out of Scope (Current Version)
- Password reset functionality
- Multi-factor authentication (MFA)
- OAuth/Social login
- User profile management beyond authentication
- Email verification
- Phone verification

## 📊 Current Status

**Version**: 1.0.0  
**Status**: Production Ready ✅  
**Last Updated**: November 10, 2025

### Achievements
- ✅ Complete Clean Architecture implementation
- ✅ Domain-driven design with User aggregate
- ✅ CQRS pattern with command/query separation
- ✅ Event-driven architecture with domain events
- ✅ Comprehensive test suite
- ✅ JWT authentication with bcrypt password hashing
- ✅ gRPC service for token validation
- ✅ TransactionManager pattern for transactional consistency

### Quality Metrics
- **Architecture Compliance**: Clean Architecture + DDD
- **Security**: JWT authentication, bcrypt password hashing, input validation
- **Performance**: Optimized database queries
- **Documentation**: Comprehensive API and architecture docs

---

**This document serves as the foundation for all other Memory Bank files and defines the core scope and requirements of the Quest Auth project.**

