# Utils Layer Standards

## Purpose

The Utils Layer provides shared technical utilities that can be used across all layers. This layer has NO dependencies on other layers.

## Location

```
pkg/utils/
├── errors/           # Error handling utilities
├── bcode/            # Business error codes
├── logger/           # Logging utilities
├── config/           # Configuration management
├── validator/        # Validation utilities
├── crypto/           # Cryptography utilities
├── time/             # Time utilities
├── pprof/            # Performance profiling
└── container/        # Dependency injection container
```

## Core Principles

### 1. No Dependencies
**Utils layer MUST NOT depend on**:
- ❌ NO `pkg/domain/`
- ❌ NO `pkg/application/`
- ❌ NO `pkg/infrastructure/`
- ❌ NO `pkg/interface/`
- ✅ ONLY standard library and external utility packages

### 2. Reusability
- Generic, reusable functions
- No business logic
- No domain-specific code

### 3. Stateless
- Pure functions when possible
- Minimal state management
- Thread-safe implementations

## Error Utilities

### Error Types
- `ValidationError` - Input validation errors
- `NotFoundError` - Resource not found
- `ConflictError` - Resource conflict
- `UnauthorizedError` - Authentication failed
- `ForbiddenError` - Permission denied
- `InternalError` - Unexpected errors

### Error Structure
- Error type identifier
- Error message
- Wrapped error (optional)
- Error code (optional)

### Error Methods
- `Error() string` - Error message
- `Unwrap() error` - Wrapped error
- `Is(error) bool` - Error comparison

## Logger Utilities

### Logger Interface
- `Debug(msg string, keysAndValues ...interface{})`
- `Info(msg string, keysAndValues ...interface{})`
- `Warn(msg string, keysAndValues ...interface{})`
- `Error(msg string, keysAndValues ...interface{})`
- `Fatal(msg string, keysAndValues ...interface{})`
- `With(keysAndValues ...interface{}) Logger`

### Log Levels
- `DEBUG` - Detailed debugging information
- `INFO` - General information
- `WARN` - Warning messages
- `ERROR` - Error messages
- `FATAL` - Fatal errors (exit application)

### Logger Implementation
- Use structured logging (zap, logrus)
- Support context values
- Thread-safe
- Configurable output format

## Config Utilities

### Configuration Structure
- Server configuration
- Database configuration
- Redis configuration
- JWT configuration
- Log configuration

### Configuration Loading
- Load from file (YAML, JSON, TOML)
- Override with environment variables
- Validate configuration
- Provide defaults

### Configuration Rules
- ✅ Type-safe configuration
- ✅ Validation on load
- ✅ Environment variable support
- ✅ Default values
- ❌ NO business logic
- ❌ NO domain-specific config

## Validator Utilities

### Validation Interface
- `Validate(i interface{}) error`
- `ValidateVar(field interface{}, tag string) error`
- `RegisterValidation(tag string, fn ValidatorFunc) error`

### Validation Tags
- `required` - Required field
- `email` - Email format
- `min=n` - Minimum length/value
- `max=n` - Maximum length/value
- `len=n` - Exact length
- `oneof=a b c` - One of values
- `url` - URL format
- `uuid` - UUID format

### Custom Validations
- Register custom validation functions
- Reusable across application
- Clear error messages

## Crypto Utilities

### Hash Functions
- `HashSHA256(input string) string`
- `HashMD5(input string) string`
- `GenerateRandomString(length int) (string, error)`

### Encryption
- Symmetric encryption (AES)
- Asymmetric encryption (RSA)
- Key generation

### Password Hashing
- Use bcrypt for passwords
- Configurable cost factor
- Verify hashed passwords

## Time Utilities

### Time Functions
- `Now() time.Time` - Current UTC time
- `ParseISO8601(s string) (time.Time, error)`
- `FormatISO8601(t time.Time) string`
- `StartOfDay(t time.Time) time.Time`
- `EndOfDay(t time.Time) time.Time`
- `AddDays(t time.Time, days int) time.Time`

### Time Rules
- Always use UTC internally
- Convert to local time for display
- Use ISO8601 format for serialization

## Business Code Utilities (bcode)

### Purpose
Business error codes provide standardized error codes for API responses and error handling.

### Structure
```
pkg/utils/bcode/
├── bcode.go          # Business code definitions
└── codes.go          # Predefined error codes
```

### Business Code Components
- **Code**: Unique integer identifier
- **Message**: Human-readable error message
- **HTTPStatus**: Corresponding HTTP status code
- **Details**: Additional error context (optional)

### Business Code Interface

**Methods**:
- Code() - Returns unique integer identifier
- Message() - Returns human-readable error message
- HTTPStatus() - Returns corresponding HTTP status code
- WithDetails() - Adds additional error context

### Standard Business Codes
- **Success Codes** (0-999)
  - `0` - Success
  
- **Client Error Codes** (1000-1999)
  - `1000` - Invalid request
  - `1001` - Validation failed
  - `1002` - Missing required field
  - `1003` - Invalid format
  
- **Authentication Codes** (2000-2999)
  - `2000` - Unauthorized
  - `2001` - Invalid credentials
  - `2002` - Token expired
  - `2003` - Token invalid
  
- **Authorization Codes** (3000-3999)
  - `3000` - Forbidden
  - `3001` - Insufficient permissions
  - `3002` - Access denied
  
- **Resource Codes** (4000-4999)
  - `4000` - Resource not found
  - `4001` - Resource already exists
  - `4002` - Resource conflict
  
- **Server Error Codes** (5000-5999)
  - `5000` - Internal server error
  - `5001` - Database error
  - `5002` - External service error
  - `5003` - Configuration error

### Business Code Rules
- ✅ Use consistent code ranges
- ✅ Provide clear error messages
- ✅ Map to appropriate HTTP status codes
- ✅ Support internationalization (i18n)
- ❌ NO business logic in error codes
- ❌ NO sensitive information in messages

### Usage Pattern

**Define business codes**:
- CodeSuccess (0, "Success", 200)
- CodeInvalidRequest (1000, "Invalid request", 400)
- CodeUnauthorized (2000, "Unauthorized", 401)
- CodeResourceNotFound (4000, "Resource not found", 404)

**Use in error handling**:
- Return business code with optional details
- Add context information (resource type, ID, etc.)

## Performance Profiling (pprof)

### Purpose
Performance profiling utilities for debugging and optimization.

### Structure
```
pkg/utils/pprof/
├── pprof.go          # Profiling utilities
└── middleware.go     # HTTP middleware for pprof
```

### Profiling Types
- **CPU Profiling**: Analyze CPU usage
- **Memory Profiling**: Analyze memory allocation
- **Goroutine Profiling**: Analyze goroutine usage
- **Block Profiling**: Analyze blocking operations
- **Mutex Profiling**: Analyze mutex contention

### Profiling Interface

**Methods**:
- Start() - Start profiling
- Stop() - Stop profiling
- EnableHTTP(addr) - Enable HTTP profiling endpoints

### HTTP Endpoints
- `/debug/pprof/` - Profiling index
- `/debug/pprof/profile` - CPU profile
- `/debug/pprof/heap` - Memory heap profile
- `/debug/pprof/goroutine` - Goroutine profile
- `/debug/pprof/block` - Block profile
- `/debug/pprof/mutex` - Mutex profile
- `/debug/pprof/trace` - Execution trace

### Profiling Rules
- ✅ Enable only in development/staging
- ✅ Protect endpoints with authentication
- ✅ Use separate port for profiling
- ✅ Document profiling procedures
- ❌ NO profiling in production (unless necessary)
- ❌ NO public access to profiling endpoints

### Usage Pattern

**Enable pprof HTTP endpoints**:
- Import net/http/pprof package
- Start profiling server on separate port (e.g., localhost:6060)

**CPU profiling**:
- Create profile file
- Start CPU profile
- Defer stop CPU profile

**Memory profiling**:
- Create profile file
- Write heap profile
- Close file

### Profiling Best Practices
- Profile in realistic conditions
- Use benchmarks for comparison
- Focus on hotspots
- Profile before and after optimization
- Document profiling results

## Dependency Injection Container

### Purpose
Dependency injection container for managing application dependencies and lifecycle.

### Structure
```
pkg/utils/container/
├── container.go      # Container implementation
├── provider.go       # Provider interface
└── scope.go          # Dependency scopes
```

### Container Interface

**Methods**:
- Register - Register dependency provider
- Resolve - Resolve dependency by name
- ResolveAs - Resolve with type assertion
- Has - Check if dependency exists
- Remove - Remove dependency

### Provider Types
- **Singleton**: Single instance shared across application
- **Transient**: New instance on each resolution
- **Scoped**: Single instance per scope (e.g., per request)

### Dependency Scopes
- **Application Scope**: Entire application lifetime
- **Request Scope**: Single HTTP request lifetime
- **Custom Scope**: User-defined scope

### Container Features
- **Constructor Injection**: Inject dependencies via constructor
- **Property Injection**: Inject dependencies via properties
- **Method Injection**: Inject dependencies via methods
- **Lazy Loading**: Create instances only when needed
- **Circular Dependency Detection**: Detect and prevent circular dependencies

### Container Rules
- ✅ Register dependencies at startup
- ✅ Use interfaces for dependencies
- ✅ Resolve dependencies explicitly
- ✅ Handle dependency lifecycle
- ❌ NO service locator pattern
- ❌ NO global container access
- ❌ NO runtime registration (except for testing)

### Usage Pattern

**Create container**:
- Initialize new container instance

**Register dependencies**:
- RegisterSingleton for shared instances
- RegisterTransient for per-request instances
- Use factory functions to create dependencies

**Resolve dependencies**:
- Resolve by name with error handling
- ResolveAs with type assertion for type safety

### Container Best Practices
- Register all dependencies at startup
- Use constructor injection
- Avoid circular dependencies
- Use interfaces for loose coupling
- Document dependency graph
- Test with mock dependencies

### Integration with Wire

While the container provides runtime dependency injection, prefer compile-time dependency injection with Wire for production:

- **Container**: Development, testing, dynamic scenarios
- **Wire**: Production, static dependency graph, compile-time safety

**Wire for production**:
- Use wire.Build with provider functions
- Compile-time dependency graph validation

**Container for testing**:
- Register mock dependencies
- Flexible test setup

## Testing Utils Layer

### Unit Tests
- Test utility functions
- Test error handling
- Test edge cases
- No external dependencies
- Fast and isolated

### Test Focus
- Function correctness
- Error handling
- Edge cases
- Thread safety

## Best Practices

### ✅ DO: Keep utilities generic and reusable
- Generic functions
- No business logic
- No domain concepts

### ❌ DON'T: Put business logic in utils
- Business logic in domain layer
- Utils for technical concerns only
- No domain-specific code

### ✅ DO: Use interfaces for flexibility
- Define interfaces
- Multiple implementations
- Easy to test

### ❌ DON'T: Create dependencies
- Utils layer is independent
- No dependencies on other layers
- Only standard library and external utilities

### ✅ DO: Make utilities stateless when possible
- Pure functions
- No global state
- Thread-safe

### ❌ DON'T: Use global state
- Avoid global variables
- Pass dependencies explicitly
- Use dependency injection

## Summary

Utils layer responsibilities:
- ✅ Provide generic technical utilities
- ✅ Error handling utilities
- ✅ Business error codes (bcode)
- ✅ Logging interfaces and implementations
- ✅ Configuration management
- ✅ Validation utilities
- ✅ Cryptography utilities
- ✅ Time utilities
- ✅ Performance profiling (pprof)
- ✅ Dependency injection container
- ❌ NO business logic
- ❌ NO dependencies on other layers
- ❌ NO domain-specific code

**Key Components**:
- **bcode**: Standardized business error codes for API responses
- **pprof**: Performance profiling and debugging tools
- **container**: Runtime dependency injection for flexible dependency management
