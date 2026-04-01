package bartender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	apiURL          = "https://api.openai.com/v1/chat/completions"
	model           = "gpt-4.1-nano"
	maxTokens       = 200
	cooldownPerUser = 10 * time.Second
	contextMessages = 20
)

// ChatMsg is a minimal chat message for building context.
type ChatMsg struct {
	Nickname string
	Text     string
}

// Bartender handles the tavern bartender AI persona.
type Bartender struct {
	apiKey    string
	soul      string
	mu        sync.Mutex
	cooldowns map[string]time.Time // fingerprint → last response time
}

// New creates a bartender. Returns nil if apiKey is empty.
func New(apiKey, soul string) *Bartender {
	if apiKey == "" {
		return nil
	}
	return &Bartender{
		apiKey:    apiKey,
		soul:      soul,
		cooldowns: make(map[string]time.Time),
	}
}

// ShouldRespond checks if a message triggers the bartender.
func ShouldRespond(text, room string) bool {
	if room != "lounge" {
		return false
	}
	lower := strings.ToLower(text)
	return strings.Contains(lower, "@bartender")
}

// CanRespond checks the per-user cooldown. Returns true if allowed.
func (b *Bartender) CanRespond(fingerprint string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	last, ok := b.cooldowns[fingerprint]
	if ok && time.Since(last) < cooldownPerUser {
		return false
	}
	b.cooldowns[fingerprint] = time.Now()
	return true
}

type apiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type apiRequest struct {
	Model     string       `json:"model"`
	Messages  []apiMessage `json:"messages"`
	MaxTokens int          `json:"max_tokens"`
}

type apiChoice struct {
	Message apiMessage `json:"message"`
}

type apiResponse struct {
	Choices []apiChoice `json:"choices"`
	Error   *apiError   `json:"error,omitempty"`
}

type apiError struct {
	Message string `json:"message"`
}

// Respond generates a bartender response given recent chat context.
func (b *Bartender) Respond(recentMessages []ChatMsg, triggerNick, triggerText string) (string, error) {
	// Build conversation context
	var contextParts []string
	for _, m := range recentMessages {
		contextParts = append(contextParts, fmt.Sprintf("%s: %s", m.Nickname, m.Text))
	}
	chatContext := strings.Join(contextParts, "\n")

	messages := []apiMessage{
		{Role: "system", Content: b.soul},
		{Role: "user", Content: fmt.Sprintf("Recent tavern chat:\n%s\n\n%s says to you: %s", chatContext, triggerNick, triggerText)},
	}

	reqBody := apiRequest{
		Model:     model,
		Messages:  messages,
		MaxTokens: maxTokens,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+b.apiKey)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("api call: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}

	var apiResp apiResponse
	if err := json.Unmarshal(respBytes, &apiResp); err != nil {
		return "", fmt.Errorf("unmarshal: %w", err)
	}

	if apiResp.Error != nil {
		return "", fmt.Errorf("api error: %s", apiResp.Error.Message)
	}

	if len(apiResp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}

	reply := strings.TrimSpace(apiResp.Choices[0].Message.Content)
	if reply == "" {
		return "", fmt.Errorf("empty response")
	}

	log.Printf("bartender: replied to %s (%d chars)", triggerNick, len(reply))
	return reply, nil
}
