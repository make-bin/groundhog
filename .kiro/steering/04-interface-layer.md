# Interface Layer Standards

## Purpose

The Interface Layer handles external communication protocols (HTTP, gRPC, CLI). It transforms external requests into application service calls and formats responses.

## Location

```
pkg/interface/
├── http/              # HTTP handlers (Gin framework)
│   ├── handler/       # HTTP request handlers
│   ├── middleware/    # HTTP middleware
│   ├── router/        # Route definitions
│   └── response/      # Response helpers
├── grpc/              # gRPC service implementations
│   ├── server/        # gRPC server setup
│   └── handler/       # gRPC service handlers
└── cli/               # CLI commands
    └── cmd/           # Command definitions
```

## Core Principles

### 1. Protocol Handling
- Handle HTTP/gRPC/CLI protocols
- Transform requests to DTOs
- Format responses
- NO business logic

### 2. Dependency Rule
**Interface layer can depend on**:
- ✅ `pkg/application/` - Application services and DTOs
- ✅ `pkg/domain/` - Domain errors (for error handling)
- ✅ `pkg/utils/` - Technical utilities
- ❌ NO `pkg/infrastructure/` - Interface doesn't know about database

### 3. Separation of Concerns
- Protocol-specific code only
- Delegate to application services
- Handle authentication/authorization
- Format responses consistently

## HTTP Layer

### HTTP Handlers

HTTP handlers process HTTP requests using Gin framework.

#### Handler Structure
```
pkg/interface/http/handler/
├── {resource_a}_handler.go
├── {resource_b}_handler.go
└── {feature}_handler.go
```

#### Handler Rules
- ✅ Handle HTTP requests and responses
- ✅ Bind requests to DTOs
- ✅ Call application services
- ✅ Format responses consistently
- ✅ Handle errors and convert to HTTP status codes
- ❌ NO business logic
- ❌ NO direct repository access
- ❌ NO database operations

#### Handler Responsibilities
- **Request Binding**: Parse and validate HTTP requests
- **DTO Conversion**: Convert HTTP requests to DTOs
- **Service Invocation**: Call application services
- **Response Formatting**: Format responses consistently
- **Error Handling**: Convert errors to HTTP status codes
- **Logging**: Log HTTP requests and responses

#### HTTP Methods
- `GET` - Retrieve resources
- `POST` - Create resources
- `PUT` - Update resources (full)
- `PATCH` - Update resources (partial)
- `DELETE` - Delete resources

### Response Helpers

Standardized response formatting.

#### Response Structure
```
pkg/interface/http/response/
└── response.go
```

#### Response Format
```json
{
  "code": 200,
  "message": "Success",
  "data": {...},
  "error": "error message"
}
```

#### Response Methods
- `Success(c *gin.Context, code int, message string, data interface{})`
- `Error(c *gin.Context, code int, message string, err error)`
- `ValidationError(c *gin.Context, errors map[string]string)`

### Middleware

HTTP middleware for cross-cutting concerns.

#### Middleware Structure
```
pkg/interface/http/middleware/
├── {auth}_middleware.go
├── {cors}_middleware.go
├── {logger}_middleware.go
└── {rate_limit}_middleware.go
```

#### Middleware Types
- **Authentication**: JWT token validation
- **Authorization**: Role/permission checks
- **CORS**: Cross-origin resource sharing
- **Logging**: Request/response logging
- **Rate Limiting**: API rate limiting
- **Recovery**: Panic recovery

#### Middleware Rules
- ✅ Handle cross-cutting concerns
- ✅ Call `c.Next()` to continue chain
- ✅ Call `c.Abort()` to stop chain
- ✅ Set context values for downstream handlers
- ❌ NO business logic
- ❌ NO database access

### Router

Route definitions and setup.

#### Router Structure
```
pkg/interface/http/router/
└── router.go
```

#### Route Organization
- Group routes by resource
- Apply middleware to groups
- Version API routes (`/api/v1/`)
- Separate public and protected routes

#### Route Patterns
- `/api/v1/{resources}` - List resources (GET)
- `/api/v1/{resources}` - Create resource (POST)
- `/api/v1/{resources}/:id` - Get resource (GET)
- `/api/v1/{resources}/:id` - Update resource (PUT)
- `/api/v1/{resources}/:id` - Delete resource (DELETE)
- `/api/v1/{resources}/:id/{action}` - Resource action (POST)

## gRPC Layer

### gRPC Handlers

gRPC service implementations.

#### Handler Structure
```
pkg/interface/grpc/handler/
├── {resource_a}_grpc_handler.go
├── {resource_b}_grpc_handler.go
└── {feature}_grpc_handler.go
```

#### gRPC Handler Rules
- ✅ Implement protobuf service interfaces
- ✅ Convert protobuf messages to DTOs
- ✅ Call application services
- ✅ Convert DTOs to protobuf messages
- ✅ Handle errors and convert to gRPC status
- ❌ NO business logic
- ❌ NO direct repository access

#### gRPC Error Codes
- `codes.OK` - Success
- `codes.NotFound` - Resource not found
- `codes.AlreadyExists` - Duplicate resource
- `codes.InvalidArgument` - Invalid input
- `codes.PermissionDenied` - Permission denied
- `codes.Unauthenticated` - Authentication failed
- `codes.Internal` - Internal error

### Protocol Buffers

#### Proto Structure
```
proto/
├── {resource_a}/
│   └── v1/
│       └── {resource_a}.proto
├── {resource_b}/
│   └── v1/
│       └── {resource_b}.proto
└── common/
    └── v1/
        └── common.proto
```

#### Proto Rules
- Use proto3 syntax
- Version API definitions
- Define request/response messages
- Use common types for shared definitions

## CLI Layer

### CLI Commands

Command-line interface implementations.

#### Command Structure
```
pkg/interface/cli/cmd/
├── {resource_a}_cmd.go
├── {resource_b}_cmd.go
└── {operation}_cmd.go
```

#### CLI Rules
- ✅ Parse command-line arguments
- ✅ Call application services
- ✅ Format output for terminal
- ✅ Handle errors gracefully
- ❌ NO business logic
- ❌ NO direct repository access

## Error Handling

### HTTP Error Mapping
- `400 Bad Request` - Validation error
- `401 Unauthorized` - Authentication failed
- `403 Forbidden` - Permission denied
- `404 Not Found` - Resource not found
- `409 Conflict` - Duplicate resource
- `500 Internal Server Error` - Unexpected error

### Error Response Format
```json
{
  "code": 400,
  "message": "Validation failed",
  "error": "username is required"
}
```

## Testing Interface Layer

### HTTP Handler Tests
- Use `httptest` for testing
- Mock application services
- Test request binding
- Test response formatting
- Test error handling
- ❌ NO real database
- ❌ NO real external services

### Test Focus
- Request parsing
- DTO binding
- Service invocation
- Response formatting
- Error handling
- Middleware behavior

## Best Practices

### ✅ DO: Handle protocol-specific concerns
- HTTP/gRPC/CLI handling only
- Delegate to application services
- Format responses consistently

### ❌ DON'T: Put business logic in handlers
- Business logic in domain layer
- Orchestration in application layer
- Protocol handling in interface layer

### ✅ DO: Use consistent response format
- Standard response structure
- Consistent error format
- Clear success/error indication

### ❌ DON'T: Expose internal errors
- Convert to user-friendly messages
- Hide implementation details
- Log detailed errors internally

### ✅ DO: Validate requests early
- Use binding tags
- Validate before service call
- Return clear validation errors

### ❌ DON'T: Skip validation
- Always validate input
- Use DTO validation tags
- Fail fast on invalid input

## Summary

Interface layer responsibilities:
- ✅ Handle HTTP/gRPC/CLI protocols
- ✅ Transform requests to DTOs
- ✅ Call application services
- ✅ Format responses consistently
- ✅ Handle authentication/authorization
- ❌ NO business logic
- ❌ NO direct repository access
- ❌ NO database operations
