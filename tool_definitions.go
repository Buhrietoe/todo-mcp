package main

import "github.com/modelcontextprotocol/go-sdk/mcp"

func getToolDefinitions() []mcp.Tool {
	// Define tool annotations and schemas
	openWorldHint := false
	readAnnotations := &mcp.ToolAnnotations{ReadOnlyHint: true, IdempotentHint: true, OpenWorldHint: &openWorldHint}
	writeDestructive := true
	writeAnnotations := &mcp.ToolAnnotations{DestructiveHint: &writeDestructive, IdempotentHint: false, OpenWorldHint: &openWorldHint}

	// Input schema for write tool
	writeInputSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"content": map[string]any{
				"type":        "string",
				"description": "The TODO list content to save",
			},
		},
		"required": []string{"content"},
	}

	return []mcp.Tool{
		{
			Name: "todo_read",
			Description: `Read the current TODO list from storage.
Returns the complete TODO content as markdown text. Use this to retrieve the current state before making any modifications.`,
			Title:        "Read TODO List",
			InputSchema:  map[string]any{"type": "object"},
			OutputSchema: map[string]any{"type": "object"},
			Annotations:  readAnnotations,
		},
		{
			Name: "todo_write",
			Description: `Write the complete TODO list to storage.
Replaces the entire TODO content with the provided markdown text. Always use todo_read first to avoid data loss.`,
			Title:        "Write TODO List",
			InputSchema:  writeInputSchema,
			OutputSchema: map[string]any{"type": "object"},
			Annotations:  writeAnnotations,
		},
		{
			Name:         "prompt",
			Description:  `List available prompt names.`,
			Title:        "List prompts",
			InputSchema:  map[string]any{"type": "object"},
			OutputSchema: map[string]any{"type": "object"},
			Annotations:  readAnnotations,
		},
	}
}
