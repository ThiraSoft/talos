package talos

import (
	"context"

	"google.golang.org/genai"
)

var (
	// DEFAULT_MODEL    string   = "gemini-2.0-flash-lite"
	DEFAULT_MODEL    string   = "gemini-2.5-flash-preview-05-20"
	DEFAULT_PROVIDER Provider = PROVIDER_GOOGLE // Default provider to use if not specified
)

// OutputNotification sends a notification to the agent's output channel.
// It includes the agent's name, the message content, and the type of message (e.g., text or audio).
// This function is used send responses or notifications from the agent to the output channel.
func (a *Agent) OutputNotification(messageContent, messageType string) {
	a.OutputChan <- AgentNotification{
		AgentName:      a.Name,
		MessageContent: messageContent,
		MessageType:    messageType,
	}
}

// Async starts a goroutine that listens for input messages and processes them asynchronously.
// It handles both text and audio messages, calling the appropriate chat methods.
// The input messages are expected to be of type AgentNotification, which contains the message content and type.
// This allows the agent to handle multiple messages concurrently without blocking the main execution flow.
func (a *Agent) Async() {
	go func() {
		for {
			select {
			case input := <-a.InputChan:
				switch input.MessageType {
				case "TEXT":
					a.ChatWithRetry(input.MessageContent, 5)
				case "AUDIO":
					a.ChatWithRetryWithAudio(input.Bytes, 5)
				}
			case <-a.Ctx.Done():
				return
			}
		}
	}()
}

func NewAgent(
	name, desc, instructions string,
	provider Provider,
	model string,
) *Agent {
	ctx, ctxCancelFunc := context.WithCancel(context.Background())

	na := Agent{
		Name:             name,
		Description:      desc,
		Provider:         provider,
		Model:            model,
		ChatSession:      &genai.Chat{},
		Configuration:    &genai.GenerateContentConfig{},
		Temperature:      float32(1.0),
		PartsBuffer:      make([]*genai.Part, 0, 1000), // For tools responses
		CallToolFunction: nil,
		OutputChan:       make(chan AgentNotification, 255), // Buffered channel for responses
		Ctx:              ctx,
		CtxCancelFunc:    ctxCancelFunc,
	}

	baseInstructions := "You are an AI agent named " + na.Name + ".\n"
	na.Configuration.SystemInstruction = genai.NewContentFromText(baseInstructions+instructions, genai.RoleModel)
	na.Configuration.Temperature = &na.Temperature
	na.Configuration.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockThresholdBlockNone,
		},

		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockThresholdBlockNone,
		},

		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockThresholdBlockNone,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockThresholdBlockNone,
		},
		{
			Category:  genai.HarmCategoryCivicIntegrity,
			Threshold: genai.HarmBlockThresholdBlockNone,
		},
	}

	// init chat session
	na.ChatSession, _ = Client.Chats.Create(
		Ctx,
		na.Model,
		na.Configuration,
		[]*genai.Content{},
	)

	return &na
}

func (a *Agent) AddFunctionDeclarations(declarations ...*genai.FunctionDeclaration) {
	if len(a.Configuration.Tools) == 0 {
		a.Configuration.Tools = append(a.Configuration.Tools, &genai.Tool{})
	} else if a.Configuration.Tools[0] == nil {
		a.Configuration.Tools[0] = &genai.Tool{}
	}

	a.Configuration.Tools[0].FunctionDeclarations = append(a.Configuration.Tools[0].FunctionDeclarations, declarations...)
}

// UpdateInstructions updates the agent's instructions and the system instruction in the configuration.
func (a *Agent) GetInstructions() string {
	return a.Configuration.SystemInstruction.Parts[0].Text
}

func (a *Agent) SetInstructions(newInstructions string) {
	a.Configuration.SystemInstruction = genai.NewContentFromText(newInstructions, genai.RoleModel)
}
