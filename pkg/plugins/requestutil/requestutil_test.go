package requestutil

import (
	"encoding/json"
	"errors"
	"testing"

	"sigs.k8s.io/gateway-api-inference-extension/pkg/epp/scheduling/types"
)

func TestPromptBytes_Completions(t *testing.T) {
	req := &types.LLMRequest{
		Body: &types.LLMRequestBody{
			Completions: &types.CompletionsRequest{Prompt: "hello"},
		},
	}

	bytes, err := PromptBytes(req)
	if err != nil {
		t.Fatalf("PromptBytes returned error: %v", err)
	}

	if string(bytes) != "hello" {
		t.Fatalf("expected prompt bytes to equal completions prompt, got %q", string(bytes))
	}
}

func TestPromptBytes_ChatCompletions(t *testing.T) {
	req := &types.LLMRequest{
		Body: &types.LLMRequestBody{
			ChatCompletions: &types.ChatCompletionsRequest{
				Messages: []types.Message{
					{Role: "user", Content: "hi"},
					{Role: "assistant", Content: "hello"},
				},
			},
		},
	}

	expected, err := json.Marshal(req.Body.ChatCompletions.Messages)
	if err != nil {
		t.Fatalf("failed to marshal expected messages: %v", err)
	}

	bytes, err := PromptBytes(req)
	if err != nil {
		t.Fatalf("PromptBytes returned error: %v", err)
	}

	if string(bytes) != string(expected) {
		t.Fatalf("unexpected prompt bytes for chat completions: got %q want %q", string(bytes), string(expected))
	}
}

func TestPromptBytes_Errors(t *testing.T) {
	t.Run("nil request", func(t *testing.T) {
		if _, err := PromptBytes(nil); !errors.Is(err, errNilRequest) {
			t.Fatalf("expected errNilRequest, got %v", err)
		}
	})

	t.Run("nil body", func(t *testing.T) {
		if _, err := PromptBytes(&types.LLMRequest{}); !errors.Is(err, errNilRequestBody) {
			t.Fatalf("expected errNilRequestBody, got %v", err)
		}
	})

	t.Run("unsupported body", func(t *testing.T) {
		req := &types.LLMRequest{Body: &types.LLMRequestBody{}}
		if _, err := PromptBytes(req); !errors.Is(err, errUnsupportedInputs) {
			t.Fatalf("expected errUnsupportedInputs, got %v", err)
		}
	})
}

func TestPromptLength(t *testing.T) {
	req := &types.LLMRequest{
		Body: &types.LLMRequestBody{
			Completions: &types.CompletionsRequest{Prompt: "abc"},
		},
	}

	length, err := PromptLength(req)
	if err != nil {
		t.Fatalf("PromptLength returned error: %v", err)
	}

	if length != 3 {
		t.Fatalf("expected length 3, got %d", length)
	}
}
