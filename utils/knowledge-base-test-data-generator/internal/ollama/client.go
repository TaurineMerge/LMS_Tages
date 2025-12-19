package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Client is a client for the Ollama API.
type Client struct {
	url   string
	model string
}

// NewClient creates a new Ollama client.
func NewClient(url, model string) *Client {
	return &Client{
		url:   strings.TrimRight(url, "/"),
		model: model,
	}
}

//ollamaRequest defines the structure for a request to the Ollama API.
type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// ollamaResponse defines the structure for a response from the Ollama API.
type ollamaResponse struct {
	Response string `json:"response"`
}

// Generate sends a prompt to Ollama and returns the generated text.
func (c *Client) Generate(prompt string) (string, error) {
	reqBody := ollamaRequest{
		Model:  c.model,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal ollama request: %w", err)
	}

	resp, err := http.Post(fmt.Sprintf("%s/api/generate", c.url), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to send request to ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama request failed with status: %s", resp.Status)
	}

	var ollamaResp ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("failed to decode ollama response: %w", err)
	}

	// Clean up the response
	return strings.TrimSpace(ollamaResp.Response), nil
}
