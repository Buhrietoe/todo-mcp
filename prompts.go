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
		Content: &mcp.TextContent{Text: "Review this list of todo tasks below. Implement the first item, taking into consideration the rest of the work that needs to be done. When the first task is complete, remove it from the todo list by re-reading the list with todo_read then writing the new list with todo_write. If it is found that more tasks need to be added to the todo list to help split up work, then add them."},
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
