# CRUSH.md

## Build / Run
- `go build ./...` – build all Go packages.
- `go run ./...` – run the main package.
- `make build` – build binary `todo-mcp`.
- `make run` – build then execute the binary.

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
- `todos.json` – sample data.

*Generated for agentic tooling.*
