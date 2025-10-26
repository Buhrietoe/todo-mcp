package main

import (
    "context"
    "io"
    "log"
    "encoding/json"
    "fmt"
    "testing"

    "github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestHandleReadWrite(t *testing.T) {
    s := &TodoServer{logger: log.New(io.Discard, "", 0), todos: make(map[string]string)}
    ctx := context.Background()
    req := &mcp.CallToolRequest{}
    // Write
    writeArgs := struct{ Content string `json:"content"` }{Content: "- [ ] Task 1"}
    _, err := s.handleWrite(ctx, req, writeArgs)
    if err != nil {
        t.Fatalf("handleWrite error: %v", err)
    }
    // Read
    res, err := s.handleRead(ctx, req)
    if err != nil {
        t.Fatalf("handleRead error: %v", err)
    }
    txt := res.Content[0].(*mcp.TextContent).Text
    if txt != "- [ ] Task 1" {
        t.Fatalf("expected content '- [ ] Task 1', got %q", txt)
    }
}

func TestEndToEnd(t *testing.T) {
    ctx := context.Background()
    clientTransport, serverTransport := mcp.NewInMemoryTransports()
    todo := &TodoServer{logger: log.New(io.Discard, "", 0), todos: make(map[string]string)}
    impl := &mcp.Implementation{Name: "todo", Version: "0.1.0"}
    server := mcp.NewServer(impl, &mcp.ServerOptions{})
    // Register tools
    for _, tool := range getTools() {
        t := tool // copy to avoid range variable reuse
        switch t.Name {
        case "todo_read":
            server.AddTool(&t, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
                return todo.handleRead(ctx, req)
            })
        case "todo_write":
            server.AddTool(&t, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
                var args struct{ Content string `json:"content"` }
                if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
                    return nil, fmt.Errorf("failed to parse arguments: %w", err)
                }
                return todo.handleWrite(ctx, req, args)
            })
        }
    }
    // No need to add tools manually; server uses CallTool implementation
    serverSession, err := server.Connect(ctx, serverTransport, nil)
    if err != nil {
        t.Fatalf("server connect: %v", err)
    }
    client := mcp.NewClient(&mcp.Implementation{Name: "client"}, nil)
    clientSession, err := client.Connect(ctx, clientTransport, nil)
    if err != nil {
        t.Fatalf("client connect: %v", err)
    }
    // Write via client
    _, err = clientSession.CallTool(ctx, &mcp.CallToolParams{Name: "todo_write", Arguments: map[string]any{"content": "- [ ] End2End"}})
    if err != nil {
        t.Fatalf("client write: %v", err)
    }
    // Read via client
    res, err := clientSession.CallTool(ctx, &mcp.CallToolParams{Name: "todo_read"})
    if err != nil {
        t.Fatalf("client read: %v", err)
    }
    txt := res.Content[0].(*mcp.TextContent).Text
    if txt != "- [ ] End2End" {
        t.Fatalf("expected '- [ ] End2End', got %q", txt)
    }
    clientSession.Close()
    serverSession.Wait()
}
