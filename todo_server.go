package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	storageFile = "TODO.md"
)

// TodoServer implements the mcp.Server interface and provides todo management functionality.
type TodoServer struct {
	mu       sync.RWMutex
	content  string
	fallback string
	logger   *log.Logger
}

func (s *TodoServer) logf(format string, args ...any) { s.logger.Printf(format, args...) }

// loadFromFile loads persisted todos from storageFile if it exists.
func (s *TodoServer) loadFromFile() error {
	data, err := os.ReadFile(storageFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if len(data) == 0 {
		return nil
	}

	s.mu.Lock()
	s.content = string(data)
	defer s.mu.Unlock()
	return nil
}

// persistToFile writes the default todo content to storageFile.
func (s *TodoServer) persistToFile() error {
	s.mu.RLock()
	content := s.content
	defer s.mu.RUnlock()
	// Open the file atomically, creating it if it does not exist
	f, err := os.OpenFile(storageFile, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o644)
	if err != nil {
		if os.IsExist(err) {
			// File exists, open for truncation
			f, err = os.OpenFile(storageFile, os.O_WRONLY|os.O_TRUNC, 0o644)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	defer f.Close()
	// Write content
	if _, err := f.WriteString(content); err != nil {
		return err
	}
	return nil
}

// Initialize implements the mcp.Server interface's Initialize method, providing server metadata.
func (s *TodoServer) Initialize(_ context.Context, req *mcp.InitializeRequest) (*mcp.InitializeResult, error) {
	// Provide basic server info and capabilities.
	// Log server initialization
	s.logf("Server initialized with protocol version %s", "2025-06-18")
	return &mcp.InitializeResult{
		ProtocolVersion: "2025-06-18",
		ServerInfo:      &mcp.Implementation{Name: "todo", Version: "1.0.0"},
		Instructions: `Task Management

Use the todo_read and todo_write tools to manage work items. Use todo_read to retrieve the current list, then use todo_write to save an updated list. This approach supports multi-step tasks, cross-file work, notes, and ambiguous scopes.

Workflow:
- Begin: read current list
- Select: pick the next incomplete task
- Update: modify list as needed, removing the completed task
- Finish: write back the revised list and verify completion
When a task is completed, remove its entry from the list using todo_write.

Important:
- Always read first; writing without reading may lose data
- todo_write replaces the entire file
- NEVER view or edit the TODO.md file with other tools; only use todo_read and todo_write to update TODO.md

Guidelines:
- Keep entries concise and actionable.
- Do not skip the tools for complex tasks.

Template:
- [ ] Implement feature X
  - [ ] Update API
  - [ ] Write tests
  - [ ] Run tests
  - [ ] Run lint
- [ ] Blocked: waiting on credentials`,
		Capabilities: &mcp.ServerCapabilities{
			Tools: &mcp.ToolCapabilities{ListChanged: true},
		},
	}, nil
}

// ListTools returns the list of tools supported by the server.
func (s *TodoServer) ListTools(_ context.Context, req *mcp.ListToolsRequest) (*mcp.ListToolsResult, error) {
	s.logf("ListTools called")
	tools := getTools()
	// Convert []mcp.Tool to []*mcp.Tool
	ptrs := make([]*mcp.Tool, len(tools))
	for i := range tools {
		ptrs[i] = &tools[i]
	}
	return &mcp.ListToolsResult{Tools: ptrs}, nil
}

// CallTool dispatches a tool call to the appropriate handler.
func (s *TodoServer) CallTool(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	switch req.Params.Name {
	case "todo_read":
		s.logf("todo_read called")

		res, err := s.handleRead(ctx, req)
		return res, err
	case "todo_write":
		var args struct {
			Content string `json:"content"`
		}
		if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
			return nil, fmt.Errorf("failed to parse arguments: %w", err)
		}
		res, err := s.handleWrite(ctx, req, args)
		return res, err
	default:
		return nil, fmt.Errorf("unknown tool %s", req.Params.Name)
	}
}

// handleRead returns the current todo content.
func (s *TodoServer) handleRead(_ context.Context, _ *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Reload content from storage file to ensure up-to-date view
	if err := s.loadFromFile(); err != nil {
		// Log error but continue with existing cached content
		s.logf("failed to reload todo file: %v", err)
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	content := s.content
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: content}}}, nil
}

// handleWrite updates the todo content and persists it.
func (s *TodoServer) handleWrite(_ context.Context, req *mcp.CallToolRequest, args struct {
	Content string `json:"content"`
}) (*mcp.CallToolResult, error) {
	// Validate request
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}
	s.logf("todo_write called with %d chars", len(args.Content))
	if len(args.Content) == 0 {
		// Clear todo entry (ignore session)
		// Truncate storage file to clear persisted data
		if err := os.WriteFile(storageFile, []byte{}, 0o644); err != nil {
			s.logger.Printf("failed to clear storage file: %v", err)
			return nil, fmt.Errorf("failed to clear storage: %w", err)
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "Cleared todo list"}}}, nil
	}
	// Update content and persist
	s.mu.Lock()
	s.content = args.Content
	s.mu.Unlock()
	if err := s.persistToFile(); err != nil {
		return nil, fmt.Errorf("failed to persist todos: %w", err)
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Updated (%d chars)", len(args.Content))}}}, nil
}

// getTools constructs and returns the tool definitions for the server.
func getTools() []mcp.Tool {
	return getToolDefinitions()
}
