package main

import "github.com/modelcontextprotocol/go-sdk/mcp"

func getTodoTaskPrompt() mcp.Prompt {
	return mcp.Prompt{
		Name:        "todo-task",
		Title:       "todo-task",
		Description: "Do first todo task",
	}
}

func getTodoPromptMessage() *mcp.PromptMessage {
	return &mcp.PromptMessage{
		Role:    "user",
		Content: &mcp.TextContent{Text: "Read our list of todo tasks with todo_read. Implement the first item. When the first task is complete, remove it from the todo list by re-reading the list with todo_read then writing the new list with todo_write. If it is found that more tasks need to be added to the todo list, then add them."},
	}
}

func getTodoPrompts() []mcp.Prompt {
	return []mcp.Prompt{getTodoTaskPrompt()}
}
