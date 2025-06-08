package talos

import (
	"fmt"

	"google.golang.org/genai"
)

// SetHistory sets the chat history for the agent.
// It takes a slice of genai.Content and initializes the ChatSession.
// If the history is nil, it initializes an empty history.
func (a *Agent) SetHistory(history []*genai.Content) {
	var h []*genai.Content

	if history == nil {
		h = []*genai.Content{}
	} else {
		h = history
	}

	a.ChatSession, _ = Client.Chats.Create(
		Ctx,
		a.Model,
		a.Configuration,
		h,
	)

	// Log the history set
	logger(fmt.Sprintf("History set for agent %s: %d messages", a.Name, len(history)), DEBUG_LEVEL_ALL)
}

// AddTextToHistory adds a text message to the agent's history.
// It creates a new genai.Part with the provided text and appends it to the PartsBuffer.
// This function put the text in the agent's buffer, which will be sent to the chat on next send.
func (a *Agent) AddTextToHistory(text string) {
	part := &genai.Part{
		Text: text,
	}
	a.PartsBuffer = append(a.PartsBuffer, part)
}
