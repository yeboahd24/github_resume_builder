package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/resume-builder/internal/model"
)

type GitHubClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewGitHubClient() *GitHubClient {
	return &GitHubClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    "https://api.github.com",
	}
}

func (c *GitHubClient) GetProfile(ctx context.Context, token string) (*model.GitHubProfile, error) {
	var profile struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
		Bio       string `json:"bio"`
		Company   string `json:"company"`
		Location  string `json:"location"`
	}

	if err := c.doRequest(ctx, "GET", "/user", token, &profile); err != nil {
		return nil, err
	}

	return &model.GitHubProfile{
		ID:        profile.ID,
		Login:     profile.Login,
		Name:      profile.Name,
		Email:     profile.Email,
		AvatarURL: profile.AvatarURL,
		Bio:       profile.Bio,
		Company:   profile.Company,
		Location:  profile.Location,
	}, nil
}

func (c *GitHubClient) GetRepositories(ctx context.Context, token string) ([]model.Repository, error) {
	var repos []struct {
		Name        string    `json:"name"`
		FullName    string    `json:"full_name"`
		Description string    `json:"description"`
		HTMLURL     string    `json:"html_url"`
		Stars       int       `json:"stargazers_count"`
		Forks       int       `json:"forks_count"`
		Language    string    `json:"language"`
		Topics      []string  `json:"topics"`
		Private     bool      `json:"private"`
		Fork        bool      `json:"fork"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		PushedAt    time.Time `json:"pushed_at"`
	}

	if err := c.doRequest(ctx, "GET", "/user/repos?per_page=100&sort=updated", token, &repos); err != nil {
		return nil, err
	}

	result := make([]model.Repository, len(repos))
	for i, r := range repos {
		result[i] = model.Repository{
			Name:           r.Name,
			FullName:       r.FullName,
			Description:    r.Description,
			URL:            r.HTMLURL,
			Stars:          r.Stars,
			Forks:          r.Forks,
			Language:       r.Language,
			Topics:         r.Topics,
			LastCommitDate: r.PushedAt,
			CreatedAt:      r.CreatedAt,
			UpdatedAt:      r.UpdatedAt,
			IsPrivate:      r.Private,
			IsFork:         r.Fork,
		}
	}

	return result, nil
}

func (c *GitHubClient) doRequest(ctx context.Context, method, path, token string, result interface{}) error {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("github api error: status %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(result)
}
