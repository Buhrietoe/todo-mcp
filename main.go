package main

import (
    "context"
    "log"
    "os"
    "os/signal"

    "github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
    defer cancel()

    // Initialize TodoServer implementation.
    ts := &TodoServer{}
    // Create MCP server with instructions.
    opts := &mcp.ServerOptions{Instructions: "Use the todo_read and todo_write tools to manage your todo list."}
    server := mcp.NewServer(&mcp.Implementation{Name: "todo", Version: "1.0.0"}, opts)
    // Register tools.
    mcp.AddTool(server, &mcp.Tool{Name: "todo_read", Description: "Read the entire TODO list", Title: "Read TODO"}, ts.handleRead)
    mcp.AddTool(server, &mcp.Tool{Name: "todo_write", Description: "Write or overwrite the entire TODO list", Title: "Write TODO"}, ts.handleWrite)

    // Run server over Stdio with logging.
    t := &mcp.LoggingTransport{Transport: &mcp.StdioTransport{}, Writer: os.Stderr}
    if err := server.Run(ctx, t); err != nil {
        log.Fatal(err)
    }
}
