# Infrastructure Layer Standards

## Purpose

The Infrastructure Layer provides technical implementations for domain interfaces. It handles database access, external services, caching, and other infrastructure concerns.

## Location

```
pkg/infrastructure/
├── persistence/       # Repository implementations
│   ├── po/           # Persistence Objects (GORM models)
│   ├── mapper/       # Domain ↔ PO mappers
│   └── *_repository_impl.go
├── datastore/        # Database connections
├── migration/        # Database migration management
├── cache/            # Cache implementations
├── clients/          # External service clients
├── transaction/      # Transaction management
└── service/          # Infrastructure services
```

## Core Principles

### 1. Implementation Details
- Infrastructure layer implements domain interfaces
- Contains all technical concerns
- Isolated from business logic

### 2. Dependency Rule
**Infrastructure layer can depend on**:
- ✅ `pkg/domain/` - Implements domain interfaces
- ✅ `pkg/utils/` - Technical utilities
- ✅ External frameworks (GORM, Redis, etc.)
- ❌ NO `pkg/application/` - Infrastructure doesn't know about use cases
- ❌ NO `pkg/interface/` - Infrastructure doesn't know about HTTP/gRPC

### 3. Separation of Concerns
- Persistence Objects (PO) separate from Domain Models
- Mappers handle conversion
- Repository implementations isolated

## Persistence Objects (PO)

Persistence Objects are database-specific models with ORM tags.

### PO Structure
```
pkg/infrastructure/persistence/po/
├── {aggregate_a}_po.go
├── {aggregate_b}_po.go
└── {aggregate_c}_po.go
```

### PO Rules
- ✅ Use GORM tags for database mapping
- ✅ Include all database-specific fields (ID, timestamps, soft delete)
- ✅ Define relationships with GORM associations
- ✅ Implement TableName() method
- ❌ NO business logic
- ❌ NO validation (use domain models)
- ❌ NO exported methods except TableName()

### PO Components
- **Primary Key**: Auto-increment ID
- **Business ID**: UUID or unique identifier
- **Fields**: Database columns with GORM tags
- **Timestamps**: CreatedAt, UpdatedAt, DeletedAt
- **Relationships**: GORM associations (has many, belongs to, many2many)

### GORM Tags
- `gorm:"primaryKey"` - Primary key
- `gorm:"uniqueIndex"` - Unique constraint
- `gorm:"type:varchar(50)"` - Column type
- `gorm:"not null"` - NOT NULL constraint
- `gorm:"default:value"` - Default value
- `gorm:"index"` - Index
- `gorm:"many2many:table_name"` - Many-to-many relationship

## Mappers

Mappers convert between Domain Models and Persistence Objects.

### Mapper Structure
```
pkg/infrastructure/persistence/mapper/
├── {aggregate_a}_mapper.go
├── {aggregate_b}_mapper.go
└── {aggregate_c}_mapper.go
```

### Mapper Rules
- ✅ Pure conversion functions
- ✅ Handle nil cases
- ✅ Validate during PO → Domain conversion
- ✅ Handle nested objects (relationships)
- ❌ NO business logic
- ❌ NO repository calls
- ❌ NO database operations

### Mapper Methods
- `DomainToPO(entity *model.Entity) *po.EntityPO`
- `POToDomain(entityPO *po.EntityPO) (*model.Entity, error)`
- `POListToDomainList(entityPOs []*po.EntityPO) ([]*model.Entity, error)`

### Conversion Rules
- **Domain → PO**: Extract values from value objects
- **PO → Domain**: Reconstruct value objects with validation
- **Relationships**: Convert nested objects recursively
- **Nil Handling**: Return nil for nil input

## Repository Implementations

Repository implementations provide data access using GORM for AGGREGATES ONLY.

### Repository Structure
```
pkg/infrastructure/persistence/
├── {aggregate_a}_repository_impl.go  # Implementation for Aggregate A
├── {aggregate_b}_repository_impl.go  # Implementation for Aggregate B
└── {aggregate_c}_repository_impl.go  # Implementation for Aggregate C
```

### Repository Implementation Rules - CRITICAL

**MUST FOLLOW**:
- ✅ Implement domain repository interfaces (one per aggregate)
- ✅ Use GORM for database operations
- ✅ Use mappers for Domain ↔ PO conversion
- ✅ Handle GORM errors and convert to domain errors
- ✅ Use Preload/Joins to load entire aggregate (including internal entities)
- ✅ Use WithContext for cancellation and tracing
- ✅ Return complete aggregate with all internal entities loaded
- ✅ Save entire aggregate (root + all internal entities) in one operation

**STRICTLY FORBIDDEN**:
- ❌ NO business logic in repository implementations
- ❌ NO direct domain model exposure to GORM (use PO + mappers)
- ❌ NO repositories for entities (only for aggregates)
- ❌ NO partial aggregate loading (must load complete aggregate)
- ❌ NO transaction management in repository (belongs in application layer)

### Standard Repository Operations

**Create Operation** (persist entire aggregate):
- Convert aggregate to PO (including all internal entities)
- Save aggregate root and all internal entities using GORM Create
- Convert GORM errors to domain errors

**FindByID Operation** (load complete aggregate):
- Load aggregate root WITH all internal entities using Preload
- Use WithContext for cancellation support
- Convert GORM errors to domain errors (ErrRecordNotFound → ErrUserNotFound)
- Convert PO to complete aggregate using mapper

**Update Operation** (save entire aggregate):
- Convert aggregate to PO
- Update aggregate root using GORM Updates
- Update internal entities (delete old, create new relationships)
- Convert GORM errors to domain errors

**List Operation** (return complete aggregates):
- Build query with type-safe filters
- Count total matching records
- Load aggregates with internal entities using Preload
- Apply pagination (Offset, Limit)
- Convert all POs to aggregates using mapper

### GORM Patterns for Aggregate Loading

**Load Complete Aggregate**:
- ✅ CORRECT: Use Preload for all internal entities and related data
- ❌ WRONG: Load without Preload (missing internal entities)

**Save Complete Aggregate**:
- ✅ CORRECT: Use FullSaveAssociations or manually manage relationships
- Save aggregate root and internal entities together

**Common GORM Operations**:
- WithContext(ctx) - Context support for cancellation
- Preload("Relation") - Eager load relationships
- Joins("Relation") - Join tables for filtering
- Where("field = ?", value) - Query conditions
- Create(&po) - Insert aggregate
- Updates(map) - Update fields
- Save(&po) - Save with associations
- Delete(&po) - Soft delete
- Offset/Limit - Pagination

### Error Handling

**Convert GORM errors to domain errors**:
- gorm.ErrRecordNotFound → domain.ErrEntityNotFound
- gorm.ErrDuplicatedKey → domain.ErrEntityAlreadyExists
- gorm.ErrInvalidData → domain.ErrInvalidData
- Database connection errors → domain.ErrInternalError
- Log all unexpected errors for debugging

### Repository Implementation Anti-Patterns

**Anti-Pattern 1: Entity Repository Implementation**
- ❌ WRONG: Repository for entity (UserRoleRepositoryImpl)
- ✅ CORRECT: Manage through aggregate, save entire aggregate including internal entities

**Anti-Pattern 2: Partial Aggregate Loading**
- ❌ WRONG: Loading aggregate without Preload (missing internal entities)
- ✅ CORRECT: Load complete aggregate with all internal entities using Preload

**Anti-Pattern 3: Transaction Management in Repository**
- ❌ WRONG: Transaction management in repository (Begin, Commit, Rollback)
- ✅ CORRECT: Transaction in application layer, repository uses transaction from context

## Database Connection

### DataStore Abstraction

The infrastructure layer uses a unified DataStore abstraction for database operations. This provides database-agnostic repository implementations with consistent transaction management.

**See**: [Datastore Layer Standards](./03.01-infrastructure-layer-datastore.md) for detailed documentation on:
- DataStore interface and operations
- Transaction management with context
- Query options (filtering, pagination, sorting)
- Multiple database implementations (PostgreSQL, OpenGauss, in-memory)
- Repository integration patterns
- Error handling and conversion

**Key Points**:
- All repository implementations use DataStore interface
- Context-based transaction propagation
- Support for multiple database backends
- Consistent CRUD operations across implementations

## Database Migration Management

### Migration Infrastructure

The infrastructure layer provides database migration management using golang-migrate. This ensures schema versioning, automated migrations on startup, and safe schema evolution.

**See**: [Migration Layer Standards](./03.02-infrastructure-layer-migration.md) for detailed documentation on:
- MigrationManager interface and implementation
- Migration configuration and validation
- Startup integration with RunMigrations
- Migration sources (filesystem, embedded)
- Dirty state detection and recovery
- Password sanitization in logs
- Migration file conventions

**Key Points**:
- Migrations run automatically on application startup
- Configuration-driven execution (can be disabled)
- Fail fast on errors with recovery instructions
- Support for filesystem and embedded migration sources
- Retry logic for connection failures
- Structured logging with sanitized credentials

## Transaction Management

### Transaction Structure
```
pkg/infrastructure/transaction/
└── transaction_manager.go
```

### Transaction Interface
- `Begin(ctx context.Context) (Transaction, error)`
- `Commit() error`
- `Rollback() error`
- `GetDB() *gorm.DB`

### Transaction Rules
- Application layer manages transaction lifecycle
- Infrastructure provides transaction primitives
- Use context for cancellation
- Always rollback on error

## Cache Implementation

### Cache Structure
```
pkg/infrastructure/cache/
├── redis_cache.go
└── memory_cache.go
```

### Cache Interface
- `Get(ctx context.Context, key string, dest interface{}) error`
- `Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error`
- `Delete(ctx context.Context, key string) error`
- `Exists(ctx context.Context, key string) (bool, error)`

### Cache Patterns
- **Cache-Aside**: Check cache, query database, update cache
- **Write-Through**: Update database and cache together
- **TTL**: Set appropriate expiration time
- **Invalidation**: Delete cache on updates

## External Service Clients

### Client Structure
```
pkg/infrastructure/clients/
├── {external_service_a}_client.go
├── {external_service_b}_client.go
└── {external_service_c}_client.go
```

### Client Rules
- Implement domain interfaces
- Handle external API calls
- Convert external errors to domain errors
- Implement retry logic
- Handle timeouts

## Testing Infrastructure Layer

### Integration Tests
- Use real database for testing
- Test repository implementations
- Test mapper conversions
- Clean up test data
- ❌ NO mocks for database
- ❌ NO business logic tests

### Test Database
- Use test database or test containers
- Run migrations before tests
- Clean up after tests
- Isolate test data

## Best Practices

### ✅ DO: Separate PO from domain models
- PO with GORM tags in infrastructure
- Pure domain models in domain layer
- Mappers for conversion

### ❌ DON'T: Put business logic in repositories
- Repository only handles data access
- Business logic in domain layer
- Validation in domain models

### ✅ DO: Use mappers for conversion
- Dedicated mapper functions
- Handle nil cases
- Validate during conversion

### ❌ DON'T: Expose PO to other layers
- PO stays in infrastructure layer
- Return domain models from repositories
- Use mappers for conversion

### ✅ DO: Handle GORM errors properly
- Convert to domain errors
- Log errors
- Provide meaningful error messages

### ❌ DON'T: Leak database errors
- Don't expose GORM errors
- Convert to domain errors
- Hide implementation details

## Summary

Infrastructure layer responsibilities:
- ✅ Implement domain repository interfaces (one per aggregate ONLY)
- ✅ Handle database operations with GORM
- ✅ Manage database migrations on startup
- ✅ Load complete aggregates (root + all internal entities)
- ✅ Save entire aggregates in one operation
- ✅ Manage external service integrations
- ✅ Provide caching implementations
- ✅ Convert domain models ↔ persistence objects using mappers
- ✅ Convert GORM errors to domain errors
- ❌ NO business logic
- ❌ NO direct domain model exposure to ORM
- ❌ NO repositories for entities (only for aggregates)
- ❌ NO partial aggregate loading
- ❌ NO transaction management in repositories

**Critical Rules**:
- Repository implementations ONLY for aggregate roots
- Always load complete aggregate with all internal entities
- Use Preload/Joins to load relationships
- Save entire aggregate (root + entities) together
- Convert all GORM errors to domain errors
- No transaction management in repository layer
- Run migrations on startup before database connection
- Fail fast on migration errors with recovery instructions
