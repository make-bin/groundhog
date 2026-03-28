# Tech Stack Standards

## Purpose

This document defines the standard technology stack and frameworks used across the project. All code MUST follow these technology choices and usage patterns.

## Core Technologies

### Programming Language
- **Go 1.21+**
- Use Go modules for dependency management
- Follow official Go style guide
- Use `gofmt` for code formatting
- Use `golangci-lint` for linting

### Project Structure
```
project/
├── cmd/                    # Application entry points
│   └── server/
│       ├── main.go         # Application entry point
│       ├── container.go    # Dependency container setup
│       └── wire.go         # Wire setup (optional)
├── pkg/                    # Application code (DDD layers)
│   ├── domain/
│   ├── application/
│   ├── infrastructure/
│   ├── interface/
│   ├── utils/
│   └── server/             # Server initialization
│       └── server.go
├── proto/                  # Protocol buffer definitions
├── migrations/             # Database migrations
├── configs/                # Configuration files
├── scripts/                # Build and deployment scripts
├── docs/                   # Documentation
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Web Framework

### Gin (HTTP)
- **Version**: v1.9+
- **Usage**: HTTP API endpoints
- **Location**: `pkg/interface/http/`
- **Package**: `github.com/gin-gonic/gin`

#### Middleware
- CORS: Custom middleware
- Authentication: JWT middleware
- Logging: Custom logger middleware
- Recovery: `gin.Recovery()`
- Rate Limiting: Custom middleware

## gRPC Framework

### gRPC-Go
- **Version**: v1.58+
- **Usage**: gRPC service implementations
- **Location**: `pkg/interface/grpc/`
- **Package**: `google.golang.org/grpc`

#### Protocol Buffers
- **Version**: proto3
- **Location**: `proto/`
- **Generated code**: `proto/gen/`

## Database

### PostgreSQL
- **Version**: 14+
- **Driver**: `github.com/lib/pq`
- **ORM**: GORM v2

#### GORM
- **Version**: v1.25+
- **Usage**: Repository implementations ONLY
- **Location**: `pkg/infrastructure/persistence/`
- **Package**: `gorm.io/gorm`, `gorm.io/driver/postgres`

#### Migration
- **Tool**: golang-migrate
- **Location**: `migrations/`
- **Package**: `github.com/golang-migrate/migrate/v4`

#### Migration Files
```
migrations/
├── 000001_create_{table_a}.up.sql
├── 000001_create_{table_a}.down.sql
├── 000002_create_{table_b}.up.sql
└── 000002_create_{table_b}.down.sql
```

**Migration Naming Convention**:
- Format: `{version}_{description}.{direction}.sql`
- Version: Sequential number (000001, 000002, etc.)
- Description: Snake_case description
- Direction: `up` (apply) or `down` (rollback)

## Cache

### Redis
- **Version**: 7+
- **Client**: go-redis v9
- **Usage**: Caching, session storage
- **Location**: `pkg/infrastructure/cache/`
- **Package**: `github.com/redis/go-redis/v9`

## Authentication & Authorization

### JWT
- **Library**: golang-jwt v5
- **Usage**: Token generation and validation
- **Location**: `pkg/application/service/` (authentication service)
- **Package**: `github.com/golang-jwt/jwt/v5`

### Password Hashing
- **Library**: golang.org/x/crypto/bcrypt
- **Usage**: Password hashing
- **Location**: `pkg/domain/model/value_object.go`
- **Package**: `golang.org/x/crypto/bcrypt`

## Configuration Management

### Viper
- **Version**: v1.17+
- **Usage**: Configuration loading
- **Location**: `pkg/utils/config/`
- **Package**: `github.com/spf13/viper`

#### Configuration File Format
- **Format**: YAML
- **Location**: `configs/config.yaml`

#### Configuration Sections
- Server configuration (host, port, timeouts)
- Database configuration (connection, pool settings)
- Redis configuration (address, password, db)
- JWT configuration (secret, token TTL)
- Log configuration (level, format)

## Logging

### Zap
- **Version**: v1.26+
- **Usage**: Structured logging
- **Location**: `pkg/utils/logger/`
- **Package**: `go.uber.org/zap`

#### Log Levels
- DEBUG, INFO, WARN, ERROR, FATAL

#### Log Format
- JSON format for production
- Console format for development
- Structured key-value pairs

## Validation

### go-playground/validator
- **Version**: v10.15+
- **Usage**: Struct validation
- **Location**: `pkg/application/dto/`
- **Package**: `github.com/go-playground/validator/v10`

#### Validation Tags
- `required` - Required field
- `email` - Email validation
- `min=n, max=n` - Length validation
- `oneof=a b c` - Enum validation

## Testing

### Testing Framework
- **Standard library**: `testing`
- **Assertions**: testify
- **Mocking**: testify/mock
- **Package**: `github.com/stretchr/testify`

#### Test Structure
```
pkg/
├── domain/
│   ├── aggregate/
│   │   └── {aggregate}/
│   │       ├── {aggregate}.go
│   │       └── {aggregate}_test.go
│   ├── entity/
│   │   ├── {entity}.go
│   │   └── {entity}_test.go
│   └── vo/
│       ├── {vo}.go
│       └── {vo}_test.go
├── application/
│   └── service/
│       ├── {service}.go
│       └── {service}_test.go
└── infrastructure/
    └── persistence/
        ├── {repository}_impl.go
        └── {repository}_impl_test.go
```

#### Test Types
- **Unit Tests**: Domain layer (no external dependencies)
- **Unit Tests with Mocks**: Application layer (mock repositories)
- **Integration Tests**: Infrastructure layer (real database)
- **E2E Tests**: Interface layer (HTTP/gRPC)

#### Test Coverage
- **Target**: 80%+ coverage
- **Command**: `go test -cover ./...`
- **Report**: `go test -coverprofile=coverage.out ./...`

## API Documentation

### Swagger/OpenAPI
- **Library**: swaggo/swag
- **Version**: v1.16+
- **Usage**: API documentation generation
- **Package**: `github.com/swaggo/swag`, `github.com/swaggo/gin-swagger`

#### Swagger Annotations
- `@Summary` - Endpoint summary
- `@Description` - Detailed description
- `@Tags` - Group endpoints by resource
- `@Accept` - Request content type (json, xml)
- `@Produce` - Response content type (json, xml)
- `@Param` - Request parameters (path, query, body)
- `@Success` - Success response with status code
- `@Failure` - Error response with status code
- `@Router` - Route path and HTTP method
- `@Security` - Authentication requirements

## Dependency Injection

### Inject Library - REQUIRED
- **Package**: `github.com/barnettZQG/inject`
- **Version**: Latest
- **Usage**: Runtime dependency injection for ALL server initialization
- **Location**: `pkg/utils/container/container.go`
- **Documentation**: https://github.com/barnettZQG/inject

## Server Initialization

### Overview
Server initialization uses container-based dependency injection to manage all application dependencies and their lifecycle.

**See**: [Dependency Injection Standards](./07-dependency-injection.md) for detailed documentation.

### Key Components
- **Main Entry Point**: `cmd/server/main.go` - Application startup
- **Container Setup**: `pkg/utils/container/container.go` - Dependency registration
- **Server**: `pkg/server/server.go` - HTTP/gRPC server initialization

### Initialization Flow
1. Load configuration
2. Initialize logger
3. Create dependency container
4. Register all dependencies (infrastructure → domain → application → interface)
5. Initialize server
6. Start server
7. Graceful shutdown on signal

## Build & Deployment

### Makefile
Standard Makefile commands:
- `make build` - Build application
- `make test` - Run tests
- `make lint` - Run linter
- `make run` - Run application
- `make migrate-up` - Run migrations
- `make migrate-down` - Rollback migrations
- `make docker-build` - Build Docker image
- `make generate` - Generate code (swagger, mocks)
- `make clean` - Clean build artifacts

### Docker
- **Base image**: golang:1.21-alpine
- **Multi-stage build**: Yes
- **Location**: `Dockerfile`

#### Docker Stages
- Build stage: Compile Go application
- Run stage: Minimal runtime image

## Code Quality Tools

### Linting
- **Tool**: golangci-lint
- **Version**: v1.54+
- **Config**: `.golangci.yml`

#### Enabled Linters
- gofmt, golint, govet
- errcheck, staticcheck
- unused, gosimple
- structcheck, varcheck
- ineffassign, deadcode
- typecheck

### Formatting
- **Tool**: gofmt, goimports
- **Usage**: Automatic formatting
- **Command**: `gofmt -w .`, `goimports -w .`

## Version Control

### Git Workflow
- **Branching**: GitFlow
- **Main branches**: `main`, `develop`
- **Feature branches**: `feature/*`
- **Release branches**: `release/*`
- **Hotfix branches**: `hotfix/*`

### Commit Message Format
```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types**:
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation
- `style` - Formatting
- `refactor` - Code refactoring
- `test` - Tests
- `chore` - Maintenance

## CI/CD

### GitHub Actions
- **Location**: `.github/workflows/`
- **Workflows**: test, lint, build, deploy

#### Workflow Triggers
- Push to main/develop
- Pull requests
- Manual trigger

#### Workflow Steps
- Checkout code
- Set up Go
- Run tests
- Run linter
- Build application
- Build Docker image
- Deploy (if applicable)

## Environment Variables

### Required Variables

#### Server Configuration
- `SERVER_HOST` - Server host address
- `SERVER_PORT` - Server port number
- `SERVER_READ_TIMEOUT` - Read timeout
- `SERVER_WRITE_TIMEOUT` - Write timeout

#### Database Configuration
- `DB_HOST` - Database host
- `DB_PORT` - Database port
- `DB_USER` - Database user
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name
- `DB_SSLMODE` - SSL mode

#### Cache Configuration
- `REDIS_ADDR` - Redis address
- `REDIS_PASSWORD` - Redis password
- `REDIS_DB` - Redis database number

#### Authentication Configuration
- `JWT_SECRET` - JWT secret key
- `JWT_ACCESS_TOKEN_TTL` - Access token TTL
- `JWT_REFRESH_TOKEN_TTL` - Refresh token TTL

#### Logging Configuration
- `LOG_LEVEL` - Log level (debug, info, warn, error)
- `LOG_FORMAT` - Log format (json, console)

## Performance Guidelines

### Database
- Use connection pooling
- Set appropriate `max_idle_conns` and `max_open_conns`
- Use prepared statements
- Add indexes for frequently queried columns
- Use pagination for large result sets

### Caching
- Cache frequently accessed data
- Set appropriate TTL
- Use cache-aside pattern
- Invalidate cache on updates

### API
- Use pagination for list endpoints
- Implement rate limiting
- Use compression (gzip)
- Set appropriate timeouts

## Security Guidelines

### Authentication
- Use JWT for stateless authentication
- Store tokens securely
- Implement token refresh mechanism
- Use HTTPS in production

### Password
- Hash passwords with bcrypt
- Enforce password complexity
- Implement password reset flow
- Never log passwords

### API Security
- Validate all inputs
- Sanitize outputs
- Implement CORS properly
- Use rate limiting
- Log security events

## Monitoring & Observability

### Metrics
- **Tool**: Prometheus (optional)
- **Metrics**: Request count, latency, error rate
- **Endpoint**: `/metrics`

### Tracing
- **Tool**: OpenTelemetry (optional)
- **Usage**: Distributed tracing

### Health Checks
- **Endpoint**: `/health`
- **Response**: `{"status": "ok"}`
- **Checks**: Database connection, Redis connection



## Summary

This tech stack provides:
- ✅ Modern Go frameworks and libraries
- ✅ Robust database and caching solutions
- ✅ Comprehensive testing tools
- ✅ API documentation generation
- ✅ Runtime dependency injection with container
- ✅ Flexible dependency management
- ✅ Code quality tools
- ✅ CI/CD integration
- ✅ Security best practices
- ✅ Performance optimization
- ✅ Monitoring and observability
- ✅ Graceful server initialization and shutdown

**Key Features**:
- Container-based dependency injection for all layers
- Proper initialization order and lifecycle management
- Graceful shutdown with resource cleanup
- Flexible dependency configuration
- Support for testing with mock dependencies
