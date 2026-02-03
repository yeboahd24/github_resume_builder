package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/yourusername/resume-builder/internal/client"
	"github.com/yourusername/resume-builder/internal/model"
)

type GitHubService struct {
	client *client.GitHubClient
	cache  *client.CacheClient
}

func NewGitHubService(ghClient *client.GitHubClient, cache *client.CacheClient) *GitHubService {
	return &GitHubService{
		client: ghClient,
		cache:  cache,
	}
}

func (s *GitHubService) FetchUserData(ctx context.Context, token string) (*model.GitHubProfile, []model.Repository, error) {
	cacheKey := fmt.Sprintf("github:repos:%s", token[:10])

	// Try cache first
	if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
		var data struct {
			Profile *model.GitHubProfile
			Repos   []model.Repository
		}
		if json.Unmarshal([]byte(cached), &data) == nil {
			return data.Profile, data.Repos, nil
		}
	}

	profile, err := s.client.GetProfile(ctx, token)
	if err != nil {
		return nil, nil, err
	}

	repos, err := s.client.GetRepositories(ctx, token)
	if err != nil {
		return nil, nil, err
	}

	// Cache for 1 hour
	if data, err := json.Marshal(map[string]interface{}{"Profile": profile, "Repos": repos}); err == nil {
		s.cache.Set(ctx, cacheKey, string(data), time.Hour)
	}

	return profile, repos, nil
}

func (s *GitHubService) ExtractSkills(repos []model.Repository) []string {
	skillSet := make(map[string]bool)

	for _, repo := range repos {
		if repo.Language != "" {
			skillSet[repo.Language] = true
		}

		for _, topic := range repo.Topics {
			skillSet[topic] = true
		}
	}

	skills := make([]string, 0, len(skillSet))
	for skill := range skillSet {
		skills = append(skills, skill)
	}

	return skills
}
