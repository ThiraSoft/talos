package agentic

import (
	"google.golang.org/genai"
)

type Agent struct {
	Name        string
	Description string
	Role        string

	Model       string
	Temperature float32

	Provider      Provider
	ChatSession   *genai.Chat
	History       []*genai.Content
	Configuration *genai.GenerateContentConfig
	Tools         []*genai.Tool
}

func NewAgent(name, desc, instructions string) *Agent {
	na := Agent{
		Name:        name,
		Description: desc,
		Provider:    GOOGLE,
		// Model:       "gemini-2.5-flash-preview-04-17",
		// Model:       "gemini-2.5-flash-preview-05-20",
		Model:         "gemini-2.0-flash-lite",
		History:       []*genai.Content{},
		ChatSession:   &genai.Chat{},
		Configuration: &genai.GenerateContentConfig{},
		Temperature:   float32(1.0),
		Tools:         []*genai.Tool{},
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
		na.History,
	)

	return &na
}

func (a *Agent) AddTools(tool ...*genai.Tool) {
	for _, t := range tool {
		a.Tools = append(a.Tools, t)
		a.Configuration.Tools = append(a.Configuration.Tools, t)
	}
}

// UpdateInstructions updates the agent's instructions and the system instruction in the configuration.
func (a *Agent) GetInstructions() string {
	return a.Configuration.SystemInstruction.Parts[0].Text
}

func (a *Agent) SetInstructions(newInstructions string) {
	a.Configuration.SystemInstruction = genai.NewContentFromText(newInstructions, genai.RoleModel)
}
