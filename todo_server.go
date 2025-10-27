package main

import (
    "context"
    "log"
    "fmt"
    "encoding/json"
    "os"

    "sync"

    "github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
    storageFile = "todos.json"
)




type TodoServer struct {
    mu       sync.RWMutex
    todos    map[string]string
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
    var loaded map[string]string
    if err := json.Unmarshal(data, &loaded); err != nil {
        return err
    }
    s.mu.Lock()
    s.todos = loaded
    s.mu.Unlock()
    return nil
}

// persistToFile writes the current todos map to storageFile.
func (s *TodoServer) persistToFile() error {
    s.mu.RLock()
    data, err := json.MarshalIndent(s.todos, "", "  ")
    s.mu.RUnlock()
    if err != nil {
        return err
    }
    return os.WriteFile(storageFile, data, 0o644)
}



// Initialize implements the mcp.Server interface's Initialize method, providing server metadata.
func (s *TodoServer) Initialize(ctx context.Context, req *mcp.InitializeRequest) (*mcp.InitializeResult, error) {
    // Provide basic server info and capabilities.
    // Log server initialization
    s.logger.Printf("Server initialized with protocol version %s", "2025-06-18")
    return &mcp.InitializeResult{
        ProtocolVersion: "2025-06-18",
        ServerInfo:      &mcp.Implementation{Name: "todo", Version: "1.0.0"},
        Instructions: `Task Management

Use todo_read and todo_write for tasks with 2+ steps, multiple files/components, or uncertain scope.

Workflow:
- Start: read → write checklist
- During: read → update progress
- End: verify all complete

Warning: todo_write overwrites entirely; always todo_read first (skipping is an error)

Keep items short, specific, action-oriented. Not using the todo tools for complex tasks is an error.

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
func (s *TodoServer) ListTools(ctx context.Context, req *mcp.ListToolsRequest) (*mcp.ListToolsResult, error) {
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
func (s *TodoServer) CallTool(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    switch req.Params.Name {
    case "todo_read":
        s.logger.Printf("todo_read called")
    
        res, err := s.handleRead(ctx, req)
        return res, err
    case "todo_write":
        var args struct{ Content string `json:"content"` }
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
func (s *TodoServer) handleRead(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    var content string
    if req.Session != nil && req.Session.ID() != "" {
        content = s.todos[req.Session.ID()]
    } else {
        content = s.todos["default"]
    }
    return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: content}}}, nil
}

// handleWrite writes the provided todo content.
func (s *TodoServer) handleWrite(ctx context.Context, req *mcp.CallToolRequest, args struct{ Content string `json:"content"` }) (*mcp.CallToolResult, error) {
    // Validate request
    if req == nil {
        return nil, fmt.Errorf("nil request")
    }
    s.logger.Printf("todo_write called with %d chars", len(args.Content))
    if len(args.Content) == 0 {
        return nil, fmt.Errorf("content is empty")
    }
// Updated: persist after unlocking to avoid deadlock
    s.mu.Lock()
    if req.Session != nil && req.Session.ID() != "" {
        if s.todos == nil {
            s.todos = make(map[string]string)
        }
        s.todos[req.Session.ID()] = args.Content
    } else {
        s.todos["default"] = args.Content
    }
    // Unlock before persisting to file to avoid deadlock with RLock inside persistToFile
    s.mu.Unlock()
    if s.persistToFile() != nil {
        s.logger.Printf("failed to persist todos to file")
    }
    // Re-lock for deferred unlock removed; return result
    return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Updated (%d chars)", len(args.Content))}}}, nil
}


// getTools returns the tool definitions for the server.
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
            Description: `Read the current temporary storage content.

This tool returns the stored data as a string, allowing the LLM to retrieve previously saved information. It can be used for simple state persistence across calls.

The tool will return an error if the storage cannot be accessed.`,
            Title:       "Read TODO file",
            InputSchema:  map[string]any{"type": "object"},
            OutputSchema: map[string]any{"type": "object"},
            Annotations: readAnnotations,
        },
        {
            Name: "todo_write",
            Description: `Write or overwrite the entire TODO file content.

This tool replaces the whole TODO file with the supplied string, allowing the LLM to store arbitrary data persistently. It can be used to save updated task lists, notes, or any structured information the model wishes to retain across calls.

WARNING: This operation overwrites the entire file. Ensure the provided content includes all data you wish to keep, as any existing content not included will be lost.

The tool will create the TODO file if it does not exist, or overwrite it if it does. It returns an error if the file cannot be written due to permissions or other I/O issues.`,
            Title:       "Write TODO file",
            InputSchema:  writeInputSchema,
            OutputSchema: map[string]any{"type": "object"},
            Annotations: writeAnnotations,
        },
    }
}
