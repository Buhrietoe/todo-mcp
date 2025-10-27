# todo-mcp

A minimal Go server implementing Model Context Protocol (MCP) tools for persistent TODO list storage.

## Overview

The server provides two MCP tools:
- `todo_read` – returns the current TODO content for a session (or a default entry).
- `todo_write` – overwrites the TODO content for a session (or the default entry) and persists it to `todos.json`.

It uses `github.com/modelcontextprotocol/go-sdk/mcp` to expose these tools over STDIO, making it easy to integrate with any MCP‑compatible client.

## Audit Findings

All previously identified audit findings have been addressed:
- Unused parameters removed (replaced with `_`).
- Persistence errors now returned to the caller.
- Concurrency handling verified.
- GoDoc comments added to exported functions.
- Table‑driven tests implemented.

## Installation

```sh
# Install the binary to your $GOPATH/bin (or $HOME/go/bin)

go install -v github.com/Buhrietoe/todo-mcp@latest
```

Example MCP config:

```json
"todo": {
    "type": "stdio",
    "command": "todo-mcp"
}
```

Or use the provided Makefile shortcuts:

```sh
make build   # builds the binary `todo-mcp`
make run     # builds then runs the binary
```

## Testing

```sh
go test ./...
```

## Persistence

Todos are stored in `todos.json` in the working directory. The server loads this file on start (if present) and writes updates after each `todo_write` call.

## MCP Tool Definitions

- **todo_read** – read‑only, idempotent. Returns the stored string.
- **todo_write** – destructive, non‑idempotent. Requires a `content` field (string) and overwrites the stored value.

Both tools are registered automatically in `main.go` via `getTools()`.

## License

MIT License (see LICENSE file).
