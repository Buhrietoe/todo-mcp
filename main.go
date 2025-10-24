package main

import (
    "context"
    "fmt"
    "log"

    "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Removed duplicated TodoServer implementation; now in todo_server.go


func main() {
    ctx := context.Background()
    clientTransport, serverTransport := mcp.NewInMemoryTransports()

    // Create the MCP server.
    server := mcp.NewServer(&mcp.Implementation{Name: "todo", Version: "1.0.0"}, nil)
    todo := &TodoServer{}
    // Register tools.
    for _, tool := range getTools() {
        switch tool.Name {
        case "todo_read":
            mcp.AddTool[any, any](server, &tool, todo.handleRead)
        case "todo_write":
            mcp.AddTool[struct{ Content string `json:"content"` }, any](server, &tool, todo.handleWrite)
        }
    }


    // Run server side.
    serverSession, err := server.Connect(ctx, serverTransport, nil)
    if err != nil {
        log.Fatal(err)
    }

    // Connect client.
    client := mcp.NewClient(&mcp.Implementation{Name: "client"}, nil)
    clientSession, err := client.Connect(ctx, clientTransport, nil)
    if err != nil {
        log.Fatal(err)
    }

    // Example usage: write then read.
    _, err = clientSession.CallTool(ctx, &mcp.CallToolParams{Name: "todo_write", Arguments: map[string]any{"content": "- [ ] Sample task"}})
    if err != nil {
        log.Fatal(err)
    }
    res, err := clientSession.CallTool(ctx, &mcp.CallToolParams{Name: "todo_read"})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(res.Content[0].(*mcp.TextContent).Text)

    // Clean up.
    clientSession.Close()
    serverSession.Wait()
}
