package main

import (
    "context"
    "io"
    "log"
    "testing"

    "github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestHandleReadWrite(t *testing.T) {
    s := &TodoServer{logger: log.New(io.Discard, "", 0), todos: make(map[string]string)}
    ctx := context.Background()
    req := &mcp.CallToolRequest{}
    // Write
    writeArgs := struct{ Content string `json:"content"` }{Content: "- [ ] Task 1"}
    _, _, err := s.handleWrite(ctx, req, writeArgs)
    if err != nil {
        t.Fatalf("handleWrite error: %v", err)
    }
    // Read
    res, _, err := s.handleRead(ctx, req, nil)
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
    server := mcp.NewServer(&mcp.Implementation{Name: "todo", Version: "1.0.0"}, nil)
    for _, tool := range getTools() {
        switch tool.Name {
        case "todo_read":
            mcp.AddTool(server, &tool, todo.handleRead)
        case "todo_write":
            mcp.AddTool(server, &tool, todo.handleWrite)
        }
    }
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
