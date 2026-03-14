package dashboard

import (
	"errors"
	"net/http"

	"github.com/AelcioJozias/vibe-invest/backend/internal/shared/apperrors"
	"github.com/AelcioJozias/vibe-invest/backend/internal/shared/httpx"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Summary(w http.ResponseWriter, r *http.Request) {
	summary, err := h.service.GetSummary(r.Context(), r.URL.Query().Get("referenceMonth"))
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrValidation):
			httpx.WriteProblem(w, http.StatusBadRequest, "Bad Request", "referenceMonth must be in YYYY-MM format", r.URL.Path)
		default:
			httpx.WriteProblem(w, http.StatusInternalServerError, "Internal Server Error", err.Error(), r.URL.Path)
		}
		return
	}

	httpx.WriteJSON(w, http.StatusOK, summary)
}
