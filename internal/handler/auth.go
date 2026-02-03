package handler

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/yourusername/resume-builder/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
	jwtService  *service.JWTService
	frontendURL string
}

func NewAuthHandler(authService *service.AuthService, jwtService *service.JWTService, frontendURL string) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		jwtService:  jwtService,
		frontendURL: frontendURL,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	state := generateState()
	
	// In production, store state in session/cookie with expiry
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		MaxAge:   300,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	url := h.authService.GetAuthURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil {
		respondError(w, http.StatusBadRequest, "missing state cookie")
		return
	}

	state := r.URL.Query().Get("state")
	if state != stateCookie.Value {
		respondError(w, http.StatusBadRequest, "invalid state")
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		respondError(w, http.StatusBadRequest, "missing code")
		return
	}

	user, err := h.authService.HandleCallback(r.Context(), code)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "authentication failed")
		return
	}

	token, err := h.jwtService.GenerateToken(user.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	// Redirect to frontend with token
	redirectURL := h.frontendURL + "/callback?token=" + token
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func generateState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
