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

    
    // Create MCP server with the TodoServer handler and options.
    server := mcp.NewServer(ts, &mcp.ServerOptions{Instructions: "Use the todo_read and todo_write tools to manage your todo list."})

    // Run server over Stdio with logging.
    t := &mcp.LoggingTransport{Transport: &mcp.StdioTransport{}, Writer: os.Stderr}
    if err := server.Run(ctx, t); err != nil {
        log.Fatal(err)
    }
}
