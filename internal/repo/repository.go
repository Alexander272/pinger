package repo

import (
	"github.com/Alexander272/Pinger/internal/repo/postgres"
	"github.com/jmoiron/sqlx"
)

type Address interface {
	postgres.Address
}
type Statistic interface {
	postgres.Statistic
}

type Repository struct {
	Address
	Statistic
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Address:   postgres.NewAddressRepo(db),
		Statistic: postgres.NewStatisticRepo(db),
	}
}
