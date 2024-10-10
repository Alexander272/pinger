package services

import "github.com/Alexander272/Pinger/internal/repo"

type Services struct {
	Address
	Scheduler
}

type Deps struct {
	repo *repo.Repository
}

func NewServices(deps *Deps) *Services {
	addresses := NewAddressService(deps.repo.Address)
	scheduler := NewSchedulerService()

	return &Services{
		Address:   addresses,
		Scheduler: scheduler,
	}
}
