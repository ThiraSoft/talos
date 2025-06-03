package talos

import (
	"google.golang.org/genai"
)

var (
	// DEFAULT_MODEL string = "gemini-2.5-flash-preview-04-17"
	// DEFAULT_MODEL string = "gemini-2.0-flash-lite"
	// DEFAULT_MODEL    string   = "gemini-2.0-flash"
	DEFAULT_MODEL    string   = "gemini-2.5-flash-preview-05-20"
	DEFAULT_PROVIDER Provider = PROVIDER_GOOGLE // Default provider to use if not specified
)

type Agent struct {
	Name        string
	Description string
	Role        string

	Model       string
	Temperature float32

	Provider         Provider
	ChatSession      *genai.Chat
	History          []*genai.Content
	Configuration    *genai.GenerateContentConfig
	Tools            []*genai.Tool
	PartsBuffer      []*genai.Part                                               // For tools responses
	CallToolFunction func(caller *Agent, fn *genai.FunctionCall) (string, error) // Function to call tool functions
}

func NewAgent(name, desc, instructions string, provider Provider, model string) *Agent {
	na := Agent{
		Name:             name,
		Description:      desc,
		Provider:         provider,
		Model:            model,
		History:          make([]*genai.Content, 0, 10000),
		ChatSession:      &genai.Chat{},
		Configuration:    &genai.GenerateContentConfig{},
		Temperature:      float32(1.0),
		Tools:            make([]*genai.Tool, 0, 255),
		PartsBuffer:      make([]*genai.Part, 0, 10000), // For tools responses
		CallToolFunction: nil,
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
	na.Configuration.Tools = na.Tools

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
	a.Tools = append(a.Tools, tool...)
	a.Configuration.Tools = append(a.Configuration.Tools, tool...)
}

// UpdateInstructions updates the agent's instructions and the system instruction in the configuration.
func (a *Agent) GetInstructions() string {
	return a.Configuration.SystemInstruction.Parts[0].Text
}

func (a *Agent) SetInstructions(newInstructions string) {
	a.Configuration.SystemInstruction = genai.NewContentFromText(newInstructions, genai.RoleModel)
}
