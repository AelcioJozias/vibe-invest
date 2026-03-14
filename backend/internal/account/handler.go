package account

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/AelcioJozias/vibe-invest/backend/internal/shared/apperrors"
	"github.com/AelcioJozias/vibe-invest/backend/internal/shared/httpx"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.service.List(r.Context(), r.URL.Query().Get("searchString"))
	if err != nil {
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal Server Error", err.Error(), r.URL.Path)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, accounts)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var request CreateRequest
	if err := httpx.DecodeJSON(r, &request); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad Request", err.Error(), r.URL.Path)
		return
	}

	account, err := h.service.Create(r.Context(), request)
	if err != nil {
		h.writeError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, account)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad Request", "invalid account id", r.URL.Path)
		return
	}

	var request UpdateRequest
	if err := httpx.DecodeJSON(r, &request); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad Request", err.Error(), r.URL.Path)
		return
	}

	account, err := h.service.Update(r.Context(), id, request)
	if err != nil {
		h.writeError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, account)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad Request", "invalid account id", r.URL.Path)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		h.writeError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) writeError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, apperrors.ErrValidation):
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad Request", err.Error(), r.URL.Path)
	case errors.Is(err, apperrors.ErrNotFound):
		httpx.WriteProblem(w, http.StatusNotFound, "Not Found", "account not found", r.URL.Path)
	default:
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal Server Error", err.Error(), r.URL.Path)
	}
}

func parseID(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}
