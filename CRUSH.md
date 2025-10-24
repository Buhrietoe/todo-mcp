# CRUSH.md

## Build / Run
- `go build ./...` – build all Go packages.
- `go run ./...` – run the main package (if present).

## Test
- `go test ./...` – run all Go tests.
- To run a single Go test: `go test ./... -run TestName`.
- To test MCP server functionality: `npx @modelcontextprotocol/inspector ./todo-mcp <arg1> <arg2>`

## Lint / Format
- `go vet ./...` – static analysis.
- `golint ./...` – lint Go code.

## Code Style Guidelines
- **Imports**: group standard library, third‑party, and local imports; separate groups with a blank line.
- **Formatting**: `go fmt` for Go, `rustfmt` for Rust.
- **Naming**:
  - Types: `CamelCase`.
  - Functions/variables: `camelCase` (exported functions start with uppercase).
  - Constants: `SCREAMING_SNAKE_CASE`.
- **Error handling**: return `error` as the last return value; wrap errors with context using `fmt.Errorf` or `anyhow`.
- **Documentation**: exported identifiers must have a comment starting with the name.
- **Testing**: use table‑driven tests; name test functions `TestXxx`.

## Project Structure
- `examples/` – language‑specific example implementations.
- `PLAN.md` – high‑level project plan.
- `README.md` – project overview.

*Generated for agentic tooling.*
