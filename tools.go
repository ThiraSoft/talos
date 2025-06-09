package talos

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"google.golang.org/genai"
)

func (a *Agent) CallTool(fn *genai.FunctionCall) (string, error) {
	// Custom function to call tool functions first, if provided.
	if a.CallToolFunction != nil {
		resp, err := a.CallToolFunction(a, fn)
		return resp, err
	}

	// If no custom function is provided, use the default tool handler.
	if fn.Name == "send_message" {
		resp, err := SendMessage(fn)
		if err != nil {
			return "", fmt.Errorf("error calling send_message tool: %w", err)
		}
		return resp, nil
	}

	if fn.Name == "write_file" {
		resp, err := WriteFile(fn)
		if err != nil {
			return "", fmt.Errorf("error calling write_file tool: %w", err)
		}
		return resp, nil
	}

	return "unknown tool function", fmt.Errorf("unknown tool function: %s", fn.Name)
}

// Tool_Definition_SendMessage is the definition of the send_message tool function.
var Tool_Definition_SendMessage *genai.FunctionDeclaration = &genai.FunctionDeclaration{
	Name:        "send_message",
	Description: "Allow you to send a message to someone.",
	Parameters: &genai.Schema{
		Type: "object",
		Properties: map[string]*genai.Schema{
			"from": {
				Type:        "string",
				Description: "The name of the sender.",
			},
			"to": {
				Type:        "string",
				Description: "The name of the receiver.",
			},
			"message": {
				Type:        "string",
				Description: "The message to send.",
			},
		},
		Required: []string{"from", "to", "message"},
	},
}

// SendMessage sends a message to another agent.
func SendMessage(tool *genai.FunctionCall) (string, error) {
	message := tool.Args["message"].(string)
	// from := tool.Args["from"].(string)
	to := tool.Args["to"].(string)
	logger.Debug("SendMessage called", slog.String("to", to), slog.String("message", message))

	var response string = ""
	var err error

	for _, agent := range Agents { // Global Agents is updated by the flow when it starts
		if agent.Name == to {
			logger.Debug("Found agent", slog.String("agent_name", agent.Name))
			response, err = agent.ChatWithRetry(message, 5)
			if err != nil {
				return "Error while asking " + to + ": " + err.Error(), fmt.Errorf("error while asking %s : %w", to, err)
			}
			return "Response from " + agent.Name + " : " + response, nil
		}
	}
	return response, nil
}

// Tool_Definition_WriteFile is the definition of the write_file tool function.
var Tool_Definition_WriteFile *genai.FunctionDeclaration = &genai.FunctionDeclaration{
	Name:        "write_file",
	Description: "Write a file given a file_name and a content.",
	Parameters: &genai.Schema{
		Type: "object",
		Properties: map[string]*genai.Schema{
			"file_name": {
				Type:        "string",
				Description: "The name of the file to create",
			},
			"content": {
				Type:        "string",
				Description: "The string content of the file to create.",
			},
		},
		Required: []string{"file_name", "content"},
	},
}

// WriteFile writes a file with the given file_name and content.
func WriteFile(tool *genai.FunctionCall) (string, error) {
	// Validation
	filename, ok := tool.Args["file_name"].(string)
	if !ok {
		return "Invalid arguments 'file_name' for write_file tool", fmt.Errorf("invalid arguments 'file_name' for write_file tool")
	}
	content, ok := tool.Args["content"].(string)
	if !ok {
		return "Invalid arguments 'content' for write_file tool", fmt.Errorf("invalid arguments 'content' for write_file tool")
	}

	// Create the directory if it does not exist
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		logger.Error("Error creating directory", slog.String("directory", filepath.Dir(filename)), slog.String("error", err.Error()))
		return "Error creating directory: " + err.Error(), err
	}

	// Convert the message to bytes
	byteContent := []byte(content)

	// Mehtod to write a file with the content
	err := os.WriteFile(filename, byteContent, 0644)
	if err != nil {
		return "Error writing file: " + err.Error(), err
	}

	return "Page written", nil
}
