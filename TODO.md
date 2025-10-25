# TODO - Go MCP Todo Server Implementation

This document outlines the steps to implement a Go-based MCP server for a todo list, mirroring the functionality of the Rust example.

## 1. Bug Fixes
- [x] Fix server initialization in `main.go`. The `TodoServer` instance should be passed to `mcp.NewServer` as the handler, instead of registering tool handlers individually. This will ensure the server's state (mutex and fallback content) is used correctly.
- [x] Fix `TestEndToEnd` in `todo_server_test.go`. Similar to `main.go`, the test should use the `TodoServer` as the handler for the test server, not individual functions.
- [x] Fix `TestHandleReadWrite` to initialize the logger to prevent potential panics.

## 2. Feature Parity with Rust Example
- [x] Implement session-based state management. The current implementation only has a single global todo list. This should be updated to store todo lists per session ID, with a fallback for requests without a session ID.
- [x] Enhance tool descriptions in `getTools()` to match the detail in the Rust example.
- [x] Improve server instructions in `Initialize()` to provide a more detailed workflow for the user, as seen in the Rust example.

## 3. Refinement
- [ ] Add more robust error handling, such as checking for content size limits in `handleWrite`.
- [ ] Add logging for all public methods of `TodoServer`.
- [ ] Add comments to the code where necessary to clarify the implementation.

---
*Original checklist below for reference.*

## Project Setup
- [x] Initialize the Go module.
- [x] Add the `github.com/modelcontextprotocol/go-sdk` dependency.

## MCP Server Scaffolding
- [x] Create a `main.go` file.
- [x] Define a `TodoServer` struct.

## State Management
- [x] Implement a simple in-memory state management solution for the todo list content. A `sync.RWMutex` protecting a `map[string]string` to store todo lists by session ID would be a good starting point.
- [x] Add a fallback mechanism for when no session ID is present, similar to the Rust example.

## Tool Implementation
- [x] Implement the `todo_read` tool.
- [x] Implement the `todo_write` tool.

## Server Initialization
- [x] In the `Initialize` method, provide the server's metadata.

## Main Function
- [x] In `main.go`, create an instance of the `TodoServer`.
- [x] Start the MCP server using the `mcp.NewServer` and `Run` functions from the SDK.

## Testing
- [x] Create a `todo_server_test.go` file.
- [x] Write unit tests for the `todo_read` and `todo_write` tools.
- [x] Test the server's initialization and error handling.
