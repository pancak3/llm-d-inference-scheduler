package requestutil

import (
	"encoding/json"
	"errors"
	"fmt"

	"sigs.k8s.io/gateway-api-inference-extension/pkg/epp/scheduling/types"
)

var (
	errNilRequest        = errors.New("llm request is nil")
	errNilRequestBody    = errors.New("llm request body is nil")
	errUnsupportedInputs = errors.New("llm request body is missing completions or chat completions inputs")
)

// PromptBytes returns a byte representation of the user-provided prompt for the request.
// For completions requests, it returns the raw prompt string bytes. For chat completions requests,
// the slice of message payloads is JSON encoded to retain ordering and content.
func PromptBytes(request *types.LLMRequest) ([]byte, error) {
	if request == nil {
		return nil, errNilRequest
	}
	if request.Body == nil {
		return nil, errNilRequestBody
	}

	switch {
	case request.Body.Completions != nil:
		return []byte(request.Body.Completions.Prompt), nil
	case request.Body.ChatCompletions != nil:
		bytes, err := json.Marshal(request.Body.ChatCompletions.Messages)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal chat completion messages: %w", err)
		}
		return bytes, nil
	default:
		return nil, errUnsupportedInputs
	}
}

// PromptString returns the prompt in string form.
func PromptString(request *types.LLMRequest) (string, error) {
	bytes, err := PromptBytes(request)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// PromptLength returns the number of bytes in the prompt payload that would be used for prefix cache computations.
func PromptLength(request *types.LLMRequest) (int, error) {
	bytes, err := PromptBytes(request)
	if err != nil {
		return 0, err
	}

	return len(bytes), nil
}
