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
		"Write a professional 3-5 sentence resume summary for a %s position. The candidate has %d GitHub repositories showcasing expertise in: %v. Highlight technical depth, impact, and career goals. Be specific and achievement-oriented.",
		targetRole, repoCount, skills[:min(5, len(skills))],
	)

	reqBody := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{"role": "system", "content": "You are an expert resume writer. Write compelling, achievement-focused summaries that highlight technical expertise and career impact."},
			{"role": "user", "content": prompt},
		},
		"max_tokens":  200,
		"temperature": 0.8,
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

func (c *LLMClient) EnhanceProjectDescription(ctx context.Context, repoName, description, language string, topics []string) (string, []string, error) {
	if !c.enabled || c.apiKey == "" {
		return description, []string{}, fmt.Errorf("llm not enabled")
	}

	prompt := fmt.Sprintf(
		"Project: %s\nLanguage: %s\nTopics: %v\nOriginal description: %s\n\nWrite a professional 1-sentence project description and 2-3 bullet points highlighting technical achievements, impact, or key features. Format as JSON: {\"description\": \"...\", \"highlights\": [\"...\", \"...\"]}",
		repoName, language, topics, description,
	)

	reqBody := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{"role": "system", "content": "You are a technical resume writer. Create compelling project descriptions that highlight technical skills and impact. Always respond with valid JSON."},
			{"role": "user", "content": prompt},
		},
		"max_tokens":  250,
		"temperature": 0.7,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return description, []string{}, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewBuffer(body))
	if err != nil {
		return description, []string{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return description, []string{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return description, []string{}, fmt.Errorf("openai api error: status %d", resp.StatusCode)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return description, []string{}, err
	}

	if len(result.Choices) == 0 {
		return description, []string{}, fmt.Errorf("no response from llm")
	}

	// Parse the JSON response
	var enhanced struct {
		Description string   `json:"description"`
		Highlights  []string `json:"highlights"`
	}

	if err := json.Unmarshal([]byte(result.Choices[0].Message.Content), &enhanced); err != nil {
		// If parsing fails, return original
		return description, []string{}, err
	}

	return enhanced.Description, enhanced.Highlights, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
