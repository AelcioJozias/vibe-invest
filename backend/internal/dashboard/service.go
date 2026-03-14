package dashboard

import (
	"context"
	"fmt"
	"time"

	"github.com/AelcioJozias/vibe-invest/backend/internal/shared/apperrors"
	"github.com/AelcioJozias/vibe-invest/backend/internal/shared/timeutil"
)

type Repository interface {
	Summary(ctx context.Context, referenceMonth time.Time) (Response, error)
}

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		now:  time.Now,
	}
}

func (s *Service) GetSummary(ctx context.Context, rawReferenceMonth string) (Response, error) {
	referenceMonth, err := s.resolveReferenceMonth(rawReferenceMonth)
	if err != nil {
		return Response{}, err
	}

	summary, err := s.repo.Summary(ctx, referenceMonth)
	if err != nil {
		return Response{}, err
	}

	summary.ReferenceMonth = referenceMonth.Format("2006-01")
	return summary, nil
}

func (s *Service) resolveReferenceMonth(rawReferenceMonth string) (time.Time, error) {
	if rawReferenceMonth == "" {
		return timeutil.CurrentReferenceMonth(s.now()), nil
	}

	referenceMonth, err := timeutil.ParseReferenceMonth(rawReferenceMonth)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid referenceMonth: %w", apperrors.ErrValidation)
	}

	return referenceMonth, nil
}
