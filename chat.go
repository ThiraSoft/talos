package talos

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"google.golang.org/genai"
)

// ChatWithRetry is a wrapper around the Chat function that retries on specific errors
func (a *Agent) ChatWithRetry(input string, maxRetries int) (string, error) {
	var lastErr error
	count := 1

	// Loop
	for attempt := 0; attempt < maxRetries; attempt++ {
		// Délai *5 entre les tentatives
		if attempt > 0 {
			backoffDuration := time.Duration(5*attempt) * time.Second
			time.Sleep(backoffDuration)
		}

		// Tentative d'appel à Chat
		response, err := a.Chat(input)
		if err == nil {
			return response, nil
		}

		// Stocker la dernière erreur
		lastErr = err

		// Log de l'erreur
		logger(
			fmt.Sprintf("Tentative %d échouée : %v", attempt+1, err),
			DEBUG_LEVEL_ALL,
			DEBUG_LEVEL_ERRORS,
		)

		// Conditions spécifiques de retry
		if isRetryableError(err) {
			count++
			continue
		}

		// Arrêter si l'erreur n'est pas retraitable
		break
	}

	return fmt.Sprintf("Echec de l'appel : %s", lastErr), fmt.Errorf("échec après %d tentatives : %w", count, lastErr)
}

// ChatWithRetry is a wrapper around the Chat function that retries on specific errors
func (a *Agent) ChatWithRetryWithAudio(audioBytes []byte, maxRetries int) (string, error) {
	var lastErr error
	count := 1

	// Loop
	for attempt := 0; attempt < maxRetries; attempt++ {
		// Délai *5 entre les tentatives
		if attempt > 0 {
			backoffDuration := time.Duration(5*attempt) * time.Second
			time.Sleep(backoffDuration)
		}

		// Tentative d'appel à Chat
		response, err := a.ChatWithAudio(audioBytes)
		if err == nil {
			return response, nil
		}

		// Stocker la dernière erreur
		lastErr = err

		// Log de l'erreur
		logger(
			fmt.Sprintf("Tentative %d échouée : %v", attempt+1, err),
			DEBUG_LEVEL_ALL,
			DEBUG_LEVEL_ERRORS,
		)

		// Conditions spécifiques de retry
		if isRetryableError(err) {
			count++
			continue
		}

		// Arrêter si l'erreur n'est pas retraitable
		break
	}

	return fmt.Sprintf("Echec de l'appel : %s", lastErr), fmt.Errorf("échec après %d tentatives : %w", count, lastErr)
}

// Fonction pour déterminer si l'erreur justifie un retry
func isRetryableError(err error) bool {
	// Exemples de conditions de retry
	return errors.Is(err, context.DeadlineExceeded) ||
		strings.Contains(err.Error(), "timeout") ||
		strings.Contains(err.Error(), "connection reset") ||
		strings.Contains(err.Error(), "500") ||
		strings.Contains(err.Error(), "503") ||
		strings.Contains(err.Error(), "429")
}

func (a *Agent) Chat(input string) (string, error) {
	fullResponse := ""
	cs := a.ChatSession

	// for each part in the buffer, append it to the parts slice
	parts := []genai.Part{}
	if len(a.PartsBuffer) != 0 {
		for _, content := range a.PartsBuffer {
			parts = append(parts, *content)
		}
		a.PartsBuffer = make([]*genai.Part, 0, 10000) // Reset buffer after sending
	}
	parts = append(parts, genai.Part{Text: input})

	// send the message to the chat session
	res, err := cs.SendMessage(Ctx, parts...)

	// Display the agent's name
	logger(
		a.Name+" : ",
		DEBUG_LEVEL_ALL,
		DEBUG_LEVEL_ERRORS,
	)
	if err != nil {
		logger(
			fmt.Sprintf("Error receiving response: %s", err),
			DEBUG_LEVEL_ALL,
		)
		return fmt.Sprintf("error receiving response from chat session : %s", err), fmt.Errorf("error receiving response from chat session: %w", err)
	}

	if res == nil {
		logger("Received nil chunk, skipping...", DEBUG_LEVEL_ALL, DEBUG_LEVEL_ERRORS)
		return "", fmt.Errorf("received nil response from chat session")
	}
	if len(res.Candidates) == 0 {
		logger("No candidates in chunk, skipping...", DEBUG_LEVEL_ALL, DEBUG_LEVEL_ERRORS)
		return "", fmt.Errorf("no candidates in response from chat session")
	}
	if res.Candidates[0].Content == nil {
		logger("No content in candidate, skipping...", DEBUG_LEVEL_ALL, DEBUG_LEVEL_ERRORS)
		return "", fmt.Errorf("no content in candidate from chat session")
	}

	logger("PARTS : "+fmt.Sprint(len(res.Candidates[0].Content.Parts)), DEBUG_LEVEL_ALL)
	if len(res.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no parts in response from chat session")
	}

	for _, part := range res.Candidates[0].Content.Parts {
		logger("Part text: "+part.Text, DEBUG_LEVEL_ALL)
		if part.FunctionCall != nil {
			logger(fmt.Sprintf("Part FunctionCall: %s %v", part.FunctionCall.Name, part.FunctionCall.Args), DEBUG_LEVEL_ALL)
		}
	}

	for _, p := range res.Candidates[0].Content.Parts {
		// Add the response to the agent's responses channel
		if p.Text != "" {
			a.OutputNotification(p.Text, "TEXT")
			fullResponse += p.Text
		}

		toolResponse, err := a.toolHandler(p)
		if err != nil {
			return fmt.Sprintf("error handling tool response: %s", err), fmt.Errorf("error handling tool response: %w", err)
		}

		fullResponse += toolResponse
	}

	return fullResponse, nil
}

func (a *Agent) toolHandler(part *genai.Part) (string, error) {
	// Vérifier s'il y a un appel de fonction
	if part.FunctionCall != nil {
		fn := part.FunctionCall
		resp, err := a.CallTool(fn)
		if err != nil {
			logger("Erreur lors de l'utilisation du tool : "+fmt.Sprintf("Tool : %s \n resp : %s \n err : %s", fn.Name, resp, err), DEBUG_LEVEL_ALL, DEBUG_LEVEL_ERRORS)
		}

		// Add the response to the agent's history
		a.PartsBuffer = append(a.PartsBuffer,
			&genai.Part{
				FunctionResponse: &genai.FunctionResponse{
					Name:     fn.Name,
					Response: map[string]any{"Response": resp},
				},
			},
		)

		return resp, err
	}
	return "", nil
}

func (a *Agent) ChatWithAudio(audioBytes []byte) (string, error) {
	fullResponse := ""
	cs := a.ChatSession

	// for each part in the buffer, append it to the parts slice
	parts := []genai.Part{}
	if len(a.PartsBuffer) != 0 {
		for _, content := range a.PartsBuffer {
			parts = append(parts, *content)
		}
		a.PartsBuffer = make([]*genai.Part, 0, 10000) // Reset buffer after sending
	}

	// Add audio
	newPart := genai.Part{
		InlineData: &genai.Blob{
			MIMEType: "audio/mp3",
			Data:     audioBytes,
		},
	}

	parts = append(parts, newPart)

	// send the message to the chat session
	res, err := cs.SendMessage(Ctx, parts...)

	// Display the agent's name
	logger(
		a.Name+" : ",
		DEBUG_LEVEL_ALL,
		DEBUG_LEVEL_ERRORS,
	)
	if err != nil {
		logger(
			fmt.Sprintf("Error receiving response: %s", err),
			DEBUG_LEVEL_ALL,
		)
		return fmt.Sprintf("error receiving response from chat session : %s", err), fmt.Errorf("error receiving response from chat session: %w", err)
	}

	if res == nil {
		logger("Received nil chunk, skipping...", DEBUG_LEVEL_ALL, DEBUG_LEVEL_ERRORS)
		return "", fmt.Errorf("received nil response from chat session")
	}
	if len(res.Candidates) == 0 {
		logger("No candidates in chunk, skipping...", DEBUG_LEVEL_ALL, DEBUG_LEVEL_ERRORS)
		return "", fmt.Errorf("no candidates in response from chat session")
	}
	if res.Candidates[0].Content == nil {
		logger("No content in candidate, skipping...", DEBUG_LEVEL_ALL, DEBUG_LEVEL_ERRORS)
		return "", fmt.Errorf("no content in candidate from chat session")
	}

	logger("PARTS : "+fmt.Sprint(len(res.Candidates[0].Content.Parts)), DEBUG_LEVEL_ALL)
	if len(res.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no parts in response from chat session")
	}

	for _, part := range res.Candidates[0].Content.Parts {
		logger("Part text: "+part.Text, DEBUG_LEVEL_ALL)
		if part.FunctionCall != nil {
			logger(fmt.Sprintf("Part FunctionCall: %s %v", part.FunctionCall.Name, part.FunctionCall.Args), DEBUG_LEVEL_ALL)
		}
	}

	for _, p := range res.Candidates[0].Content.Parts {
		// Add the response to the agent's responses channel
		if p.Text != "" {
			a.OutputNotification(p.Text, "TEXT")
			fullResponse += p.Text
		}

		toolResponse, err := a.toolHandler(p)
		if err != nil {
			return fmt.Sprintf("error handling tool response: %s", err), fmt.Errorf("error handling tool response: %w", err)
		}

		fullResponse += toolResponse
	}

	return fullResponse, nil
}
