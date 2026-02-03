package service

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"

	"github.com/yourusername/resume-builder/internal/client"
	"github.com/yourusername/resume-builder/internal/crypto"
	"github.com/yourusername/resume-builder/internal/model"
	"github.com/yourusername/resume-builder/internal/repository"
)

type AuthService struct {
	userRepo      *repository.UserRepository
	githubClient  *client.GitHubClient
	encryptor     *crypto.Encryptor
	oauthConfig   *oauth2.Config
}

func NewAuthService(
	userRepo *repository.UserRepository,
	githubClient *client.GitHubClient,
	encryptor *crypto.Encryptor,
	clientID, clientSecret, redirectURL string,
) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		githubClient: githubClient,
		encryptor:    encryptor,
		oauthConfig: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"read:user", "user:email", "repo"},
			Endpoint:     github.Endpoint,
		},
	}
}

func (s *AuthService) GetAuthURL(state string) string {
	return s.oauthConfig.AuthCodeURL(state)
}

func (s *AuthService) HandleCallback(ctx context.Context, code string) (*model.User, error) {
	token, err := s.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	profile, err := s.githubClient.GetProfile(ctx, token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	encryptedToken, err := s.encryptor.Encrypt(token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt token: %w", err)
	}

	user, err := s.userRepo.GetByGitHubID(ctx, profile.ID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user = &model.User{
			GitHubID:       profile.ID,
			Username:       profile.Login,
			Email:          profile.Email,
			Name:           profile.Name,
			AvatarURL:      profile.AvatarURL,
			EncryptedToken: encryptedToken,
		}
		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, err
		}
	} else {
		if err := s.userRepo.UpdateToken(ctx, user.ID, encryptedToken); err != nil {
			return nil, err
		}
		user.EncryptedToken = encryptedToken
	}

	return user, nil
}

func (s *AuthService) GetDecryptedToken(encryptedToken string) (string, error) {
	return s.encryptor.Decrypt(encryptedToken)
}

func (s *AuthService) GetUserToken(ctx context.Context, userID int64) (string, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", fmt.Errorf("user not found")
	}
	return s.encryptor.Decrypt(user.EncryptedToken)
}
