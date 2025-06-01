package talos

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"google.golang.org/genai"
)

// ChatWithRetry is a wrapper around the Chat function that retries on specific errors
func (a *Agent) ChatWithRetry(input string, maxRetries int) (string, error) {
	var lastErr error
	var count int = 1

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
		log.Printf("Tentative %d échouée : %v", attempt+1, err)

		// Conditions spécifiques de retry
		if isRetryableError(err) {
			count++
			continue
		}

		// Arrêter si l'erreur n'est pas retraitable
		break
	}

	return "", fmt.Errorf("échec après %d tentatives : %w", count, lastErr)
}

// Fonction pour déterminer si l'erreur justifie un retry
func isRetryableError(err error) bool {
	// Exemples de conditions de retry
	return errors.Is(err, context.DeadlineExceeded) ||
		strings.Contains(err.Error(), "timeout") ||
		strings.Contains(err.Error(), "connection reset") ||
		strings.Contains(err.Error(), "503") ||
		strings.Contains(err.Error(), "429")
}

func (a *Agent) Chat(input string) (string, error) {
	fullResponse := a.Name + " : "
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
	fmt.Println("\n======================")
	fmt.Println(" " + a.Name + " : ")
	fmt.Println("======================")
	if err != nil {
		fmt.Println("Error receiving response:", err)
	}

	if res == nil {
		fmt.Println("Received nil chunk, skipping...")
		return "", fmt.Errorf("received nil response from chat session")
	}
	if len(res.Candidates) == 0 {
		fmt.Println("No candidates in chunk, skipping...")
		return "", fmt.Errorf("no candidates in response from chat session")
	}
	if res.Candidates[0].Content == nil {
		fmt.Println("No content in candidate, skipping...")
		return "", fmt.Errorf("no content in candidate from chat session")
	}

	fmt.Println("PARTS : ", len(res.Candidates[0].Content.Parts))
	for _, part := range res.Candidates[0].Content.Parts {
		fmt.Println("Part text: ", part.Text)
		if part.FunctionCall != nil {
			fmt.Println("Part FunctionCall: ", part.FunctionCall.Name, part.FunctionCall.Args)
		}
	}

	for _, p := range res.Candidates[0].Content.Parts {
		fullResponse += p.Text
		responseHandler(p)

		toolResponse, err := a.toolHandler(p)
		if err != nil {
			return "", fmt.Errorf("error handling tool response: %w", err)
		}

		fullResponse += toolResponse
	}

	return fullResponse, nil
}

func responseHandler(part *genai.Part) (string, error) {
	response := fmt.Sprint(part.Text)
	// fmt.Print(response)

	return response, nil
}

func (a *Agent) toolHandler(part *genai.Part) (string, error) {
	// Vérifier s'il y a un appel de fonction
	if part.FunctionCall != nil {
		fn := part.FunctionCall
		resp, err := CallTool(fn)
		if err != nil {
			fmt.Print("Erreur lors de l'utilisation du tool : ", err)
			return "", fmt.Errorf("error calling tool %s: %w", fn.Name, err)
		}

		// Add the response to the agent's history
		// a.ChatSession.AppendToHistory(
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
