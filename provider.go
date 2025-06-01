package talos

type (
	Provider string
)

const (
	PROVIDER_OPEN_AI   Provider = "OPEN_AI"
	PROVIDER_ANTHROPIC Provider = "ANTHROPIC"
	PROVIDER_GOOGLE    Provider = "GOOGLE"
	PROVIDER_MISTRAL   Provider = "MISTRAL"
)

func IsValidProvider(p Provider) bool {
	switch p {
	case PROVIDER_OPEN_AI, PROVIDER_ANTHROPIC, PROVIDER_GOOGLE, PROVIDER_MISTRAL:
		return true
	default:
		return false
	}
}
