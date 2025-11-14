# AGENTS.md

## Build / Run
- `go build ./...` – build all Go packages.
- `mcptools prompts ./todo-mcp` - validate prompts run
- `mcptools tools ./todo-mcp` - validate tools run

## Test
- `go test ./...` – run all tests.
- Run a single test: `go test ./... -run ^TestName$`.
- Table‑driven tests are preferred.

## Lint / Format
- `go vet ./...` – static analysis.
- `golint ./...` – lint Go code.
- `go fmt ./...` – format source files.

## Code Style Guidelines
- **Imports**: group standard library, third‑party, and local imports; separate groups with a blank line.
- **Formatting**: use `go fmt`.
- **Naming**:
  - Types: `CamelCase`.
  - Exported functions/variables: `CamelCase`.
  - Unexported: `camelCase`.
  - Constants: `SCREAMING_SNAKE_CASE`.
- **Error handling**: return `error` as the last return value; wrap with `fmt.Errorf("%w", err)` or `errors.New`.
- **Documentation**: exported identifiers need a comment starting with the name.
- **Unused params**: prefix with `_` if intentionally unused.
- **Testing**: use table‑driven tests; test functions named `TestXxx`.

## Project Structure
- `examples/` – language‑specific examples.
- `main.go` – entry point.
- `todo_server.go` – server implementation.
- `todo_server_test.go` – tests.

## Additional Agent Guidance

### Tool Usage
- `todo_read`: reads the current TODO list from `TODO.md` (markdown format). Returns the content as a string via the `mcp` tool interface.
- `todo_write`: overwrites the entire `TODO.md` file. **Must read first**; writing without reading may lose data. Passing empty content clears the file.
- Both tools are thread‑safe; the server protects access with a `sync.RWMutex`.

### Dependencies
- Go 1.25 (`go.mod`).
- `github.com/modelcontextprotocol/go-sdk` provides the MCP framework.
