package account

import (
	"context"
	"errors"
	"testing"

	"github.com/AelcioJozias/vibe-invest/backend/internal/shared/apperrors"
)

type fakeAccountRepository struct{}

func (f fakeAccountRepository) List(ctx context.Context, search string) ([]Response, error) {
	return nil, nil
}

func (f fakeAccountRepository) Create(ctx context.Context, name string) (Response, error) {
	return Response{}, nil
}

func (f fakeAccountRepository) UpdateName(ctx context.Context, id int64, name string) (Response, error) {
	return Response{}, nil
}

func (f fakeAccountRepository) Deactivate(ctx context.Context, id int64) error {
	return nil
}

func (f fakeAccountRepository) ExistsActive(ctx context.Context, id int64) (bool, error) {
	return true, nil
}

func TestCreateValidation(t *testing.T) {
	service := NewService(fakeAccountRepository{})

	_, err := service.Create(context.Background(), CreateRequest{Name: "   "})
	if err == nil {
		t.Fatal("expected validation error")
	}

	if !errors.Is(err, apperrors.ErrValidation) {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}
