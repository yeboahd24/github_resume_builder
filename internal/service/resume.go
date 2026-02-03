package service

import (
	"context"
	"fmt"

	"github.com/yourusername/resume-builder/internal/client"
	"github.com/yourusername/resume-builder/internal/model"
	"github.com/yourusername/resume-builder/internal/repository"
)

type ResumeService struct {
	resumeRepo     *repository.ResumeRepository
	userRepo       *repository.UserRepository
	githubService  *GitHubService
	rankingService *RankingService
	llmClient      *client.LLMClient
}

func NewResumeService(
	resumeRepo *repository.ResumeRepository,
	userRepo *repository.UserRepository,
	githubService *GitHubService,
	rankingService *RankingService,
	llmClient *client.LLMClient,
) *ResumeService {
	return &ResumeService{
		resumeRepo:     resumeRepo,
		userRepo:       userRepo,
		githubService:  githubService,
		rankingService: rankingService,
		llmClient:      llmClient,
	}
}

func (s *ResumeService) GenerateResume(ctx context.Context, userID int64, targetRole string, token string) (*model.Resume, error) {
	_, repos, err := s.githubService.FetchUserData(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch github data: %w", err)
	}

	rankedRepos := s.rankingService.RankRepositories(repos)
	skills := s.githubService.ExtractSkills(repos)

	topProjects := s.selectTopProjects(rankedRepos, 5)

	// Try LLM summary first, fallback to rule-based
	summary, err := s.llmClient.GenerateSummary(ctx, targetRole, len(repos), skills)
	if err != nil {
		summary = s.generateSummary(targetRole, len(repos), skills)
	}

	resume := &model.Resume{
		UserID:     userID,
		Title:      "GitHub Resume",
		TargetRole: targetRole,
		Summary:    summary,
		Projects:   topProjects,
		Skills:     skills,
		IsDefault:  true,
	}

	if err := s.resumeRepo.Create(ctx, resume); err != nil {
		return nil, err
	}

	return resume, nil
}

func (s *ResumeService) GetResume(ctx context.Context, resumeID, userID int64) (*model.Resume, error) {
	resume, err := s.resumeRepo.GetByID(ctx, resumeID)
	if err != nil {
		return nil, err
	}

	if resume == nil {
		return nil, fmt.Errorf("resume not found")
	}

	if resume.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}

	return resume, nil
}

func (s *ResumeService) ListResumes(ctx context.Context, userID int64) ([]model.Resume, error) {
	return s.resumeRepo.ListByUserID(ctx, userID)
}

func (s *ResumeService) UpdateResume(ctx context.Context, resume *model.Resume, userID int64) error {
	existing, err := s.resumeRepo.GetByID(ctx, resume.ID)
	if err != nil {
		return err
	}

	if existing == nil {
		return fmt.Errorf("resume not found")
	}

	if existing.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	return s.resumeRepo.Update(ctx, resume)
}

func (s *ResumeService) DeleteResume(ctx context.Context, resumeID, userID int64) error {
	resume, err := s.resumeRepo.GetByID(ctx, resumeID)
	if err != nil {
		return err
	}

	if resume == nil {
		return fmt.Errorf("resume not found")
	}

	if resume.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	return s.resumeRepo.Delete(ctx, resumeID)
}

func (s *ResumeService) selectTopProjects(rankedRepos []model.RankedRepository, count int) []model.ResumeProject {
	if len(rankedRepos) < count {
		count = len(rankedRepos)
	}

	projects := make([]model.ResumeProject, count)
	for i := 0; i < count; i++ {
		repo := rankedRepos[i]
		
		// Try to enhance with LLM
		description := repo.Description
		highlights := repo.Highlights
		
		if enhancedDesc, enhancedHighlights, err := s.llmClient.EnhanceProjectDescription(
			context.Background(),
			repo.Name,
			repo.Description,
			repo.Language,
			repo.Topics,
		); err == nil && enhancedDesc != "" {
			description = enhancedDesc
			if len(enhancedHighlights) > 0 {
				highlights = enhancedHighlights
			}
		}
		
		projects[i] = model.ResumeProject{
			RepoName:    repo.Name,
			Description: description,
			URL:         repo.URL,
			Stars:       repo.Stars,
			Language:    repo.Language,
			Topics:      repo.Topics,
			Highlights:  highlights,
			Position:    i,
		}
	}

	return projects
}

func (s *ResumeService) generateSummary(targetRole string, repoCount int, skills []string) string {
	return fmt.Sprintf(
		"Software engineer with %d public repositories. Experienced in %s. Seeking %s role.",
		repoCount,
		s.formatSkills(skills, 5),
		targetRole,
	)
}

func (s *ResumeService) formatSkills(skills []string, max int) string {
	if len(skills) == 0 {
		return "various technologies"
	}

	if len(skills) > max {
		skills = skills[:max]
	}

	result := ""
	for i, skill := range skills {
		if i > 0 {
			result += ", "
		}
		result += skill
	}

	return result
}
