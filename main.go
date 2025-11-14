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

const serverInstructions = `Task Management

Use the todo_read and todo_write tools to manage work items. Use todo_read to retrieve the current list (in markdown format), then todo_write to save an updated list in markdown format. This approach supports multi-step tasks, cross-file work, notes, and ambiguous scopes.

Workflow:
- Begin: read current list
- Update: modify list as needed
- Finish: write back the revised list and verify completion

Important: todo_write replaces the entire file. Always read first; writing without reading may lose data.

Guidelines:
- Keep entries concise and actionable.
- Do not skip the tools for complex tasks.

Template:
- [ ] Implement feature X
  - [ ] Update API
  - [ ] Write tests
  - [ ] Run tests
  - [ ] Run lint
- [ ] Blocked: waiting on credentials`

func main() {
	todo := &TodoServer{
		logger:   log.New(os.Stderr, "", 0),
		fallback: "No todos found.",
	}

	if err := todo.loadFromFile(); err != nil {
		todo.logger.Printf("failed to load persisted todos: %v", err)
	}

	impl := &mcp.Implementation{Name: "todo", Version: "0.1.0"}

	mcpServer := mcp.NewServer(impl, &mcp.ServerOptions{Instructions: serverInstructions})

	for _, p := range getTodoPrompts() {
		mcpServer.AddPrompt(&p, func(_ context.Context, _ *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			return &mcp.GetPromptResult{
				Description: p.Description,
				Messages: []*mcp.PromptMessage{getTodoPromptMessage()},			}, nil
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
