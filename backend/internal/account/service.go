package account

import (
	"context"
	"fmt"
	"strings"

	"github.com/AelcioJozias/vibe-invest/backend/internal/shared/apperrors"
)

type Repository interface {
	List(ctx context.Context, search string) ([]Response, error)
	Create(ctx context.Context, name string) (Response, error)
	UpdateName(ctx context.Context, id int64, name string) (Response, error)
	Deactivate(ctx context.Context, id int64) error
	ExistsActive(ctx context.Context, id int64) (bool, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, search string) ([]Response, error) {
	return s.repo.List(ctx, strings.TrimSpace(search))
}

func (s *Service) Create(ctx context.Context, request CreateRequest) (Response, error) {
	name := strings.TrimSpace(request.Name)
	if name == "" {
		return Response{}, fmt.Errorf("name is required: %w", apperrors.ErrValidation)
	}

	return s.repo.Create(ctx, name)
}

func (s *Service) Update(ctx context.Context, id int64, request UpdateRequest) (Response, error) {
	name := strings.TrimSpace(request.Name)
	if name == "" {
		return Response{}, fmt.Errorf("name is required: %w", apperrors.ErrValidation)
	}

	return s.repo.UpdateName(ctx, id, name)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Deactivate(ctx, id)
}
