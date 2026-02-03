package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/yourusername/resume-builder/internal/model"
	"github.com/yourusername/resume-builder/internal/service"
)

type ResumeHandler struct {
	resumeService *service.ResumeService
	authService   *service.AuthService
}

func NewResumeHandler(resumeService *service.ResumeService, authService *service.AuthService) *ResumeHandler {
	return &ResumeHandler{
		resumeService: resumeService,
		authService:   authService,
	}
}

func (h *ResumeHandler) Generate(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		TargetRole string `json:"target_role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	token, err := h.authService.GetUserToken(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get user token")
		return
	}

	resume, err := h.resumeService.GenerateResume(r.Context(), userID, req.TargetRole, token)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate resume")
		return
	}

	respondJSON(w, http.StatusCreated, resume)
}

func (h *ResumeHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	resumeID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid resume id")
		return
	}

	resume, err := h.resumeService.GetResume(r.Context(), resumeID, userID)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resume)
}

func (h *ResumeHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	resumes, err := h.resumeService.ListResumes(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list resumes")
		return
	}

	respondJSON(w, http.StatusOK, resumes)
}

func (h *ResumeHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	resumeID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid resume id")
		return
	}

	var resume model.Resume
	if err := json.NewDecoder(r.Body).Decode(&resume); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resume.ID = resumeID
	if err := h.resumeService.UpdateResume(r.Context(), &resume, userID); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resume)
}

func (h *ResumeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	resumeID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid resume id")
		return
	}

	if err := h.resumeService.DeleteResume(r.Context(), resumeID, userID); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
