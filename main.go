package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    
    "github.com/modelcontextprotocol/go-sdk/mcp"
    "sync"
)

type TodoServer struct {
    mu       sync.RWMutex
    todos    map[string]string
    fallback string
    logger   *log.Logger
}

func (s *TodoServer) Initialize() mcp.Metadata {
    return mcp.Metadata{
        ProtocolVersion: "1.0.0",
        ServerName:      "todo-server",
        Version:         "0.1.0",
        Instructions:    "Use session-based state management to store todos. Each session has its own todo list, with a fallback for unauthenticated access.",
        Capabilities:    []string{"read", "write"},
    }
}

func (s *TodoServer) GetTools() []mcp.Tool {
    return []mcp.Tool{
        {
            Name:      "todo_read",
            Description: "Read the current todo content for the session. Returns the content if session exists, otherwise returns fallback content.",
            Parameters: []mcp.Parameter{ {
                Name: "session_id", 
                Description: "The session ID to read from. If not provided, uses fallback content.",
                Type: "string",
                Required: false,
            }},
        },
        {
            Name:      "todo_write",
            Description: "Write or overwrite the todo content for the session. If the session ID doesn't exist, it will be created.",
            Parameters: []mcp.Parameter{ {
                Name: "session_id", 
                Description: "The session ID to write to. If not provided, the session ID will be generated.",
                Type: "string",
                Required: false,
            }, {
                Name: "content", 
                Description: "The content to write. Must be less than 1000 characters.",
                Type: "string",
                Required: true,
            }},
        },
    }
}

func (s *TodoServer) HandleRead(ctx *mcp.Context) (string, error) {
    sessionID := ctx.Params["session_id"]
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    if sessionID == "" {
        return s.fallback, nil
    }
    
    content, exists := s.todos[sessionID]
    if !exists {
        return s.fallback, nil
    }
    
    return content, nil
}

func (s *TodoServer) HandleWrite(ctx *mcp.Context) (string, error) {
    sessionID := ctx.Params["session_id"]
    content := ctx.Params["content"]
    
    if len(content) > 1000 {
        return "", &mcp.Error{Message: "Content exceeds 1000 character limit"}
    }
    
    s.mu.Lock()
    defer s.mu.Unlock()
    
    if sessionID == "" {
        sessionID = "default"
    }
    
    s.todos[sessionID] = content
    
    return "Success", nil
}

func main() {
    server := &TodoServer{
        todos:    make(map[string]string),
        fallback: "No todos found.",
        logger:   log.New(os.Stderr, "", 0),
    }
    
    // Create and run the MCP server
    mcpServer := mcp.NewServer(server, &mcp.ServerOptions{Instructions: "Use the todo_read and todo_write tools to manage your todo list."})
    
    // Run server over Stdio with logging
    t := &mcp.LoggingTransport{Transport: &mcp.StdioTransport{}, Writer: os.Stderr}
    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
    defer cancel()
    
    if err := mcpServer.Run(ctx, t); err != nil {
        log.Fatal(err)
    }
}
