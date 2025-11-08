package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	storageFile = "TODO.md"
)

// TodoServer implements the mcp.Server interface and provides todo management functionality.
type TodoServer struct {
	mu       sync.RWMutex
	content  string
	fallback string
	logger   *log.Logger
}

// loadFromFile loads persisted todos from storageFile if it exists.
func (s *TodoServer) loadFromFile() error {
	data, err := os.ReadFile(storageFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if len(data) == 0 {
		return nil
	}

	s.mu.Lock()
	s.content = string(data)
	s.mu.Unlock()
	return nil
}

// persistToFile writes the default todo content to storageFile.
func (s *TodoServer) persistToFile() error {
	s.mu.RLock()
	content := s.content
	s.mu.RUnlock()
	return os.WriteFile(storageFile, []byte(content), 0o644)
}

// Initialize implements the mcp.Server interface's Initialize method, providing server metadata.
// Initialize implements the mcp.Server interface's Initialize method, providing server metadata.
func (s *TodoServer) Initialize(_ context.Context, req *mcp.InitializeRequest) (*mcp.InitializeResult, error) {
	// Provide basic server info and capabilities.
	// Log server initialization
	s.logger.Printf("Server initialized with protocol version %s", "2025-06-18")
	return &mcp.InitializeResult{
		ProtocolVersion: "2025-06-18",
		ServerInfo:      &mcp.Implementation{Name: "todo", Version: "1.0.0"},
		Instructions: `Task Management

Use the todo_read and todo_write tools to manage work items. Use todo_read to retrieve the current list, then use todo_write to save an updated list. This approach supports multi-step tasks, cross-file work, notes, and ambiguous scopes.

Workflow:
- Begin: read current list
- Update: modify list as needed
- Finish: write back the revised list and verify completion

Important: todo_write replaces the entire file. Always read first; writing without reading may lose data.

Guidelines:
- Keep entries concise and actionable.
- Do not skip the tools for complex tasks.

Template:
- [ ] Implement feature X
  - [ ] Update API
  - [ ] Write tests
  - [ ] Run tests
  - [ ] Run lint
- [ ] Blocked: waiting on credentials`,
		Capabilities: &mcp.ServerCapabilities{
			Tools: &mcp.ToolCapabilities{ListChanged: true},
		},
	}, nil
}

// ListTools returns the list of tools supported by the server.
// ListTools returns the list of tools supported by the server.
func (s *TodoServer) ListTools(_ context.Context, req *mcp.ListToolsRequest) (*mcp.ListToolsResult, error) {
	s.logger.Printf("ListTools called")
	tools := getTools()
	// Convert []mcp.Tool to []*mcp.Tool
	ptrs := make([]*mcp.Tool, len(tools))
	for i := range tools {
		ptrs[i] = &tools[i]
	}
	return &mcp.ListToolsResult{Tools: ptrs}, nil
}

// CallTool dispatches tool calls to the appropriate handler.
// CallTool dispatches a tool call to the appropriate handler.
func (s *TodoServer) CallTool(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	switch req.Params.Name {
	case "todo_read":
		s.logger.Printf("todo_read called")

		res, err := s.handleRead(ctx, req)
		return res, err
	case "todo_write":
		var args struct {
			Content string `json:"content"`
		}
		if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
			return nil, fmt.Errorf("failed to parse arguments: %w", err)
		}
		res, err := s.handleWrite(ctx, req, args)
		return res, err
	default:
		return nil, fmt.Errorf("unknown tool %s", req.Params.Name)
	}
}

// handleRead reads the current todo content.
// handleRead returns the current todo content.
func (s *TodoServer) handleRead(_ context.Context, _ *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	content := s.content
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: content}}}, nil
}

// handleWrite writes the provided todo content.
// handleWrite updates the todo content and persists it.
func (s *TodoServer) handleWrite(_ context.Context, req *mcp.CallToolRequest, args struct {
	Content string `json:"content"`
}) (*mcp.CallToolResult, error) {
	// Validate request
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}
	s.logger.Printf("todo_write called with %d chars", len(args.Content))
	if len(args.Content) == 0 {
		// Clear todo entry (ignore session)
		// Truncate storage file to clear persisted data
		if err := os.WriteFile(storageFile, []byte{}, 0o644); err != nil {
			s.logger.Printf("failed to clear storage file: %v", err)
			return nil, fmt.Errorf("failed to clear storage: %w", err)
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "Cleared todo list"}}}, nil
	}
	// Update content and persist
	s.mu.Lock()
	s.content = args.Content
	s.mu.Unlock()
	if err := s.persistToFile(); err != nil {
		s.logger.Printf("failed to persist todos to file: %v", err)
		return nil, fmt.Errorf("failed to persist todos: %w", err)
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Updated (%d chars)", len(args.Content))}}}, nil
}

// getTools returns the tool definitions for the server.
// getTools constructs and returns the tool definitions for the server.
func getTools() []mcp.Tool {
	// Define tool annotations and schemas
	openWorldHint := false
	readAnnotations := &mcp.ToolAnnotations{ReadOnlyHint: true, IdempotentHint: true, OpenWorldHint: &openWorldHint}
	writeDestructive := true
	writeAnnotations := &mcp.ToolAnnotations{DestructiveHint: &writeDestructive, IdempotentHint: false, OpenWorldHint: &openWorldHint}

	// Input schema for write tool
	writeInputSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"content": map[string]any{
				"type":        "string",
				"description": "The TODO list content to save",
			},
		},
		"required": []string{"content"},
	}

	return []mcp.Tool{
		{
			Name: "todo_read",
			Description: `Read the current temporary storage content (in markdown format).

This tool returns the stored data as a string, allowing the LLM to retrieve previously saved information. It can be used for simple state persistence across calls.

The tool will return an error if the storage cannot be accessed.`,
			Title:        "Read TODO file",
			InputSchema:  map[string]any{"type": "object"},
			OutputSchema: map[string]any{"type": "object"},
			Annotations:  readAnnotations,
		},
		{
			Name: "todo_write",
			Description: `Write or overwrite the entire TODO file content (in markdown format).

This tool replaces the whole TODO file with the supplied string, allowing the LLM to store arbitrary data persistently. It can be used to save updated task lists, notes, or any structured information the model wishes to retain across calls.

WARNING: This operation overwrites the entire file. Ensure the provided content includes all data you wish to keep, as any existing content not included will be lost.

The tool will create the TODO file if it does not exist, or overwrite it if it does. It returns an error if the file cannot be written due to permissions or other I/O issues.`,
			Title:        "Write TODO file",
			InputSchema:  writeInputSchema,
			OutputSchema: map[string]any{"type": "object"},
			Annotations:  writeAnnotations,
		},
	}
}
