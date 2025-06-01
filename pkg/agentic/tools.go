package agentic

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"google.golang.org/genai"
)

// Tools provides a set of tools that can be used by agents to perform specific tasks.
var DefaultTools = []*genai.Tool{}

func init() {
	DefaultTools = append(
		// DefaultTools is a slice of tools that can be used by agents.
		DefaultTools,

		// Add the tools to the DefaultTools slice
		Tool_Definition_SendMessage(),
		Tool_Definition_WriteFile(),
	)
}

// Placeholder for tool calls
func CallTool(fn *genai.FunctionCall) (string, error) {
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
			return "", fmt.Errorf("error calling write_page tool: %w", err)
		}
		return resp, nil
	}

	return "unknown tool function", fmt.Errorf("unknown tool function: %s", fn.Name)
}

// filterBackticks removes (multiline) all text between triple backticks (including the backticks themselves) from the input string.
// It uses a regular expression to match the pattern of triple backticks and any text in between.
func FilterTools(input string) string {
	re := regexp.MustCompile("(?s)```tool.*?```")
	return strings.TrimSpace(re.ReplaceAllString(input, ""))
}

// FilterCode removes (multiline) all text between triple backticks (including the backticks themselves) from the input string.
func FilterCode(input string) string {
	re := regexp.MustCompile("(?s)```.*?```")
	return re.ReplaceAllString(input, "")
}

func Tool_Definition_SendMessage() *genai.Tool {
	tool := &genai.Tool{
		FunctionDeclarations: []*genai.FunctionDeclaration{
			{
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
			},
		},
	}
	return tool
}

func SendMessage(tool *genai.FunctionCall) (string, error) {
	message := tool.Args["message"].(string)
	// from := tool.Args["from"].(string)
	to := tool.Args["to"].(string)
	fmt.Println("\n Sending message to : ", to, "\n Message : ", message)

	var response string
	var err error

	// question
	for _, agent := range Agents {
		if agent.Name == to {
			fmt.Println("\nFound agent:", agent.Name)
			response, err = agent.ChatWithRetry(message, 5)
			if err != nil {
				return "Error while asking " + to + ": " + err.Error(), fmt.Errorf("Error while asking " + to + ": " + err.Error())
			}
			break
		}
	}
	return response, nil

	// // response
	// if response != "" {
	// 	for _, agent := range Agents {
	// 		if agent.Name == from {
	// 			fmt.Println("\nFound agent:", agent.Name)
	// 			_, err = agent.ChatWithRetry("Response from "+to+" : "+response, 5)
	// 			if err != nil {
	// 				return fmt.Errorf("Error while returning the response from " + to + " to " + from + ": " + err.Error())
	// 			}
	// 		}
	// 	}
	// }

	// return nil
}

func Tool_Definition_WriteFile() *genai.Tool {
	tool := &genai.Tool{
		FunctionDeclarations: []*genai.FunctionDeclaration{
			{
				Name:        "write_file",
				Description: "Write a file with a path/filename and a content.",
				Parameters: &genai.Schema{
					Type: "object",
					Properties: map[string]*genai.Schema{
						"filename": {
							Type:        "string",
							Description: "The relative path/filename of the file to create.",
						},
						"content": {
							Type:        "string",
							Description: "The content of the page to write.",
						},
					},
					Required: []string{"filename", "content"},
				},
			},
		},
	}
	return tool
}

func WriteFile(tool *genai.FunctionCall) (string, error) {
	// Validation
	filename, ok := tool.Args["filename"].(string)
	if !ok {
		return "Invalid arguments 'filename' for write_file tool", fmt.Errorf("Invalid arguments 'filename' for write_file tool")
	}
	content, ok := tool.Args["content"].(string)
	if !ok {
		return "Invalid arguments 'content' for write_file tool", fmt.Errorf("Invalid arguments 'content' for write_file tool")
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
