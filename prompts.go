package main

import "github.com/modelcontextprotocol/go-sdk/mcp"

func getTodoTaskPrompt() mcp.Prompt {
	return mcp.Prompt{
		Name:        "todo-task",
		Title:       "Process next todo item",
		Description: "Automatically process the next pending todo item from the TODO list",
	}
}

func getTodoPromptMessage() *mcp.PromptMessage {
	return &mcp.PromptMessage{
		Role:    "user",
		Content: &mcp.TextContent{Text: "Review the TODO list below and implement the first incomplete task. When complete, remove it from the list by reading the current state with todo_read, then updating with todo_write. Feel free to add subtasks if needed."},
	}
}

func getTodoPromptMessageWithContent(content string) *mcp.PromptMessage {
	base := getTodoPromptMessage()
	if txt, ok := base.Content.(*mcp.TextContent); ok {
		txt.Text += "\n\nCurrent TODO list:\n" + content
	}
	return base
}

func getTodoPrompts() []mcp.Prompt {
	return []mcp.Prompt{getTodoTaskPrompt()}
}
