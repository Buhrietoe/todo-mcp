package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"log"
	"os"
	"os/signal"
)

func main() {
	todo := &TodoServer{
		logger:   log.New(os.Stderr, "", 0),
		todos:    make(map[string]string),
		fallback: "No todos found.",
	}
	// Load persisted todos from disk if available
	if err := todo.loadFromFile(); err != nil {
		todo.logger.Printf("failed to load persisted todos: %v", err)
	}
	impl := &mcp.Implementation{Name: "todo", Version: "0.1.0"}
	mcpServer := mcp.NewServer(impl, &mcp.ServerOptions{Instructions: "Use the todo_read and todo_write tools to manage your todo list, notes, or other temporary storage."})
	// Register tools
	for _, tool := range getTools() {
		t := tool // copy to avoid range variable reuse
		switch t.Name {
		case "todo_read":
			mcpServer.AddTool(&t, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				return todo.handleRead(ctx, req)
			})
		case "todo_write":
			mcpServer.AddTool(&t, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				var args struct {
					Content string `json:"content"`
				}
				if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
					return nil, fmt.Errorf("failed to parse arguments: %w", err)
				}
				return todo.handleWrite(ctx, req, args)
			})
		}
	}
	// Run server over Stdio with logging
	t := &mcp.LoggingTransport{Transport: &mcp.StdioTransport{}, Writer: os.Stderr}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	if err := mcpServer.Run(ctx, t); err != nil {
		log.Fatal(err)
	}
}
