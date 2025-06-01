package agentic

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
	fullResponse := ""
	cs := a.ChatSession
	// stream := cs.SendMessageStream(Ctx, genai.Part{Text: input})

	res, err := cs.SendMessage(Ctx, genai.Part{Text: input})
	fmt.Println("\n" + a.Name + " : \n")
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

	for _, p := range res.Candidates[0].Content.Parts {
		fullResponse += p.Text
		responseHandler(p)
		a.toolHandler(p)
	}

	// fmt.Println("\n" + a.Name + " : \n")
	// for chunk, err := range stream {
	// 	if err != nil {
	// 		fmt.Println("Error receiving response:", err)
	// 	}
	//
	// 	// Log the chunk for debugging
	// 	if chunk == nil {
	// 		fmt.Println("Received nil chunk, skipping...")
	// 		continue
	// 	}
	// 	if len(chunk.Candidates) == 0 {
	// 		fmt.Println("No candidates in chunk, skipping...")
	// 		continue // Skip if no candidates in the chunk
	// 	}
	// 	if chunk.Candidates[0].Content == nil {
	// 		fmt.Println("No content in candidate, skipping...")
	// 		continue // Skip if no parts in the candidate content
	// 	}
	//
	// 	part := chunk.Candidates[0].Content.Parts[0]
	// 	responseHandler(part)
	// 	fullResponse += part.Text
	//
	// 	a.toolHandler(part)
	// }

	return fullResponse, nil
}

func responseHandler(part *genai.Part) (string, error) {
	response := fmt.Sprint(part.Text)
	fmt.Print(response)

	return response, nil
}

func (a *Agent) toolHandler(part *genai.Part) {
	// Vérifier s'il y a un appel de fonction
	if part.FunctionCall != nil {
		// Si c'est un appel de fonction, on l'affiche
		fmt.Printf("Tool call detected\nTool : %s\nArgs : %s\n", part.FunctionCall.Name, part.FunctionCall.Args)

		// Si c'est un appel de fonction, on l'exécute
		fn := part.FunctionCall
		resp, err := CallTool(fn)
		if err != nil {
			fmt.Print("Erreur lors de l'utilisation du tool : ", err)
		}

		// Appeler le chat avec la réponse de la fonction
		_, err = a.ChatWithRetry(resp, 5)
		if err != nil {
			fmt.Print("Erreur lors de l'appel à l'API : ", err)
			return
		}
	}
}
