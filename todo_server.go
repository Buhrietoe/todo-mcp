package main

import (
    "context"
    "fmt"
    "sync"

    "github.com/modelcontextprotocol/go-sdk/mcp"
)

type TodoServer struct {
    mu       sync.RWMutex
    fallback string
}

// Initialize implements the mcp.Server interface's Initialize method, providing server metadata.
func (s *TodoServer) Initialize(ctx context.Context, req *mcp.InitializeRequest) (*mcp.InitializeResult, error) {
    // Provide basic server info and capabilities.
    return &mcp.InitializeResult{
        ProtocolVersion: "2025-06-18",
        ServerInfo:      &mcp.Implementation{Name: "todo", Version: "1.0.0"},
        Instructions:    "Use the todo_read and todo_write tools to manage your TODO list.",
        Capabilities: &mcp.ServerCapabilities{
            Tools: &mcp.ToolCapabilities{ListChanged: true},
        },
    }, nil
}

// ListTools returns the list of tools supported by the server.
func (s *TodoServer) ListTools(ctx context.Context, req *mcp.ListToolsRequest) (*mcp.ListToolsResult, error) {
    tools := getTools()
    // Convert []mcp.Tool to []*mcp.Tool
    ptrs := make([]*mcp.Tool, len(tools))
    for i := range tools {
        ptrs[i] = &tools[i]
    }
    return &mcp.ListToolsResult{Tools: ptrs}, nil
}

// CallTool dispatches tool calls to the appropriate handler. This is a generic fallback.
func (s *TodoServer) CallTool(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // Since we register handlers directly via AddTool, this method can return a not‑implemented error.
    return nil, fmt.Errorf("CallTool not implemented for %s", req.Params.Name)
}

// handleRead reads the current todo content.
func (s *TodoServer) handleRead(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
    s.mu.RLock()
    content := s.fallback
    s.mu.RUnlock()
    return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: content}}}, nil, nil
}

// handleWrite writes the provided todo content.
func (s *TodoServer) handleWrite(ctx context.Context, req *mcp.CallToolRequest, args struct{ Content string `json:"content"` }) (*mcp.CallToolResult, any, error) {
    s.mu.Lock()
    s.fallback = args.Content
    s.mu.Unlock()
    msg := fmt.Sprintf("Updated (%d chars)", len(args.Content))
    return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: msg}}}, nil, nil
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
            Name:        "todo_read",
            Description: "Read the entire TODO file content.",
            Title:       "Read TODO",
            InputSchema:  map[string]any{"type": "object"},
            OutputSchema: map[string]any{"type": "object"},
            Annotations: readAnnotations,
        },
        {
            Name:        "todo_write",
            Description: "Write or overwrite the entire TODO file content.",
            Title:       "Write TODO",
            InputSchema:  writeInputSchema,
            OutputSchema: map[string]any{"type": "object"},
            Annotations: writeAnnotations,
        },
    }
}
