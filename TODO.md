# TODO - Go MCP Todo Server Implementation

This document outlines the steps to implement a Go-based MCP server for a todo list, mirroring the functionality of the Rust example.

## 1. Project Setup
- [ ] Initialize the Go module.
- [ ] Add the `github.com/modelcontextprotocol/go-sdk` dependency.

## 2. MCP Server Scaffolding
- [ ] Create a `main.go` file.
- [ ] Define a `TodoServer` struct.
- [ ] Implement the `mcp.Server` interface for `TodoServer`.
  - [ ] `Initialize(context.Context, *mcp.InitializeRequest) (*mcp.InitializeResponse, error)`
  - [ ] `ListTools(context.Context, *mcp.ListToolsRequest) (*mcp.ListToolsResponse, error)`
  - [ ] `CallTool(context.Context, *mcp.CallToolRequest) (*mcp.CallToolResponse, error)`
  - [ ] Implement stubs for other `mcp.Server` methods (e.g., `ListResources`, `ReadResource`, etc.) to return "not implemented" errors.

## 3. State Management
- [ ] Implement a simple in-memory state management solution for the todo list content. A `sync.RWMutex` protecting a `map[string]string` to store todo lists by session ID would be a good starting point.
- [ ] Add a fallback mechanism for when no session ID is present, similar to the Rust example.

## 4. Tool Implementation
- [ ] Implement the `todo_read` tool:
  - [ ] In the `CallTool` method, check for the "todo_read" tool name.
  - [ ] Read the todo list content from the state management system.
  - [ ] Return the content in a `mcp.CallToolResponse`.
- [ ] Implement the `todo_write` tool:
  - [ ] In the `CallTool` method, check for the "todo_write" tool name.
  - [ ] Extract the "content" argument from the `mcp.CallToolRequest`.
  - [ ] Update the todo list content in the state management system.
  - [ ] Return a success message in a `mcp.CallToolResponse`.

## 5. Server Initialization
- [ ] In the `Initialize` method, provide the server's metadata:
  - [ ] Protocol version.
  - [ ] Server capabilities (tools).
  - [ ] Server information (name, version, etc.).
  - [ ] Instructions for using the todo tools.

## 6. Main Function
- [ ] In `main.go`, create an instance of the `TodoServer`.
- [ ] Start the MCP server using the `mcp.NewServer` and `ListenAndServe` functions from the SDK.

## 7. Testing
- [ ] Create a `todo_server_test.go` file.
- [ ] Write unit tests for the `todo_read` and `todo_write` tools.
- [ ] Test the server's initialization and error handling.

## 8. Refinement
- [ ] Add logging to the server.
- [ ] Improve error handling and provide more informative error messages.
- [ ] Add comments to the code where necessary.
