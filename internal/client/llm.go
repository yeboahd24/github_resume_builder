package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type LLMClient struct {
	apiKey  string
	baseURL string
	client  *http.Client
	enabled bool
}

func NewLLMClient(apiKey string, enabled bool) *LLMClient {
	return &LLMClient{
		apiKey:  apiKey,
		baseURL: "https://api.openai.com/v1/chat/completions",
		client:  &http.Client{Timeout: 30 * time.Second},
		enabled: enabled,
	}
}

func (c *LLMClient) GenerateSummary(ctx context.Context, targetRole string, repoCount int, skills []string) (string, error) {
	if !c.enabled || c.apiKey == "" {
		return "", fmt.Errorf("llm not enabled")
	}

	prompt := fmt.Sprintf(
		"Write a professional 2-sentence resume summary for a %s role. The candidate has %d GitHub repositories and skills in: %v. Be concise and impactful.",
		targetRole, repoCount, skills[:min(5, len(skills))],
	)

	reqBody := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{"role": "system", "content": "You are a professional resume writer. Write concise, impactful summaries."},
			{"role": "user", "content": prompt},
		},
		"max_tokens":  150,
		"temperature": 0.7,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("openai api error: status %d", resp.StatusCode)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from llm")
	}

	return result.Choices[0].Message.Content, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
