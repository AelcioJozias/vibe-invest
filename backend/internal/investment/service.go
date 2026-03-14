package investment

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/AelcioJozias/vibe-invest/backend/internal/shared/apperrors"
)

type Repository interface {
	ListByAccount(ctx context.Context, accountID int64) ([]record, error)
	Create(ctx context.Context, accountID int64, input CreateRequest, referenceMonth time.Time) (record, error)
	GetByID(ctx context.Context, investmentID int64) (record, error)
	Update(ctx context.Context, investmentID int64, input UpdateRequest, referenceMonth time.Time) (record, error)
	Deactivate(ctx context.Context, investmentID int64) error
	IncrementFees(ctx context.Context, investmentID int64, amount int64, referenceMonth time.Time) (record, error)
}

type AccountChecker interface {
	ExistsActive(ctx context.Context, id int64) (bool, error)
}

type Service struct {
	repo         Repository
	accountCheck AccountChecker
	now          func() time.Time
}

func NewService(repo Repository, accountCheck AccountChecker) *Service {
	return &Service{
		repo:         repo,
		accountCheck: accountCheck,
		now:          time.Now,
	}
}

func (s *Service) ListByAccount(ctx context.Context, accountID int64) ([]Response, error) {
	exists, err := s.accountCheck.ExistsActive(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, apperrors.ErrNotFound
	}

	rows, err := s.repo.ListByAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return toResponses(rows), nil
}

func (s *Service) Create(ctx context.Context, accountID int64, input CreateRequest) (Response, error) {
	if err := validateCreate(input); err != nil {
		return Response{}, err
	}

	investment, err := s.repo.Create(ctx, accountID, normalizeCreateInput(input), s.now())
	if err != nil {
		return Response{}, err
	}

	return toResponse(investment), nil
}

func (s *Service) GetByID(ctx context.Context, investmentID int64) (Response, error) {
	investment, err := s.repo.GetByID(ctx, investmentID)
	if err != nil {
		return Response{}, err
	}

	return toResponse(investment), nil
}

func (s *Service) Update(ctx context.Context, investmentID int64, input UpdateRequest) (Response, error) {
	if err := validateUpdate(input); err != nil {
		return Response{}, err
	}

	investment, err := s.repo.Update(ctx, investmentID, normalizeUpdateInput(input), s.now())
	if err != nil {
		return Response{}, err
	}

	return toResponse(investment), nil
}

func (s *Service) Delete(ctx context.Context, investmentID int64) error {
	return s.repo.Deactivate(ctx, investmentID)
}

func (s *Service) IncrementFees(ctx context.Context, investmentID int64, input IncrementFeesRequest) (Response, error) {
	if input.Amount <= 0 {
		return Response{}, fmt.Errorf("amount must be greater than zero: %w", apperrors.ErrValidation)
	}

	investment, err := s.repo.IncrementFees(ctx, investmentID, input.Amount, s.now())
	if err != nil {
		return Response{}, err
	}

	return toResponse(investment), nil
}

func validateCreate(input CreateRequest) error {
	if input.Amount < 0 {
		return fmt.Errorf("amount must be greater than or equal to zero: %w", apperrors.ErrValidation)
	}
	if strings.TrimSpace(input.YieldRate) == "" {
		return fmt.Errorf("yieldRate is required: %w", apperrors.ErrValidation)
	}
	return nil
}

func validateUpdate(input UpdateRequest) error {
	if input.Amount < 0 {
		return fmt.Errorf("amount must be greater than or equal to zero: %w", apperrors.ErrValidation)
	}
	if strings.TrimSpace(input.YieldRate) == "" {
		return fmt.Errorf("yieldRate is required: %w", apperrors.ErrValidation)
	}
	return nil
}

func normalizeCreateInput(input CreateRequest) CreateRequest {
	input.YieldRate = strings.TrimSpace(input.YieldRate)
	input.Observation = strings.TrimSpace(input.Observation)
	return input
}

func normalizeUpdateInput(input UpdateRequest) UpdateRequest {
	input.YieldRate = strings.TrimSpace(input.YieldRate)
	input.Observation = strings.TrimSpace(input.Observation)
	return input
}

func toResponses(records []record) []Response {
	responses := make([]Response, 0, len(records))
	for _, item := range records {
		responses = append(responses, toResponse(item))
	}
	return responses
}

func toResponse(item record) Response {
	return Response{
		ID:       item.ID,
		Amount:   item.Amount,
		YieldRate: item.YieldRate,
		IsActive: item.IsActive,
	}
}
