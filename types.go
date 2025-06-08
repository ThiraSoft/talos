package talos

import (
	"context"

	"google.golang.org/genai"
)

type (
	DebugLevel string
	Provider   string
)

type Agent struct {
	Name             string
	Description      string
	Role             string
	Model            string
	Temperature      float32
	Provider         Provider
	ChatSession      *genai.Chat
	Configuration    *genai.GenerateContentConfig
	PartsBuffer      []*genai.Part                                               // For tools responses
	CallToolFunction func(caller *Agent, fn *genai.FunctionCall) (string, error) // Function to call tool functions
	OutputChan       chan AgentNotification
	InputChan        chan AgentNotification
	Ctx              context.Context    // Context for API calls, can be used to set timeouts or other options
	CtxCancelFunc    context.CancelFunc // Function to cancel the context
}

type AgentNotification struct {
	AgentName      string
	MessageContent string
	MessageType    string
	Bytes          []byte // For audio messages, can be nil if not applicable
}

type Flow struct {
	Id          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Tasks       []*Task          `json:"tasks"`
	Agents      []*Agent         `json:"agents"`  // Agents involved in the flow
	History     []*genai.Content `json:"history"` // History of the flow
}
