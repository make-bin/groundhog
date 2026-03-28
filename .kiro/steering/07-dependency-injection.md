---
inclusion: always
---

# Dependency Injection Standards

## Container-Based Dependency Injection

**Library**: `github.com/barnettZQG/inject` (runtime DI container)  
**Container Location**: `pkg/utils/container/container.go`  
**Server Initialization**: `pkg/server/server.go`

**See**: [Dependency Injection Container Implemetation Standards](./07.01-dependency-injection-container.md) for detailed documentation.

### Registration Order (STRICT)

Dependencies MUST be registered in this exact order:

1. **Configuration & Logger** - Load config, initialize logger
2. **Infrastructure Components** - Database, cache, external clients (use `ProvideWithName`)
3. **Domain Services** - Authentication, authorization services (use `Provides`)
4. **Application Services** - Use case implementations (use `Provides`)
5. **Interface Handlers** - HTTP/gRPC handlers (use `Provides`)
6. **Populate** - Call `container.Populate()` to wire all dependencies

### Injection Tags

```go
// Type-based injection (single implementation)
type userAppService struct {
    UserRepo repository.UserRepository `inject:""`
    AuthSvc  service.AuthenticationService `inject:""`
}

// Named injection (multiple instances or infrastructure)
type authService struct {
    DB     datastore.DataStore `inject:"datastore"`
    Cache  cache.Cache         `inject:"redis"`
    Logger logger.Logger       `inject:"logger"`
}
```

**Rules**:
- `inject:""` - Type-based injection for services with single implementation
- `inject:"name"` - Named injection for infrastructure with multiple instances
- Only use inject tags in service structs that RECEIVE dependencies
- Never use inject tags in infrastructure implementations

## Constructor Patterns

### Service Constructors (MUST Return Interfaces)

```go
// ✅ CORRECT - Returns interface type
func NewUserAppService() application.UserAppService {
    return &userAppService{}
}

// ❌ WRONG - Returns concrete type
func NewUserAppService() *userAppService {
    return &userAppService{}
}
```

### Infrastructure Constructors (Accept Config)

```go
// ✅ CORRECT - Accepts configuration parameters
func NewPostgresDataStore(cfg *config.DatabaseConfig) (datastore.DataStore, error) {
    // Initialize with config
    return &postgresDataStore{config: cfg}, nil
}
```

### Service Constructors (No Parameters)

```go
// ✅ CORRECT - No parameters, dependencies via inject tags
func NewAuthAppService() application.AuthAppService {
    return &authAppService{} // Dependencies injected by container
}
```

## Adding New Components

### New Application Service

1. Define service interface in `pkg/application/service/`
2. Implement with inject tags for dependencies
3. Constructor returns interface type (no parameters)
4. Register in container with `container.Provides(NewServiceName)`
5. Create DTOs in `pkg/application/dto/`
6. Create assemblers in `pkg/application/assembler/`

### New Domain Service

1. Define interface in `pkg/domain/service/`
2. Implement with inject tags for repository dependencies
3. Constructor returns interface type (no parameters)
4. Register in container with `container.Provides(NewServiceName)`

### New HTTP Handler

1. Create in `pkg/interface/http/handler/`
2. Use inject tags for application service dependencies
3. Constructor returns interface type (no parameters)
4. Register in container with `container.Provides(NewHandlerName)`
5. Register routes in `pkg/interface/http/router/router.go`

### New Infrastructure Component

1. Define interface in domain layer (if domain needs it) or utils layer
2. Implement in `pkg/infrastructure/`
3. Constructor accepts config parameters, returns interface
4. Register with `container.ProvideWithName("name", NewComponentFunc(cfg))`
5. Use named injection in dependent services: `inject:"name"`

## Common Mistakes

### ❌ Returning Concrete Types

```go
// WRONG
func NewUserService() *userService {
    return &userService{}
}

// CORRECT
func NewUserService() UserService {
    return &userService{}
}
```

### ❌ Using Inject Tags in Infrastructure

```go
// WRONG - Infrastructure implementations should not have inject tags
type postgresDataStore struct {
    Logger logger.Logger `inject:"logger"` // ❌ NO
}

// CORRECT - Pass dependencies via constructor or register with ProvideWithName
func NewPostgresDataStore(cfg *config.DatabaseConfig, logger logger.Logger) datastore.DataStore {
    return &postgresDataStore{logger: logger}
}
```

### ❌ Constructor with Parameters for Services

```go
// WRONG - Service constructors should have no parameters
func NewUserService(repo repository.UserRepository) UserService {
    return &userService{repo: repo}
}

// CORRECT - Use inject tags
type userService struct {
    Repo repository.UserRepository `inject:""`
}
func NewUserService() UserService {
    return &userService{}
}
```

## Initialization Flow

1. **Load Configuration** - `config.LoadConfig()`
2. **Initialize Logger** - `logger.NewLogger(cfg.Log)`
3. **Create Container** - `container.NewContainer() `
4. **Register Infrastructure** - Database, cache, clients (ProvideWithName)
5. **Register Domain Services** - Auth, authorization (Provides)
6. **Register Application Services** - Use cases (Provides)
7. **Register Handlers** - HTTP/gRPC handlers (Provides)
8. **Populate Container** - `container.Populate()` wires all dependencies
9. **Initialize Server** - Create Gin router, register routes
10. **Start Server** - Listen and serve
11. **Graceful Shutdown** - Handle signals, cleanup resources

## Key Principles

- **Constructor Returns Interface** - All service constructors return interface types
- **Infrastructure Accepts Config** - Infrastructure constructors accept configuration
- **Services Use Inject Tags** - Service dependencies injected via struct tags
- **Named Injection for Infrastructure** - Use ProvideWithName for multiple instances
- **Type Injection for Services** - Use Provides for single implementation services
- **No Container in Tests** - Unit tests create structs directly with mocks
