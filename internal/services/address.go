package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/Pinger/internal/models"
	"github.com/Alexander272/Pinger/internal/repo"
)

type AddressService struct {
	repo repo.Address
}

func NewAddressService(repo repo.Address) *AddressService {
	return &AddressService{
		repo: repo,
	}
}

type Address interface {
	Get(ctx context.Context) ([]*models.Address, error)
	Create(ctx context.Context, address *models.AddressDTO) error
	Update(ctx context.Context, address *models.AddressDTO) error
	Delete(ctx context.Context, ip string) error
}

func (s *AddressService) Get(ctx context.Context) ([]*models.Address, error) {
	data, err := s.repo.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get addresses. error: %w", err)
	}
	return data, nil
}

func (s *AddressService) Create(ctx context.Context, address *models.AddressDTO) error {
	if err := s.repo.Create(ctx, address); err != nil {
		return fmt.Errorf("failed to create addresses. error: %w", err)
	}
	return nil
}

func (s *AddressService) Update(ctx context.Context, address *models.AddressDTO) error {
	if err := s.repo.Update(ctx, address); err != nil {
		return fmt.Errorf("failed to update addresses. error: %w", err)
	}
	return nil
}

func (s *AddressService) Delete(ctx context.Context, ip string) error {
	if err := s.repo.Delete(ctx, ip); err != nil {
		return fmt.Errorf("failed to delete addresses. error: %w", err)
	}
	return nil
}
