---
inclusion: always
---

# DDD Architecture Overview

## Purpose

This document defines the strict Domain-Driven Design (DDD) architecture standards for the project. All code MUST follow these architectural principles and layer boundaries.

## Architecture Principles

### 1. Layered Architecture

The project follows a strict 5-layer DDD architecture:

```
┌─────────────────────────────────────────┐
│         Interface Layer                 │  ← User-facing interfaces (HTTP, gRPC, CLI)
│         pkg/interface/                  │
└─────────────────────────────────────────┘
              ↓ depends on
┌─────────────────────────────────────────┐
│        Application Layer                │  ← Use cases and orchestration
│        pkg/application/                 │
└─────────────────────────────────────────┘
              ↓ depends on
┌─────────────────────────────────────────┐
│          Domain Layer                   │  ← Business logic and rules
│          pkg/domain/                    │  ← CORE - No external dependencies
└─────────────────────────────────────────┘
              ↑ implemented by
┌─────────────────────────────────────────┐
│      Infrastructure Layer               │  ← Technical implementations
│      pkg/infrastructure/                │
└─────────────────────────────────────────┘
              ↑ uses
┌─────────────────────────────────────────┐
│          Utils Layer                    │  ← Shared utilities
│          pkg/utils/                     │
└─────────────────────────────────────────┘
```

**Layer Descriptions**:
- **Interface Layer**: Handles external protocols (HTTP, gRPC, CLI), transforms requests/responses
- **Application Layer**: Orchestrates use cases, manages transactions, coordinates domain objects
- **Domain Layer**: Contains business logic, domain models, business rules (CORE)
- **Infrastructure Layer**: Implements technical concerns (database, cache, external services)
- **Utils Layer**: Provides generic technical utilities (logging, config, validation)

### 2. Dependency Rule

**CRITICAL**: Dependencies MUST only point inward (toward the domain layer).

```
Interface Layer      → Application, Domain, Utils
Application Layer    → Domain, Utils
Domain Layer         → Utils ONLY
Infrastructure Layer → Domain, Utils
Utils Layer          → NOTHING (no dependencies)
```

### 3. Core Principles

#### Separation of Concerns
- Each layer has a single, well-defined responsibility
- Business logic MUST reside in the domain layer
- Technical concerns MUST reside in infrastructure layer

#### Dependency Inversion
- High-level modules do NOT depend on low-level modules
- Both depend on abstractions (interfaces)
- Abstractions are defined in the domain layer

#### Ubiquitous Language
- Use business terminology in code
- Domain models reflect business concepts
- Avoid technical jargon in domain layer

## Layer Responsibilities

### Interface Layer (pkg/interface/)
**Purpose**: Handle external communication protocols

**Responsibilities**:
- HTTP/REST endpoints (Gin handlers)
- gRPC service implementations
- CLI commands
- Request/response transformation
- Protocol-specific validation
- Authentication/authorization middleware

**MUST NOT**:
- Contain business logic
- Access database directly
- Implement domain rules

### Application Layer (pkg/application/)
**Purpose**: Orchestrate use cases and workflows

**Responsibilities**:
- Use case implementation
- Transaction management
- DTO (Data Transfer Object) definitions
- Assemblers (DTO ↔ Domain conversion)
- Application service coordination

**MUST NOT**:
- Contain business rules
- Access database directly
- Handle HTTP/gRPC protocols

### Domain Layer (pkg/domain/)
**Purpose**: Encapsulate business logic and rules

**Responsibilities**:
- Aggregates (aggregate roots in separate subdirectories)
- Entities (non-root entities in flat structure)
- Value Objects (immutable validated objects)
- Business rules and invariants
- Domain events
- Repository interfaces (one per aggregate)
- Domain services
- Domain errors

**MUST NOT**:
- Depend on external frameworks
- Contain infrastructure code
- Have database/HTTP dependencies
- Use ORM or JSON tags

### Infrastructure Layer (pkg/infrastructure/)
**Purpose**: Provide technical implementations

**Responsibilities**:
- Repository implementations
- Database access (GORM)
- Database migration management
- External service clients
- Caching implementations
- Transaction management
- Persistence objects (PO)
- Mappers (Domain ↔ PO)

**MUST NOT**:
- Contain business logic
- Be imported by domain layer

### Utils Layer (pkg/utils/)
**Purpose**: Provide shared technical utilities

**Responsibilities**:
- Error handling utilities
- Logging utilities
- Configuration management
- Validation helpers
- Common technical functions

**MUST NOT**:
- Contain business logic
- Depend on other layers

## Directory Structure

```
pkg/
├── interface/              # Interface Layer
│   ├── http/              # HTTP handlers (Gin)
│   ├── grpc/              # gRPC service implementations
│   └── cli/               # CLI commands
│
├── application/           # Application Layer
│   ├── service/           # Application services (use cases)
│   ├── dto/               # Data Transfer Objects
│   └── assembler/         # DTO ↔ Domain converters
│
├── domain/                # Domain Layer (CORE)
│   ├── aggregate/         # Aggregates (aggregate roots)
│   │   └── {aggregate_name}/
│   │       ├── {aggregate_name}.go
│   │       └── {aggregate_name}_test.go
│   ├── entity/            # Entities (non-root entities)
│   │   ├── {entity_name}.go
│   │   └── {entity_name}_test.go
│   ├── vo/                # Value Objects
│   │   ├── {vo_name}.go
│   │   └── {vo_name}_test.go
│   ├── repository/        # Repository interfaces
│   ├── service/           # Domain services
│   ├── event/             # Domain events
│   └── errors.go          # Domain errors
│
├── infrastructure/       # Infrastructure Layer
│   ├── persistence/      # Repository implementations
│   │   ├── po/           # Persistence Objects (GORM models)
│   │   ├── mapper/       # Domain ↔ PO mappers
│   │   └── *_repository_impl.go
│   ├── datastore/        # Database connections
│   ├── migration/        # Database migration management
│   ├── cache/            # Cache implementations
│   ├── clients/          # External service clients
│   ├── transaction/      # Transaction management
│   └── service/          # Infrastructure services
│
├── utils/                # Utils Layer
│   ├── errors/           # Error handling
│   ├── logger/           # Logging
│   ├── config/           # Configuration
│   └── validator/        # Validation utilities
│
└── server/                # Server bootstrap
    └── server.go          # Application entry point
```

## Architectural Patterns

### 1. Repository Pattern
- **Interfaces defined in `pkg/domain/repository/` - ONE per aggregate ONLY**
- **Implementations in `pkg/infrastructure/persistence/`**
- **Abstracts data access from business logic**
- **Repository operates on ENTIRE aggregate (root + all internal entities)**
- **NO repositories for entities - entities managed through aggregate**
- **Repository methods MUST return aggregate models, NEVER entities**

### 2. Aggregates
- Cluster of domain objects treated as a unit
- Each aggregate in its own subdirectory: `pkg/domain/aggregate/{aggregate_name}/`
- Aggregate root controls access to internal entities
- Enforce consistency boundaries

### 3. Entities
- Objects with unique identity (non-root)
- Defined in `pkg/domain/entity/`
- Part of an aggregate, accessed through aggregate root

### 4. Value Objects
- Immutable objects representing domain concepts
- Defined in `pkg/domain/vo/`
- Encapsulate validation logic

### 5. Domain Events
- Represent significant business occurrences
- Defined in `pkg/domain/event/`
- Enable loose coupling between aggregates

### 6. Domain Services
- Stateless operations that don't belong to entities
- Defined in `pkg/domain/service/`
- Coordinate multiple domain objects

## Repository Pattern - Strict DDD Rules

### Critical Repository Rules

**MUST FOLLOW - Non-Negotiable**:

1. **One Repository Per Aggregate Root ONLY**
   - ✅ Create repository for User aggregate
   - ✅ Create repository for Role aggregate
   - ❌ NO repository for UserRole entity
   - ❌ NO repository for RolePermission entity

2. **Repository Returns Complete Aggregates**
   - ✅ Repository methods return aggregate models with ALL internal entities loaded
   - ✅ Use GORM Preload to load relationships
   - ❌ NO partial aggregate loading
   - ❌ NO returning entities directly

3. **Use Value Objects for Parameters**
   - ✅ `FindByID(ctx context.Context, id vo.UserID)`
   - ❌ NO `FindByID(ctx context.Context, id uint)`
   - ✅ `FindByEmail(ctx context.Context, email vo.Email)`
   - ❌ NO `FindByEmail(ctx context.Context, email string)`

4. **Use Type-Safe Filter Structs**
   - ✅ `List(ctx context.Context, filter UserFilter, offset, limit int)`
   - ❌ NO `List(ctx context.Context, filters map[string]interface{})`

5. **No Transaction Management in Repository**
   - ✅ Transaction management in application layer
   - ❌ NO `WithTransaction()` in repository interface
   - ❌ NO `CommitTransaction()` in repository interface

6. **Entities Managed Through Aggregate**
   - ✅ `user.AssignRole(roleID)` then `userRepo.Update(ctx, user)`
   - ❌ NO `userRoleRepo.Create(ctx, userRole)`

### Repository Interface Checklist

Before creating a repository interface, verify:

- [ ] Is this for an aggregate root? (If NO, don't create repository)
- [ ] Does it use value objects for all ID parameters?
- [ ] Does it use type-safe filter structs?
- [ ] Does it return aggregate models (not entities)?
- [ ] Does it avoid transaction management methods?
- [ ] Does it use context.Context for all methods?

### Examples of Correct vs Incorrect

**✅ CORRECT Repository Interface**:
- Use value objects for all ID parameters (vo.UserID, vo.Email)
- Use type-safe filter structs (UserFilter)
- Return aggregate models with all internal entities
- No transaction management methods

**❌ INCORRECT Repository Interface**:
- Repository for entity (not aggregate) - UserRoleRepository
- Primitive types instead of value objects (uint, string)
- Generic map filters (map[string]interface{})
- Transaction management in repository (WithTransaction, CommitTransaction)

## Validation Rules

### Architecture Validation
Run `make arch-validate` or `make arch-validate-verbose` to ensure compliance.

### Validation Checks
- ✅ Domain layer only depends on utils
- ✅ No circular dependencies
- ✅ Dependency direction is correct
- ✅ No shared/common packages
- ✅ Repository interfaces in domain layer (one per aggregate ONLY)
- ✅ Repository implementations in infrastructure layer
- ✅ No repositories for entities
- ✅ Repository methods use value objects (not primitive types)
- ✅ Repository methods use type-safe filters (not generic maps)
- ✅ No transaction management in repository interfaces

## Anti-Patterns to Avoid

### ❌ DON'T: Create pkg/shared or pkg/common
- ❌ `pkg/shared/`
- ❌ `pkg/common/`
- ❌ `pkg/base/`

**Why**: Violates DDD principles, creates unclear boundaries

**Instead**: Use `pkg/domain/aggregate/`, `pkg/domain/entity/`, `pkg/domain/vo/` for domain concepts, `pkg/utils/` for technical utilities

### ❌ DON'T: Put business logic in application layer
**Wrong**: Business logic (validation, rules) in application service
**Correct**: Business logic in domain models and value objects

### ❌ DON'T: Access database from application layer
**Wrong**: Direct database access in application service
**Correct**: Use repository interfaces defined in domain layer

### ❌ DON'T: Put ORM tags in domain models
**Wrong**: GORM/JSON tags in domain models
**Correct**: Pure domain models + separate Persistence Objects (PO) in infrastructure layer

### ❌ DON'T: Create repositories for entities
**Wrong**: Separate repository for UserRole, RolePermission entities
**Correct**: Manage entities through their aggregate root, save entire aggregate

### ❌ DON'T: Use primitive types in repository interfaces
**Wrong**: `FindByID(ctx context.Context, id uint)` 
**Correct**: `FindByID(ctx context.Context, id vo.UserID)`

### ❌ DON'T: Use generic map filters
**Wrong**: `List(ctx context.Context, filters map[string]interface{})`
**Correct**: `List(ctx context.Context, filter UserFilter, offset, limit int)`

## Best Practices

### ✅ DO: Use value objects for validated fields
- Encapsulate validation logic in value object constructors
- Make value objects immutable
- Provide validation methods

### ✅ DO: Define repository interfaces in domain layer
- Repository interfaces belong to domain layer
- ONE repository per aggregate root ONLY
- Define contracts using aggregate models (not entities)
- Use value objects for all parameters (no primitive types)
- Use type-safe filter structs (no generic maps)
- Use context for cancellation
- Repository methods return complete aggregates with all internal entities

### ✅ DO: Implement repositories in infrastructure layer
- Repository implementations in infrastructure layer (one per aggregate)
- Use mappers to convert between aggregate models and persistence objects
- Load complete aggregate with all internal entities (use Preload/Joins)
- Save entire aggregate (root + entities) in one operation
- Convert GORM errors to domain errors
- Handle database-specific errors
- No transaction management in repository (belongs in application layer)

### ✅ DO: Use DTOs for API contracts
- DTOs for request/response in application layer
- Include validation tags (json, binding)
- Separate request and response DTOs

### ✅ DO: Use assemblers for DTO ↔ Domain conversion
- Assemblers convert between DTOs and domain models
- Pure conversion functions without business logic
- Handle nil cases properly

## Testing Strategy

### Unit Tests (Domain Layer)
- Test business logic without external dependencies
- No database, no HTTP, no external services
- Fast and isolated
- Focus on domain models, value objects, domain services

### Integration Tests (Infrastructure Layer)
- Test repository implementations with real database
- Test external service integrations
- Use test database or test containers
- Verify data persistence and retrieval

### E2E Tests (Interface Layer)
- Test complete workflows through HTTP/gRPC
- Test authentication and authorization
- Use test server
- Verify end-to-end functionality

## Migration Guide

When refactoring existing code to follow DDD:

1. **Identify domain concepts**: Extract business logic from services
2. **Create value objects**: Replace primitive types with validated objects
3. **Define repository interfaces**: Move to `pkg/domain/repository/`
4. **Separate persistence objects**: Create POs in `pkg/infrastructure/persistence/po/`
5. **Create mappers**: Implement Domain ↔ PO conversion
6. **Move business logic**: From application services to domain models
7. **Validate architecture**: Run `make arch-validate`

## References

- [Domain Layer Standards](./01-domain-layer.md)
- [Application Layer Standards](./02-application-layer.md)
- [Infrastructure Layer Standards](./03-infrastructure-layer.md)
  - [Datastore Layer Standards](./03.01-infrastructure-layer-datastore.md)
  - [Migration Layer Standards](./03.02-infrastructure-layer-migration.md)
- [Interface Layer Standards](./04-interface-layer.md)
- [Utils Layer Standards](./05-utils-layer.md)
- [Tech Stack Standards](./06-tech-stack.md)
- [Dependency Injection Standards](./07-dependency-injection.md)

## Summary

This architecture ensures:
- ✅ Clear separation of concerns
- ✅ Business logic isolated in domain layer
- ✅ Testable code without external dependencies
- ✅ Maintainable and scalable codebase
- ✅ Consistent coding standards across the project
