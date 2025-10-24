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
