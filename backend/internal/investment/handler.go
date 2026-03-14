package investment

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

func (h *Handler) ListByAccount(w http.ResponseWriter, r *http.Request) {
	accountID, err := parseID(r.PathValue("accountId"))
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad Request", "invalid account id", r.URL.Path)
		return
	}

	items, err := h.service.ListByAccount(r.Context(), accountID)
	if err != nil {
		h.writeError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, items)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	accountID, err := parseID(r.PathValue("accountId"))
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad Request", "invalid account id", r.URL.Path)
		return
	}

	var request CreateRequest
	if err := httpx.DecodeJSON(r, &request); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad Request", err.Error(), r.URL.Path)
		return
	}

	investment, err := h.service.Create(r.Context(), accountID, request)
	if err != nil {
		h.writeError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, investment)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	investmentID, err := parseID(r.PathValue("investmentId"))
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad Request", "invalid investment id", r.URL.Path)
		return
	}

	investment, err := h.service.GetByID(r.Context(), investmentID)
	if err != nil {
		h.writeError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, investment)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	investmentID, err := parseID(r.PathValue("investmentId"))
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad Request", "invalid investment id", r.URL.Path)
		return
	}

	var request UpdateRequest
	if err := httpx.DecodeJSON(r, &request); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad Request", err.Error(), r.URL.Path)
		return
	}

	investment, err := h.service.Update(r.Context(), investmentID, request)
	if err != nil {
		h.writeError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, investment)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	investmentID, err := parseID(r.PathValue("investmentId"))
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad Request", "invalid investment id", r.URL.Path)
		return
	}

	if err := h.service.Delete(r.Context(), investmentID); err != nil {
		h.writeError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) IncrementFees(w http.ResponseWriter, r *http.Request) {
	investmentID, err := parseID(r.PathValue("investmentId"))
	if err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad Request", "invalid investment id", r.URL.Path)
		return
	}

	var request IncrementFeesRequest
	if err := httpx.DecodeJSON(r, &request); err != nil {
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad Request", err.Error(), r.URL.Path)
		return
	}

	investment, err := h.service.IncrementFees(r.Context(), investmentID, request)
	if err != nil {
		h.writeError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, investment)
}

func (h *Handler) writeError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, apperrors.ErrValidation):
		httpx.WriteProblem(w, http.StatusBadRequest, "Bad Request", err.Error(), r.URL.Path)
	case errors.Is(err, apperrors.ErrNotFound):
		httpx.WriteProblem(w, http.StatusNotFound, "Not Found", "investment not found", r.URL.Path)
	default:
		httpx.WriteProblem(w, http.StatusInternalServerError, "Internal Server Error", err.Error(), r.URL.Path)
	}
}

func parseID(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}
