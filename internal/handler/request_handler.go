package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"service/internal/model"
	"service/internal/repository"
)

type RequestHandler struct {
	repo *repository.RequestRepository
}

func NewRequestHandler(repo *repository.RequestRepository) *RequestHandler {
	return &RequestHandler{repo: repo}
}

func (h *RequestHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input model.CreateRequestInput

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	input.Name = strings.TrimSpace(input.Name)
	input.Phone = strings.TrimSpace(input.Phone)
	input.Address = strings.TrimSpace(input.Address)
	input.Problem = strings.TrimSpace(input.Problem)

	if input.Name == "" || input.Phone == "" || input.Address == "" || input.Problem == "" {
		writeError(w, http.StatusBadRequest, "name, phone, address and problem are required")
		return
	}

	item, err := h.repo.Create(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create request")
		return
	}

	writeJSON(w, http.StatusCreated, item)
}

func (h *RequestHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.repo.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load requests")
		return
	}

	writeJSON(w, http.StatusOK, items)
}

func (h *RequestHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var input model.UpdateStatusInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	input.Status = strings.ToLower(strings.TrimSpace(input.Status))

	switch input.Status {
	case "new", "in_progress", "done", "cancelled":
	default:
		writeError(w, http.StatusBadRequest, "status must be new, in_progress, done or cancelled")
		return
	}

	item, err := h.repo.UpdateStatus(r.Context(), id, input.Status)
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "request not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not update request")
		return
	}

	writeJSON(w, http.StatusOK, item)
}

func (h *RequestHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	err = h.repo.Delete(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "request not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not delete request")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseID(value string) (int64, error) {
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
