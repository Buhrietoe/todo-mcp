# TODO - Go MCP Todo Server Implementation

This document outlines the steps to implement a Go-based MCP server for a todo list, mirroring the functionality of the Rust example.

## 1. Project Setup
- [x] Initialize the Go module.
- [x] Add the `github.com/modelcontextprotocol/go-sdk` dependency.

## 2. MCP Server Scaffolding
- [x] Create a `main.go` file.
- [x] Define a `TodoServer` struct.


## 3. State Management
- [x] Implement a simple in-memory state management solution for the todo list content. A `sync.RWMutex` protecting a `map[string]string` to store todo lists by session ID would be a good starting point.
- [x] Add a fallback mechanism for when no session ID is present, similar to the Rust example.

## 4. Tool Implementation
- [x] Implement the `todo_read` tool:
-   - [x] In the `CallTool` method, check for the "todo_read" tool name.
-   - [x] Read the todo list content from the state management system.
-   - [x] Return the content in a `mcp.CallToolResponse`.
- - [x] Implement the `todo_write` tool:
-   - [x] In the `CallTool` method, check for the "todo_write" tool name.
-   - [x] Extract the "content" argument from the `mcp.CallToolRequest`.
-   - [x] Update the todo list content in the state management system.
-   - [x] Return a success message in a `mcp.CallToolResponse`.

## 5. Server Initialization
- [x] In the `Initialize` method, provide the server's metadata:
-   - [ ] Protocol version.
-   - [ ] Server capabilities (tools).
-   - [ ] Server information (name, version, etc.).
-   - [ ] Instructions for using the todo tools.

## 6. Main Function
- [x] In `main.go`, create an instance of the `TodoServer`.
- [x] Start the MCP server using the `mcp.NewServer` and `Run` functions from the SDK.

## 7. Testing
- [x] Create a `todo_server_test.go` file.
- [x] Write unit tests for the `todo_read` and `todo_write` tools.
- [x] Test the server's initialization and error handling.

## 8. Refinement
- [ ] Add logging to the server.
- [ ] Improve error handling and provide more informative error messages.
- [ ] Add comments to the code where necessary.
