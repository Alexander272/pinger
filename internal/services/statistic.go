package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Alexander272/Pinger/internal/models"
	"github.com/Alexander272/Pinger/internal/repo"
)

type StatisticService struct {
	repo repo.Statistic
}

func NewStatisticService(repo repo.Statistic) *StatisticService {
	return &StatisticService{
		repo: repo,
	}
}

type Statistic interface {
	Get(ctx context.Context, req *models.GetStatisticDTO) ([]*models.Statistic, error)
	GetByIP(ctx context.Context, req *models.GetStatisticByIPDTO) ([]*models.Statistic, error)
	GetUnavailable(ctx context.Context, req *models.GetUnavailableDTO) ([]*models.Statistic, error)
	Create(ctx context.Context, dto *models.StatisticDTO) error
	Update(ctx context.Context, dto *models.StatisticDTO) error
}

func (s *StatisticService) Get(ctx context.Context, req *models.GetStatisticDTO) ([]*models.Statistic, error) {
	data, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistic. error: %w", err)
	}
	return data, nil
}

func (s *StatisticService) GetByIP(ctx context.Context, req *models.GetStatisticByIPDTO) ([]*models.Statistic, error) {
	data, err := s.repo.GetByIP(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistic by ip. error: %w", err)
	}
	return data, nil
}

func (s *StatisticService) GetUnavailable(ctx context.Context, req *models.GetUnavailableDTO) ([]*models.Statistic, error) {
	data, err := s.repo.GetUnavailable(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get unavailable ip. error: %w", err)
	}
	return data, nil
}

func (s *StatisticService) Create(ctx context.Context, dto *models.StatisticDTO) error {
	last, err := s.repo.GetLast(ctx, &models.GetStatisticByIPDTO{IP: dto.IP})
	if err != nil && !errors.Is(err, models.ErrNoRows) {
		return fmt.Errorf("failed to get last statistic. error: %w", err)
	}
	if last != nil {
		return nil
	}

	if err := s.repo.Create(ctx, dto); err != nil {
		return fmt.Errorf("failed to create statistic. error: %w", err)
	}
	return nil
}

func (s *StatisticService) Update(ctx context.Context, dto *models.StatisticDTO) error {
	if err := s.repo.Update(ctx, dto); err != nil {
		return fmt.Errorf("failed to update statistic. error: %w", err)
	}
	return nil
}
