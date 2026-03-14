package dashboard

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AelcioJozias/vibe-invest/backend/internal/shared/apperrors"
)

type fakeDashboardRepository struct {
	response Response
}

func (f fakeDashboardRepository) Summary(ctx context.Context, referenceMonth time.Time) (Response, error) {
	return f.response, nil
}

func TestGetSummaryInvalidReferenceMonth(t *testing.T) {
	service := NewService(fakeDashboardRepository{})

	_, err := service.GetSummary(context.Background(), "2026/03")
	if err == nil {
		t.Fatal("expected validation error")
	}

	if !errors.Is(err, apperrors.ErrValidation) {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestGetSummaryUsesReferenceMonth(t *testing.T) {
	service := NewService(fakeDashboardRepository{})

	summary, err := service.GetSummary(context.Background(), "2026-03")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if summary.ReferenceMonth != "2026-03" {
		t.Fatalf("expected reference month 2026-03, got %s", summary.ReferenceMonth)
	}
}
