package repo

import (
	"github.com/Alexander272/Pinger/internal/repo/postgres"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type Address interface {
	postgres.Address
}

type Repository struct {
	Address
}

func NewRepository(db *sqlx.DB, redis *redis.Client) *Repository {
	return &Repository{
		Address: postgres.NewAddressRepo(db),
	}
}
