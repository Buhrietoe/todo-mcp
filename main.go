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
    ts := &TodoServer{logger: log.New(os.Stderr, "todo: ", log.LstdFlags)}
    
    // Create MCP server with instructions.
    server := mcp.NewServer(&mcp.Implementation{Name: "todo", Version: "1.0.0"}, &mcp.ServerOptions{Instructions: "Use the todo_read and todo_write tools to manage your todo list."})
    
    // Register tools using the TodoServer as the handler
    server.AddTool("todo_read", "Read the entire TODO list", ts)
    server.AddTool("todo_write", "Write or overwrite the entire TODO list", ts)

    // Run server over Stdio with logging.
    t := &mcp.LoggingTransport{Transport: &mcp.StdioTransport{}, Writer: os.Stderr}
    if err := server.Run(ctx, t); err != nil {
        log.Fatal(err)
    }
}
