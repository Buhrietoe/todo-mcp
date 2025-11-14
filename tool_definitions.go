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
			Description: `Read the current temporary storage content (in markdown format).

This tool returns the stored data as a string, allowing the LLM to retrieve previously saved information. It can be used for simple state persistence across calls.

The tool will return an error if the storage cannot be accessed.`,
			Title:        "Read TODO file",
			InputSchema:  map[string]any{"type": "object"},
			OutputSchema: map[string]any{"type": "object"},
			Annotations:  readAnnotations,
		},
		{
			Name: "todo_write",
			Description: `Write or overwrite the entire TODO file content (in markdown format).

This tool replaces the whole TODO file with the supplied string, allowing the LLM to store arbitrary data persistently. It can be used to save updated task lists, notes, or any structured information the model wishes to retain across calls.

WARNING: This operation overwrites the entire file. Ensure the provided content includes all data you wish to keep, as any existing content not included will be lost.

The tool will create the TODO file if it does not exist, or overwrite it if it does. It returns an error if the file cannot be written due to permissions or other I/O issues.`,
			Title:        "Write TODO file",
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
