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
	defaultStorageFile = "TODO.md"
)

// getStorageFile returns the storage file path from environment or default
func getStorageFile() string {
	if file := os.Getenv("TODO_FILE"); file != "" {
		return file
	}
	return defaultStorageFile
}

// Unified server instructions for consistency across the application
const serverInstructions = `Task Management with TODO MCP

Use todo_read and todo_write tools to manage your work items efficiently. Always read first, then modify as needed, and finally write back the updated list.

Workflow:
1. Read: Use todo_read to retrieve the current TODO list
2. Select: Choose the next incomplete task to work on
3. Update: Modify the list as needed, removing completed tasks
4. Write: Use todo_write to save the revised list

IMPORTANT: todo_write replaces the entire file. Always read first; writing without reading may lose data.
IMPORTANT: NEVER view or edit the TODO.md file with other tools; only use todo_read and todo_write.

Guidelines:
- Keep entries concise and actionable
- Break complex tasks into smaller subtasks when needed
- Remove completed tasks promptly to avoid confusion

Template:
- [ ] Implement feature X
  - [ ] Update API
  - [ ] Write tests
  - [ ] Run tests
  - [ ] Run lint
- [ ] Blocked: waiting on credentials`

// TodoServer implements the mcp.Server interface and provides todo management functionality.
type TodoServer struct {
	mu       sync.RWMutex
	content  string
	fallback string
	version  string
	logger   *log.Logger
}

func (s *TodoServer) logf(format string, args ...any) { s.logger.Printf(format, args...) }

// loadFromFile loads persisted todos from storageFile if it exists.
func (s *TodoServer) loadFromFile() error {
	storageFile := getStorageFile()
	data, err := os.ReadFile(storageFile)
	if err != nil {
		if os.IsNotExist(err) {
			s.logf("Storage file %s does not exist, starting with empty content", storageFile)
			return nil
		}
		return fmt.Errorf("failed to read todo file %s: %w", storageFile, err)
	}
	if len(data) == 0 {
		s.logf("Storage file %s is empty", storageFile)
		return nil
	}

	s.mu.Lock()
	s.content = string(data)
	defer s.mu.Unlock()
	s.logf("Loaded %d characters from %s", len(data), storageFile)
	return nil
}

// persistToFile writes the current todo content to storageFile.
func (s *TodoServer) persistToFile() error {
	s.mu.RLock()
	content := s.content
	s.mu.RUnlock()
	
	storageFile := getStorageFile()
	if err := os.WriteFile(storageFile, []byte(content), 0o644); err != nil {
		return fmt.Errorf("failed to write todo file %s: %w", storageFile, err)
	}
	s.logf("Saved %d characters to %s", len(content), storageFile)
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
		Instructions:    serverInstructions,
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
		return nil, fmt.Errorf("invalid tool request: request cannot be nil")
	}
	s.logf("todo_write called with %d characters", len(args.Content))
	if len(args.Content) == 0 {
		// Clear todo entry (ignore session)
		// Truncate storage file to clear persisted data
		storageFile := getStorageFile()
		if err := os.WriteFile(storageFile, []byte{}, 0o644); err != nil {
			return nil, fmt.Errorf("failed to clear todo file %s: %w", storageFile, err)
		}
		s.logf("Cleared todo list in %s", storageFile)
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "Cleared todo list"}}}, nil
	}
	// Update content and persist
	s.mu.Lock()
	s.content = args.Content
	s.mu.Unlock()
	if err := s.persistToFile(); err != nil {
		return nil, fmt.Errorf("failed to persist todo updates: %w", err)
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Updated todo list (%d characters)", len(args.Content))}}}, nil
}

// getTools constructs and returns the tool definitions for the server.
func getTools() []mcp.Tool {
	return getToolDefinitions()
}
