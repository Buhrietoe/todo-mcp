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
		fallback: "No todos found.",
		version:  "1.0.0",
	}

	if err := todo.loadFromFile(); err != nil {
		if err != nil {
			todo.logger.Printf("Failed to load persisted todos: %v", err)
		}
	}

	// Initialize todo list using todo_write tool
	if _, err := todo.handleWrite(context.Background(), nil, struct {
		Content string `json:"content"`
	}{Content: "- [ ] Refactor tool definitions into a separate file for clarity\n- [ ] Abstract storage layer behind an interface to enable alternative persistence mechanisms\n- [ ] Add unit tests for loadFromFile and persistToFile edge cases (e.g., permission errors)\n- [ ] Review and standardize error messages and constants usage\n- [ ] Make logger configuration (flags, output) configurable via constructor\n- [ ] Document exported functions and types with proper Go comments\n- [ ] Ensure consistent use of context in all functions\n- [ ] Add linting configuration (e.g., golangci-lint) to enforce style"}); err != nil {
		todo.logger.Printf("failed to initialize todo list: %v", err)
	}
	impl := &mcp.Implementation{Name: "todo", Version: "1.0.0"}

	mcpServer := mcp.NewServer(impl, &mcp.ServerOptions{Instructions: serverInstructions})

	for _, p := range getTodoPrompts() {
		mcpServer.AddPrompt(&p, func(_ context.Context, _ *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			// Retrieve current TODO list content
			readRes, err := todo.handleRead(context.Background(), nil)
			if err != nil {
				return nil, err
			}
			var contentText string
			if len(readRes.Content) > 0 {
				if txt, ok := readRes.Content[0].(*mcp.TextContent); ok {
					contentText = txt.Text
				}
			}
			baseMsg := getTodoPromptMessageWithContent(contentText)
			return &mcp.GetPromptResult{
				Description: p.Description,
				Messages:    []*mcp.PromptMessage{baseMsg},
			}, nil
		})
	}

	for _, tool := range getTools() {
		t := tool
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

	t := &mcp.LoggingTransport{Transport: &mcp.StdioTransport{}, Writer: os.Stderr}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	if err := mcpServer.Run(ctx, t); err != nil {
		log.Fatal(err)
	}
}
