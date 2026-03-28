# Domain Layer Standards

## Purpose

The Domain Layer is the core of the application. It contains business logic, domain models, and business rules. This layer must be completely independent of external frameworks and infrastructure.

## Location

```
pkg/domain/
├── aggregate/         # Aggregates (aggregate roots)
│   ├── {aggregate_a}/
│   │   ├── {aggregate_a}.go
│   │   └── {aggregate_a}_test.go
│   ├── {aggregate_b}/
│   │   ├── {aggregate_b}.go
│   │   └── {aggregate_b}_test.go
│   └── {aggregate_c}/
│       ├── {aggregate_c}.go
│       └── {aggregate_c}_test.go
├── entity/            # Entities (non-root entities)
│   ├── {entity_a}.go
│   ├── {entity_b}.go
│   ├── {entity_c}.go
│   └── {entity_a}_test.go
├── vo/                # Value Objects
│   ├── {id_type}.go          # ID value objects
│   ├── {validated_field}.go  # Validated field value objects
│   ├── {status_type}.go      # Status/enum value objects
│   └── {id_type}_test.go
├── repository/        # Repository interfaces
│   ├── {aggregate_a}_repository.go
│   ├── {aggregate_b}_repository.go
│   └── {aggregate_c}_repository.go
├── service/           # Domain services
│   ├── {domain_service_a}.go
│   └── {domain_service_b}.go
├── event/             # Domain events
│   ├── {aggregate_a}_events.go
│   └── {aggregate_b}_events.go
└── errors.go          # Domain errors
```

## Core Principles

### 1. Business Logic Encapsulation
- All business rules MUST be in domain layer
- Domain models enforce invariants
- Value objects ensure validity

### 2. Dependency Rule
**Domain layer can ONLY depend on**:
- ✅ `pkg/utils/` - Technical utilities (errors, logger interfaces)
- ❌ NO `pkg/application/`
- ❌ NO `pkg/infrastructure/`
- ❌ NO `pkg/interface/`
- ❌ NO external frameworks (GORM, Gin, etc.)

### 3. Pure Domain Models
- No ORM tags (no `gorm:""`)
- No JSON tags (no `json:""`)
- No HTTP dependencies
- No database dependencies

## Domain Models

### Directory Structure Rules

#### Aggregates (`pkg/domain/aggregate/`)
- Each aggregate in its own subdirectory
- Aggregate root file named after the aggregate
- Include aggregate-specific tests
- **Example**: `pkg/domain/aggregate/user/user.go`

#### Entities (`pkg/domain/entity/`)
- Non-root entities in flat structure
- One file per entity
- Include entity tests
- **Example**: `pkg/domain/entity/role.go`

#### Value Objects (`pkg/domain/vo/`)
- All value objects in flat structure
- One file per value object
- Include value object tests
- **Example**: `pkg/domain/vo/email.go`

### Aggregates

Aggregates are clusters of domain objects treated as a single unit with an aggregate root.

**Location**: `pkg/domain/aggregate/{aggregate_name}/{aggregate_name}.go`

**Aggregate Rules**:
- ✅ Each aggregate in its own subdirectory
- ✅ Aggregate root controls access to internal entities
- ✅ Enforce consistency boundaries
- ✅ Use value objects for validated fields
- ✅ Provide getters for all fields
- ✅ Enforce invariants in methods
- ✅ Update `updatedAt` when state changes
- ❌ NO public fields
- ❌ NO setters without validation
- ❌ NO ORM tags
- ❌ NO cross-aggregate references (use IDs only)

**Aggregate Root Responsibilities**:
- Control access to internal entities
- Maintain aggregate consistency
- Enforce business rules
- Publish domain events
- Manage lifecycle of internal entities

**Key Methods**:
- Constructor: `NewAggregate(...)` with required parameters
- Getters: `ID()`, `Name()`, etc.
- Business methods: `Activate()`, `Deactivate()`, `AddEntity()`, etc.
- Validation: Enforce business rules in methods

**Example Structure**:
```
pkg/domain/aggregate/{aggregate_name}/
├── {aggregate_name}.go           # Aggregate root
└── {aggregate_name}_test.go      # Aggregate tests
```

### Entities

Entities are objects with unique identity that are NOT aggregate roots.

**Location**: `pkg/domain/entity/{entity_name}.go`

**Entity Rules**:
- ✅ Objects with unique identity
- ✅ Encapsulate business logic
- ✅ Use value objects for validated fields
- ✅ Provide getters for all fields
- ✅ Enforce invariants in methods
- ❌ NO public fields
- ❌ NO setters without validation
- ❌ NO ORM tags
- ❌ NO direct persistence (accessed through aggregate)

**Entity Characteristics**:
- Part of an aggregate
- Accessed through aggregate root
- Identity within aggregate context
- Lifecycle managed by aggregate root

**Key Methods**:
- Constructor: `NewEntity(...)` with required parameters
- Getters: `ID()`, `Name()`, etc.
- Business methods: Entity-specific operations
- Validation: Enforce entity-level rules

**Example Structure**:
```
pkg/domain/entity/
├── {entity_a}.go           # Entity A
├── {entity_b}.go           # Entity B
├── {entity_c}.go           # Entity C
└── {entity_a}_test.go      # Entity tests
```

### Value Objects

Value objects are immutable objects defined by their attributes.

**Location**: `pkg/domain/vo/{value_object_name}.go`

**Value Object Rules**:
- ✅ Immutable (no setters)
- ✅ Validate in constructor
- ✅ Provide `Equals()` method
- ✅ Encapsulate validation logic
- ✅ One file per value object
- ❌ NO public fields
- ❌ NO mutation methods
- ❌ NO ORM tags

**Common Value Object Types**:
- **Identifiers**: Entity IDs (UUID, string, etc.)
- **Validated Strings**: Fields requiring format validation
- **Secure Values**: Sensitive data (hashed, encrypted)
- **Enumerations**: Status types, categories
- **Measurements**: Quantities, amounts, metrics
- **Ranges**: Date ranges, time ranges, numeric ranges

**Value Object Pattern**:
- Private field for value storage
- Constructor with validation (NewValueObject)
- Getter method (Value())
- Equals method for comparison
- Immutable after creation

**Example Structure**:
```
pkg/domain/vo/
├── {entity}_id.go        # ID value object
├── {field}_vo.go         # Validated field value object
├── {status}_type.go      # Status enumeration
└── {field}_vo_test.go    # Value object tests
```

### Aggregate vs Entity vs Value Object

**When to use Aggregate**:
- Root entity that controls a cluster of objects
- Needs to maintain consistency across multiple entities
- Has its own lifecycle and identity
- Represents a business concept with clear boundaries

**When to use Entity**:
- Has unique identity but is not a root
- Part of an aggregate
- Accessed through aggregate root
- Lifecycle managed by aggregate

**When to use Value Object**:
- Defined by its attributes, not identity
- Immutable
- Represents a concept or measurement
- No independent lifecycle

### Aggregate Boundaries

**Consistency Boundary**:
- Aggregate defines transaction boundary
- All changes within aggregate are consistent
- Use domain events for cross-aggregate consistency

**Reference Rules**:
- External objects reference aggregate by ID only
- No direct references to internal entities
- Use repository to load entire aggregate

**Transaction Rules**:
- One transaction per aggregate
- Multiple aggregates = eventual consistency
- Use domain events for coordination

## Repository Interfaces

Repository interfaces define data access contracts for aggregates ONLY.

### Structure
```
pkg/domain/repository/
├── {aggregate_a}_repository.go       # Repository for Aggregate A
├── {aggregate_b}_repository.go       # Repository for Aggregate B
└── {aggregate_c}_repository.go       # Repository for Aggregate C
```

### Repository Scope - CRITICAL RULES

**MUST FOLLOW**:
- ✅ **ONE repository per aggregate root ONLY**
- ✅ **Repository operates on ENTIRE aggregate**
- ✅ **Repository MUST return aggregate models, NEVER entities**
- ✅ **Entities are accessed ONLY through their aggregate root**

**STRICTLY FORBIDDEN**:
- ❌ **NO repositories for entities** - Entities are managed by their aggregate
- ❌ **NO repositories for value objects** - Value objects have no independent lifecycle
- ❌ **NO repositories for relationship tables** - Relationships are part of aggregate
- ❌ **NO repositories for join tables** - Many-to-many managed within aggregates

### Repository Rules

**Interface Definition**:
- ✅ Define interfaces in `pkg/domain/repository/`
- ✅ Use aggregate models as parameters and return types
- ✅ Use value objects for all ID parameters (NO primitive types)
- ✅ Use type-safe filter structs (NO generic maps)
- ✅ Use `context.Context` for cancellation and tracing
- ✅ Return domain errors (defined in `pkg/domain/errors.go`)

**Strictly Forbidden**:
- ❌ NO implementation details in interface
- ❌ NO database-specific types (GORM, SQL types)
- ❌ NO ORM references or tags
- ❌ NO transaction management methods (belongs in application layer)
- ❌ NO primitive types for IDs (use value objects)
- ❌ NO generic `map[string]interface{}` filters
- ❌ NO JSON/XML tags in domain layer
- ❌ NO returning entities directly (must be part of aggregate)

### Standard Repository Methods

**Basic CRUD Operations** (operate on aggregate):
- Create - Persists a new aggregate
- FindByID - Retrieves aggregate by ID with all internal entities loaded
- Update - Persists changes to entire aggregate
- Delete - Removes aggregate (soft delete recommended)

**Query Operations** (return aggregates):
- List - Retrieves aggregates with pagination and type-safe filtering
- FindByField - Queries aggregates by specific field (using value objects)

**Existence Checks** (use value objects):
- ExistsByField - Checks if aggregate exists by field

**Count Operations**:
- Count - Returns total number of aggregates matching filter

### Type-Safe Filter Structs

**MUST use type-safe filter structs, NOT generic maps**:

✅ CORRECT - Type-safe filter struct with domain types:
- Use domain types for status fields
- Use pointers for optional filters
- Include keyword for search
- Include time range filters

❌ WRONG - Generic map filter (map[string]interface{})

### Repository Method Naming Conventions

**Query Methods**:
- `FindByID` - Single aggregate by ID
- `FindBy{Field}` - Single aggregate by unique field
- `FindBy{Criteria}` - Multiple aggregates by criteria
- `List` - Paginated list with filtering

**Existence Checks**:
- `ExistsBy{Field}` - Check existence by field

**Aggregate Operations**:
- `Create` - Create new aggregate
- `Update` - Update entire aggregate
- `Delete` - Remove aggregate

### Managing Relationships Within Aggregates

**Entities within aggregates are managed through aggregate methods**:

✅ CORRECT - Manage entities through aggregate:
- Aggregate contains internal entities as private fields
- Add/remove entities through aggregate methods (e.g., AssignRole)
- Repository saves entire aggregate including internal entities

❌ WRONG - Separate repository for entity (e.g., UserRoleRepository)

### Cross-Aggregate References

**Use IDs only, never direct references**:

✅ CORRECT - Reference by ID (roleIDs []vo.RoleID)
❌ WRONG - Direct aggregate reference (roles []*role.Role)

### Repository Interface Examples

**User Aggregate Repository** includes:
- Basic CRUD (Create, FindByID, Update, Delete)
- Queries with value objects (FindByUsername, FindByEmail)
- List with type-safe filter
- Existence checks with value objects
- Count with type-safe filter

**Role Aggregate Repository** includes:
- Basic CRUD operations
- FindByName query
- List with type-safe filter
- ExistsByName check
- Count with type-safe filter

### What NOT to Create

**❌ DO NOT create these repositories** (violates DDD):

- Entity repository (UserRoleRepository, RolePermissionRepository)
- Value object repository (EmailRepository)
- Transaction methods in repository (WithTransaction, CommitTransaction)
- Repositories using primitive types instead of value objects
- Repositories using generic map filters

### Repository Anti-Patterns

**Anti-Pattern 1: Entity Repository**
- ❌ WRONG: Separate repository for entity
- ✅ CORRECT: Manage through aggregate methods, save entire aggregate

**Anti-Pattern 2: Primitive Type Parameters**
- ❌ WRONG: FindByID(ctx, id uint)
- ✅ CORRECT: FindByID(ctx, id vo.UserID)

**Anti-Pattern 3: Generic Filters**
- ❌ WRONG: List(ctx, filters map[string]interface{})
- ✅ CORRECT: List(ctx, filter UserFilter, offset, limit int)

**Anti-Pattern 4: Transaction Management**
- ❌ WRONG: Transaction methods in repository interface
- ✅ CORRECT: TransactionManager in application layer

## Domain Services

Domain services contain business logic that doesn't belong to a single entity.

### When to Use Domain Services
- Operations involving multiple entities
- Business logic that doesn't naturally fit in an entity
- Complex calculations or validations
- Cross-aggregate operations

### Domain Service Rules
- ✅ Stateless operations
- ✅ Coordinate multiple domain objects
- ✅ Implement complex business rules
- ✅ Use repository interfaces
- ❌ NO infrastructure dependencies
- ❌ NO transaction management (use application layer)

**Common Domain Service Types**:
- Validation services (cross-entity validation, uniqueness checks)
- Authorization services (access control, permission checks)
- Calculation services (complex business calculations)
- Policy services (business policy enforcement)
- Coordination services (multi-entity operations)

## Domain Events

Domain events represent significant business occurrences.

### Event Rules
- Immutable
- Past tense naming (UserCreated, OrderPlaced)
- Contain relevant data
- Enable loose coupling between aggregates

**Event Pattern**:
- DomainEvent interface with OccurredAt() and EventType() methods
- Concrete event structs with entityID and occurredAt fields
- Immutable event data

## Domain Errors

Domain-specific errors defined in `pkg/domain/errors.go`.

### Error Categories
- **Not Found**: Entity not found
- **Already Exists**: Duplicate entity
- **Invalid State**: Invalid state transition
- **Validation**: Business rule violation
- **Permission**: Authorization failure

### Error Naming Convention
- `Err{Entity}NotFound`
- `Err{Entity}AlreadyExists`
- `ErrInvalid{Entity}State`
- `Err{Field}TooShort`
- `Err{Field}TooLong`
- `ErrInvalid{Field}Format`

## Testing Domain Layer

### Unit Tests

**Test Location**:
- Aggregate tests: `pkg/domain/aggregate/{aggregate_name}/{aggregate_name}_test.go`
- Entity tests: `pkg/domain/entity/{entity_name}_test.go`
- Value object tests: `pkg/domain/vo/{vo_name}_test.go`

**Test Rules**:
- ✅ Test business logic
- ✅ Test value object validation
- ✅ Test entity invariants
- ✅ Test aggregate consistency
- ✅ No external dependencies
- ❌ NO database
- ❌ NO HTTP
- ❌ NO mocks (domain layer is pure)

**Test Focus**:
- **Value Objects**: Validation, immutability, equality
- **Entities**: Business methods, state transitions
- **Aggregates**: Consistency, invariants, business rules
- **Domain Services**: Complex business logic

## Best Practices

### ✅ DO: Encapsulate business logic in domain models
- Business rules in entity methods
- Validation in value objects
- Complex logic in domain services

### ❌ DON'T: Put business logic outside domain
- Application layer orchestrates, doesn't implement rules
- Infrastructure layer implements technical concerns
- Interface layer handles protocols

### ✅ DO: Create repositories ONLY for aggregates
- One repository per aggregate root
- Repository returns complete aggregate with all internal entities
- Entities managed through aggregate methods
- No separate repositories for entities or relationships

### ❌ DON'T: Create repositories for entities
- Entities are part of aggregates
- Access entities through aggregate root
- Save entities by saving their aggregate
- No direct entity persistence

### ✅ DO: Use value objects for validation
- Validate in constructor
- Immutable after creation
- Type-safe domain concepts

### ❌ DON'T: Use primitive types without validation
- Wrap primitives in value objects
- Enforce business rules at type level

### ✅ DO: Make entities immutable where possible
- Private fields
- Getters only
- Business methods for state changes

### ❌ DON'T: Expose internal state
- No public fields
- No setters without validation
- Control state changes through methods

## File Organization Best Practices

### Naming Conventions
- **Aggregates**: `pkg/domain/aggregate/{aggregate_name}/{aggregate_name}.go`
- **Entities**: `pkg/domain/entity/{entity_name}.go`
- **Value Objects**: `pkg/domain/vo/{vo_name}.go`
- **Tests**: Same name with `_test.go` suffix

### File Naming Rules
- Use snake_case for file names
- Use singular form (user.go, not users.go)
- Match file name with primary type name
- One primary type per file

### Package Organization
- Aggregates: Each in its own package
- Entities: All in `entity` package
- Value Objects: All in `vo` package
- Avoid circular dependencies

### Import Paths
- Aggregates: `import "project/pkg/domain/aggregate/user"`
- Entities: `import "project/pkg/domain/entity"`
- Value Objects: `import "project/pkg/domain/vo"`

## Summary

Domain layer responsibilities:
- ✅ Encapsulate business logic
- ✅ Define aggregates (in separate subdirectories)
- ✅ Define entities (in flat structure)
- ✅ Define value objects (in flat structure)
- ✅ Define repository interfaces (one per aggregate ONLY)
- ✅ Implement domain services
- ✅ Define domain errors and events
- ❌ NO infrastructure dependencies
- ❌ NO framework dependencies
- ❌ NO ORM tags
- ❌ NO JSON tags
- ❌ NO repositories for entities
- ❌ NO repositories for value objects
- ❌ NO repositories for relationship tables

**Key Principles**:
- **Aggregates control consistency boundaries**
- **Entities accessed ONLY through aggregates**
- **Value objects ensure validity**
- **Repositories operate ONLY on aggregates, NEVER on entities**
- **Repository methods MUST return aggregate models**
- **Domain services coordinate complex logic**
- **One repository per aggregate root - strictly enforced**

**Repository Critical Rules**:
- ✅ Repository interfaces ONLY for aggregate roots
- ✅ Repository methods return aggregate models (with all internal entities loaded)
- ✅ Use value objects for all parameters (NO primitive types)
- ✅ Use type-safe filter structs (NO generic maps)
- ❌ NO entity repositories
- ❌ NO transaction methods in repository interfaces
- ❌ NO primitive type parameters (uint, string for IDs)
