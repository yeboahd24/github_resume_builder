package model

import "time"

type User struct {
	ID              int64
	GitHubID        int64
	Username        string
	Email           string
	Name            string
	AvatarURL       string
	EncryptedToken  string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Resume struct {
	ID          int64
	UserID      int64
	Title       string
	TargetRole  string
	Summary     string
	Projects    []ResumeProject
	Skills      []string
	IsDefault   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ResumeProject struct {
	RepoName    string
	Description string
	URL         string
	Stars       int
	Language    string
	Topics      []string
	Highlights  []string
	Position    int
}

type GitHubProfile struct {
	ID        int64
	Login     string
	Name      string
	Email     string
	AvatarURL string
	Bio       string
	Company   string
	Location  string
}

type Repository struct {
	Name            string
	FullName        string
	Description     string
	URL             string
	Stars           int
	Forks           int
	Language        string
	Topics          []string
	LastCommitDate  time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	IsPrivate       bool
	IsFork          bool
}

type RankedRepository struct {
	Repository
	Score      float64
	Highlights []string
}
