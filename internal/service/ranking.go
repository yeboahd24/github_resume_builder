package service

import (
	"math"
	"sort"
	"strings"
	"time"

	"github.com/yourusername/resume-builder/internal/model"
)

type RankingService struct{}

func NewRankingService() *RankingService {
	return &RankingService{}
}

func (s *RankingService) RankRepositories(repos []model.Repository) []model.RankedRepository {
	ranked := make([]model.RankedRepository, 0, len(repos))

	for _, repo := range repos {
		if repo.IsFork {
			continue
		}

		score := s.calculateScore(repo)
		highlights := s.generateHighlights(repo)

		ranked = append(ranked, model.RankedRepository{
			Repository: repo,
			Score:      score,
			Highlights: highlights,
		})
	}

	sort.Slice(ranked, func(i, j int) bool {
		return ranked[i].Score > ranked[j].Score
	})

	return ranked
}

func (s *RankingService) calculateScore(repo model.Repository) float64 {
	var score float64

	// Stars weight: 30%
	score += math.Log1p(float64(repo.Stars)) * 3.0

	// Recency weight: 25%
	daysSinceUpdate := time.Since(repo.LastCommitDate).Hours() / 24
	recencyScore := math.Max(0, 10-daysSinceUpdate/30)
	score += recencyScore * 2.5

	// Language weight: 20%
	if repo.Language != "" {
		score += 2.0
	}

	// Topics weight: 15%
	score += float64(len(repo.Topics)) * 0.5

	// Description weight: 10%
	if repo.Description != "" {
		score += 1.0
	}

	return score
}

func (s *RankingService) generateHighlights(repo model.Repository) []string {
	var highlights []string

	if repo.Stars > 10 {
		highlights = append(highlights, "Popular project with community engagement")
	}

	if time.Since(repo.LastCommitDate).Hours() < 30*24 {
		highlights = append(highlights, "Actively maintained")
	}

	if len(repo.Topics) > 0 {
		highlights = append(highlights, "Tagged: "+strings.Join(repo.Topics[:min(3, len(repo.Topics))], ", "))
	}

	return highlights
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
