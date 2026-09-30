package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"service/internal/auth"
	"service/internal/model"
	"service/internal/repository"
)

type AuthHandler struct {
	users *repository.UserRepository
}

func NewAuthHandler(users *repository.UserRepository) *AuthHandler {
	return &AuthHandler{users: users}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input model.RegisterInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	input.Password = strings.TrimSpace(input.Password)

	if input.Email == "" || input.Password == "" {
		writeError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	if len(input.Password) < 6 {
		writeError(w, http.StatusBadRequest, "password must be at least 6 characters")
		return
	}

	passwordHash, err := auth.HashPassword(input.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not hash password")
		return
	}

	user, err := h.users.Create(r.Context(), input.Email, passwordHash)
	if err != nil {
		writeError(w, http.StatusBadRequest, "user already exists or could not be created")
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Role)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not generate token")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"user": map[string]any{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		},
		"token": token,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input model.LoginInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	input.Email = strings.ToLower(strings.TrimSpace(input.Email))

	user, err := h.users.GetByEmail(r.Context(), input.Email)
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not login")
		return
	}

	if !auth.CheckPassword(input.Password, user.PasswordHash) {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Role)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not generate token")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"user": map[string]any{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		},
		"token": token,
	})
}
