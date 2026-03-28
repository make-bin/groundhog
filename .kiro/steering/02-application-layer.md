# Application Layer Standards

## Purpose

The Application Layer orchestrates use cases and workflows. It coordinates domain objects, manages transactions, and transforms data between external interfaces and the domain layer.

## Location

```
pkg/application/
├── service/           # Application services (use cases)
├── dto/               # Data Transfer Objects
└── assembler/         # DTO ↔ Domain converters
```

## Core Principles

### 1. Use Case Orchestration
- Application services implement use cases
- Coordinate multiple domain objects
- Manage transaction boundaries
- NO business logic (delegate to domain)

### 2. Dependency Rule
**Application layer can depend on**:
- ✅ `pkg/domain/` - Domain models, repositories, services
- ✅ `pkg/utils/` - Technical utilities
- ❌ NO `pkg/infrastructure/` - Use interfaces from domain
- ❌ NO `pkg/interface/` - Application doesn't know about HTTP/gRPC

### 3. Data Transformation
- Use DTOs for external communication
- Use Assemblers for DTO ↔ Domain conversion
- Keep domain models pure

## Application Services

Application services implement use cases and orchestrate domain logic.

### Service Structure
```
pkg/application/service/
├── {aggregate_a}_app_service.go
├── {aggregate_b}_app_service.go
└── {feature}_app_service.go
```

### Application Service Rules
- ✅ Orchestrate use cases
- ✅ Manage transactions
- ✅ Convert DTOs to domain models
- ✅ Call domain services and repositories
- ✅ Handle logging and monitoring
- ❌ NO business logic (delegate to domain)
- ❌ NO direct database access
- ❌ NO HTTP/gRPC handling

### Typical Use Case Flow
1. Receive DTO from interface layer
2. Convert DTO to domain models/value objects
3. Validate using domain services
4. Execute business logic (domain methods)
5. Persist changes using repositories
6. Convert domain models to DTOs
7. Return DTOs to interface layer

### Service Responsibilities
- **Orchestration**: Coordinate domain objects
- **Transaction Management**: Begin, commit, rollback
- **DTO Conversion**: Transform between DTOs and domain models
- **Error Handling**: Convert domain errors to application errors
- **Logging**: Log use case execution
- **Validation**: Call domain validation services

## Data Transfer Objects (DTOs)

DTOs define the contract between application layer and external interfaces.

### DTO Structure
```
pkg/application/dto/
├── {aggregate_a}_dto.go
├── {aggregate_b}_dto.go
└── {feature}_dto.go
```

### DTO Rules
- ✅ Use for API contracts
- ✅ Include validation tags (binding, json)
- ✅ Use pointers for optional fields in update requests
- ✅ Separate request and response DTOs
- ❌ NO business logic
- ❌ NO domain model references
- ❌ NO database tags (gorm)

### DTO Types
- **Request DTOs**: Input from external interfaces
  - `CreateEntityRequest`
  - `UpdateEntityRequest`
  - `ListEntitiesRequest`
- **Response DTOs**: Output to external interfaces
  - `EntityResponse`
  - `ListEntitiesResponse`
- **Query DTOs**: Search and filter parameters
  - Pagination parameters
  - Filter criteria
  - Sort options

### DTO Validation Tags
- `json:"field_name"` - JSON serialization
- `binding:"required"` - Required field
- `binding:"email"` - Email validation
- `binding:"min=3,max=50"` - Length validation
- `binding:"omitempty"` - Optional field

## Assemblers

Assemblers convert between DTOs and domain models.

### Assembler Structure
```
pkg/application/assembler/
├── {aggregate_a}_assembler.go
├── {aggregate_b}_assembler.go
└── {feature}_assembler.go
```

### Assembler Rules
- ✅ Pure conversion functions
- ✅ Handle nil cases
- ✅ Extract values from value objects
- ✅ Convert domain types to primitive types
- ❌ NO business logic
- ❌ NO validation
- ❌ NO repository calls

### Assembler Patterns
- **Domain to DTO**: `ToEntityResponse(entity *model.Entity) *dto.EntityResponse`
- **DTO to Domain**: Convert in application service (not assembler)
- **List Conversion**: `ToEntityResponseList(entities []*model.Entity) []*dto.EntityResponse`
- **Pagination**: `ToListEntitiesResponse(entities, total, page, pageSize)`

## Transaction Management

### Transaction Rules
- ✅ Manage transactions in application layer
- ✅ Use defer for rollback
- ✅ Commit explicitly
- ❌ NO transactions in domain layer
- ❌ NO transactions in interface layer

### Transaction Pattern
1. Begin transaction
2. Execute operations (may span multiple repositories)
3. Commit on success
4. Rollback on error (via defer)

### Transaction Scope
- Single use case = single transaction
- Multiple repository operations in one transaction
- Domain events published after commit

## Error Handling

### Error Conversion
- Convert domain errors to application errors
- Add context to errors
- Log errors appropriately
- ❌ NO HTTP status codes in application layer
- ❌ NO error details exposure

### Error Types
- **Validation Error**: Invalid input
- **Not Found Error**: Entity not found
- **Conflict Error**: Duplicate entity
- **Forbidden Error**: Permission denied
- **Internal Error**: Unexpected error

## Testing Application Layer

### Unit Tests
- Mock repositories and domain services
- Test use case orchestration
- Test error handling
- Test DTO conversion
- ❌ NO real database
- ❌ NO HTTP testing (use interface layer tests)

### Test Focus
- Use case flow
- Transaction management
- Error handling
- DTO to domain conversion
- Repository interaction
- Domain service calls

### Mocking Strategy
- Mock repository interfaces
- Mock domain services
- Mock external dependencies
- Use testify/mock or similar

## Best Practices

### ✅ DO: Orchestrate domain logic
- Application service coordinates
- Domain layer implements rules
- Clear separation of concerns

### ❌ DON'T: Put business logic in application service
- Business logic belongs in domain
- Application service only orchestrates
- Validation in value objects

### ✅ DO: Use assemblers for conversion
- Dedicated conversion functions
- Consistent conversion logic
- Reusable across services

### ❌ DON'T: Manual conversion in service
- Avoid inline conversion
- Use assembler functions
- Maintain consistency

### ✅ DO: Handle transactions properly
- Begin at use case start
- Commit on success
- Rollback on error
- Use defer for cleanup

### ❌ DON'T: Leak transactions
- Always commit or rollback
- Use defer for safety
- Handle errors properly

## Summary

Application layer responsibilities:
- ✅ Orchestrate use cases
- ✅ Manage transactions
- ✅ Convert DTOs ↔ Domain models
- ✅ Coordinate domain services and repositories
- ✅ Handle application-level errors
- ❌ NO business logic
- ❌ NO direct database access
- ❌ NO HTTP/gRPC handling
