# AGENTS.md - lx0318 Codebase Guidelines

## Project Overview
Go microservices project using gRPC, Nacos (configuration), Consul (service discovery), MySQL (GORM), and Message Queues.

## Build & Run Commands

### Build
```bash
go build -o bin/your-service ./path/to/main.go
```

### Run gRPC Server
```bash
go run ./rpc/basic/cmd/main.go
go run ./rpc/basic/cmd/consul.go  # Consul-based service
```

### Database Migration
Auto-migration runs on application startup via `rpc/basic/init/mysql.go`

### Test
Run tests in a package:
```bash
go test ./rpc/...
go test ./config
```

## Project Structure

```
lx0318/
├── bff-api/          # Backend-for-Frontend layer
├── config/           # Configuration management (Nacos)
├── mq/               # Message queue handlers
├── pkg/              # Shared packages (empty - use rpc/model instead)
├── proto/            # Protocol Buffer definitions
├── rpc/              # gRPC services and handlers
│   ├── basic/        # Base initialization (config, mysql)
│   ├── handler/      # gRPC service implementations
│   └── model/        # Data models (Member, Product, etc.)
├── nacos.yaml        # Nacos configuration
└── go.mod
```

## Code Style Guidelines

### Imports
- Group imports: standard library, then third-party, then internal
- Use absolute imports: `import "lx0318/rpc/model"`
- No aliases unless disambiguation needed

### Naming Conventions
- **Types**: PascalCase (e.g., `MemberRegister`, `Product`, `Server`)
- **Variables/Functions**: camelCase (e.g., `memberID`, `GetUser`, `GRPCServer`)
- **Constants**: UPPER_SNAKE_CASE (e.g., `DefaultTimeout`)
- **Files**: lowercase with hyphens or camelCase

### Error Handling
- Use `panic(err)` for initialization failures only
- Return errors explicitly: `return nil, fmt.Errorf("context: %w", err)`
- Wrap errors with context using `%w` verb
- Use `log.Fatalf` in main functions for fatal errors

### Comments
- Document all exported functions/types
- Include parameter descriptions and return values
- Add inline comments for non-obvious logic
- Use Chinese comments where team prefers it (observed in codebase)

### GORM Models
- Embed `gorm.Model` for base fields (ID, CreatedAt, UpdatedAt, DeletedAt)
- Use struct tags for column definitions:
  ```go
  Username string `gorm:"type:varchar(50);comment:用户名"`
  Status   int    `gorm:"type:tinyint;default:1;comment:状态"`
  ```

### Context Usage
- Pass `context.Context` as first parameter: `func(ctx context.Context, ...)`
- Use `context.WithTimeout` or `context.WithCancel` for operations
- Cancel contexts properly: `defer cancel()`

## Testing
- Test files: `*_test.go`
- Use table-driven tests for multiple scenarios
- Mock external dependencies
- Test package: `go test ./... -v`

## gRPC Patterns
- Register servers: `pb.RegisterServiceServer(s, &Handler{})`
- Implement unimplemented server interfaces
- Use context for request lifecycle
- Handle errors with gRPC status codes

## Configuration
- Nacos: Primary configuration source
- Fallback to local `nacos.yaml` for development
- Viper used for configuration management

## linter Rules
- No `as any` or type suppression
- No `@ts-ignore` (Go equivalent: no unsafe type assertions)
- Run `golangci-lint` if available

## Additional Notes
- Go version: 1.25.8
- UTF-8 encoding for all files
- Use `go fmt` before committing
- Avoid global state where possible
