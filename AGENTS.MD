# AGENTS.md

## Project Overview
A Go server implementing Model Context Protocol (MCP) tools for persistent TODO list storage. The server provides `todo_read` and `todo_write` tools that allow LLMs to manage task lists across conversations.

## Essential Commands

### Build / Run
- `go build ./...` – build all Go packages.
- `go run main.go` – run the server directly.
- `mcptools prompts ./todo-mcp` - validate prompts run
- `mcptools tools ./todo-mcp` - validate tools run

### Test
- `go test ./...` – run all tests.
- Run a single test: `go test ./... -run ^TestName$`.
- Table-driven tests are preferred.

### Lint / Format
- `go vet ./...` – static analysis.
- `golint ./...` – lint Go code (if golangci-lint is available).
- `go fmt ./...` – format source files.

## Code Organization and Structure

### File Structure
- `main.go` – entry point, MCP server setup, and initialization
- `todo_server.go` – core server implementation (`TodoServer` struct)
- `todo_server_test.go` – unit tests and integration tests
- `tool_definitions.go` – MCP tool schemas and definitions
- `prompts.go` – MCP prompt definitions
- `TODO.md` – persistent todo storage file (created automatically)

### Key Components
- **TodoServer**: Main server struct with mutex protection for concurrent access
- **Storage**: Uses `TODO.md` file in working directory for persistence
- **MCP Tools**: Exposes `todo_read` and `todo_write` via MCP protocol

## Code Style Guidelines

### Go Conventions
- **Imports**: group standard library, third-party, and local imports; separate groups with a blank line.
- **Formatting**: use `go fmt`.
- **Naming**:
  - Types: `CamelCase`.
  - Exported functions/variables: `CamelCase`.
  - Unexported: `camelCase`.
  - Constants: `SCREAMING_SNAKE_CASE`.
- **Error handling**: return `error` as the last return value; wrap with `fmt.Errorf("%w", err)` or `errors.New`.
- **Documentation**: exported identifiers need a comment starting with the name.
- **Unused params**: prefix with `_` if intentionally unused.
- **Testing**: use table-driven tests; test functions named `TestXxx`.

### Struct Patterns
- Use `sync.RWMutex` for concurrent-safe access to shared state
- Embedding is minimal - prefer composition over embedding
- Clear separation between MCP interface methods and internal handlers

## Testing Approach and Patterns

### Test Structure
- **Unit Tests**: Focus on individual functions (`handleRead`, `handleWrite`)
- **Integration Tests**: End-to-end testing with in-memory MCP transports
- **Table-driven Tests**: Preferred for parameterized test cases

### Key Testing Patterns
- Use `io.Discard` for test loggers to avoid noise
- Test both success and error paths
- Verify atomic file operations for persistence
- Test concurrent access scenarios if applicable

## MCP Tool Usage

### todo_read Tool
- **Purpose**: Read current TODO content from `TODO.md`
- **Behavior**: 
  - Read-only, idempotent, returns cached content
  - Reloads from file on each call to ensure freshness
  - Returns empty string if no todos exist
- **Usage**: Always read before writing to understand current state

### todo_write Tool
- **Purpose**: Overwrite entire TODO content with new content
- **Behavior**:
  - Destructive operation - completely replaces existing content
  - Creates `TODO.md` file if it doesn't exist
  - Atomic write operation to prevent corruption
  - Thread-safe with mutex protection
- **Usage**: MUST read first before writing to avoid data loss

## Important Gotchas and Patterns

### Critical Rules
1. **ALWAYS read before writing**: Writing without reading may lose data since `todo_write` replaces the entire file
2. **NEVER directly access TODO.md**: Only use `todo_read` and `todo_write` tools - never view or edit the file directly
3. **Content format**: Uses markdown format for todo lists (supports nested items with `- [ ]` syntax)

### Persistence Patterns
- File is loaded on server startup (`loadFromFile`)
- Content is cached in memory after load
- Writes are atomic - uses `O_CREATE|O_WRONLY|O_EXCL` then `O_WRONLY|O_TRUNC`
- Mutex protects against concurrent access (`sync.RWMutex`)

### Error Handling
- File not found errors are treated as "no todos" (not an error)
- Empty files are treated as valid (no todos)
- All file operations return appropriate errors for permission issues
- JSON unmarshal errors in tool calls return descriptive messages

### MCP Integration
- Uses `github.com/modelcontextprotocol/go-sdk/mcp` framework
- Tools are registered with appropriate annotations (read-only, destructive hints)
- Prompts are dynamically populated with current todo content
- Server provides instructions to users via MCP protocol

## Dependencies and Requirements

### Go Version
- Requires Go 1.25+ (specified in go.mod)

### External Dependencies
- `github.com/modelcontextprotocol/go-sdk v1.1.0` - MCP framework
- Indirect dependencies: google/jsonschema-go, yosida95/uritemplate, golang.org/x/oauth2

### Build Requirements
- Standard Go toolchain
- No external build tools required (no Makefile, just standard Go commands)

## Development Workflow

### Adding New Tools
1. Define tool in `tool_definitions.go` with proper schema and annotations
2. Add handler logic in `TodoServer.CallTool` method
3. Implement the actual handler function (e.g., `handleNewTool`)
4. Add corresponding tests in `todo_server_test.go`
5. Update documentation if behavior affects users

### Testing Changes
```bash
# Run all tests
go test ./...

# Run specific test
go test ./... -run TestHandleReadWrite

# Test with verbose output
go test -v ./...

# Run tests and show coverage
go test -cover ./...
```

### Common Patterns
- Context is passed to all handlers but not currently used for cancellation
- Logging uses standard Go `log.Logger` with timestamps via MCP transport
- All exported functions have proper Go comments
- Tool arguments are unmarshaled from JSON with strict typing
