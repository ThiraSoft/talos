package agentic

type (
	Provider string
)

const (
	OPEN_AI   Provider = "OPEN_AI"
	ANTHROPIC Provider = "ANTHROPIC"
	GOOGLE    Provider = "GOOGLE"
	MISTRAL   Provider = "MISTRAL"
)

func IsValidProvider(p Provider) bool {
	switch p {
	case OPEN_AI, ANTHROPIC, GOOGLE, MISTRAL:
		return true
	default:
		return false
	}
}
