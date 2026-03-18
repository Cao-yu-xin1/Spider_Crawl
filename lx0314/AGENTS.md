# AGENTS.md - lx0314 Codebase Guidelines

## Project Overview

This is a Go microservices project with a layered architecture:
- `bff-api/` - Backend-for-Frontend layer (Gin HTTP server)
- `srv/` - Service/Backend layer (gRPC services with GORM models)
- `nacos/` - Nacos configuration management
- `consul/` - Consul service discovery (placeholder)
- `mq/` - Message queue handling (placeholder)
- `proto/` - Protocol Buffer definitions (placeholder)

## Build and Configuration

### Dependencies
- Module: `lx0314`
- Go version: 1.25.8
- Key dependencies: gin (HTTP), grpc (gRPC), gorm (ORM), nacos-sdk-go, viper

### Key Commands
```bash
# Run bff-api service
go run bff-api/basic/cmd/main.go

# Run srv service  
go run srv/basic/cmd/main.go

# Install dependencies
go mod download
go mod tidy

# Run go vet for static analysis
go vet ./...

# Format all code
go fmt ./...
```

### Project Structure
```
lx0314/
├── bff-api/           # BFF layer
│   ├── basic/        # Shared basic setup
│   │   ├── cmd/     # Entry points
│   │   ├── config/  # Config structs
│   │   └── init/    # Initialization logic
│   ├── handler/      # Request/response/service handlers
│   ├── middleware/   # HTTP middleware
│   ├── pkg/          # Shared packages
│   └── router/       # Route definitions
├── srv/              # Backend services
│   ├── basic/        # Shared basic setup
│   ├── handler/      # gRPC handlers
│   ├── model/        # Database models (GORM)
│   ├── pkg/          # Shared packages
│   └── model/        # Data models
├── nacos/            # Nacos config client
├── consul/           # Consul integration
├── mq/               # Message queue consumers
├── proto/            # .proto definitions
├── go.mod            # Go modules file
└── nacos.yaml        # Nacos configuration
```

## Code Style Guidelines

### Naming Conventions
- **Structs**: PascalCase (e.g., `User`, `Product`, `AppConfig`)
- **Interfaces**: PascalCase with "er" suffix or adjective (e.g., `Service`, `Processor`)
- **Functions/Methods**: camelCase (e.g., `MysqlInit`, `GetConfig`, `handleRequest`)
- **Packages**: lowercase, singular (e.g., `model`, `handler`, `config`)
- **Variables**: camelCase (e.g., `userName`, `dbConnection`, `maxRetries`)

### File Organization
- Group by functionality/module
- One main struct per model file
- Separate handler layers: `request`, `response`, `service`
- Keep initialization logic in `init/` directory

### Error Handling
- Return errors explicitly from functions
- Log errors with context before returning
- Use standard Go error wrapping: `fmt.Errorf("context: %w", err)`
- Handle database errors via GORM error checking

### Comments
- All public functions/types should have doc comments
- Service functions comment what they do
- Model fields include gorm tags with comments
- Complex logic requires inline comments

### Database Models (GORM)
- Embed `gorm.Model` for standard fields (ID, CreatedAt, UpdatedAt, DeletedAt)
- Use gorm tags for column definitions: `gorm:"type:varchar(50);comment:描述"`
- Pluralize table names (GORM default behavior)
- Foreign keys: `FieldNameID int64`

### HTTP Layer (bff-api)
- Use Gin framework for HTTP routes
- Request handlers in `bff-api/handler/request/`
- Response formats in `bff-api/handler/response/`
- Service logic in `bff-api/handler/service/`
- Middleware in `bff-api/middleware/`

### gRPC Layer (srv)
-/grpc services in `srv/basic/cmd/main.go`
- Implement gRPC handlers in `srv/handler/`
- Database models in `srv/model/`
- Use gorm for all database operations

### Configuration
- Nacos config loaded via `nacos.App()` in `nacos/` package
- Environment-specific config in `nacos.yaml`
- Global config stored in package vars (e.g., `nacos.GlobalConfig`)
- Use `viper` for configuration management

### Logging
- Use `log.Println()` for simple logging
- Consider `zap` for production (already in dependencies)
- Log errors with meaningful context

## Testing

- No formal test structure yet
- Test files should use `_test.go` suffix
- Test package: `{package}_test`
- Use `go test ./...` to run all tests

## CI/CD

- No CI/CD configuration found yet
- Recommended: GitHub Actions or GitLab CI
- Build matrix: test with `go vet ./...` and `go test ./...`

## Development Notes

- Current codebase is in early development
- Directory structure is partially created (run `go run main.go` to create)
- MySQL is the current database (GORM + mysql driver)
- Nacos is used for configuration management
- Code follows Go idioms with clear separation of concerns

## Editor Settings

- Navigate to `E:\gowork\src\lx\zg5\lx0314`
- Use Go 1.25.8+ withgo module support
- Format with `go fmt`
- Run `go vet` before commits
