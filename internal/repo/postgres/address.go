package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/Pinger/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AddressRepo struct {
	db *sqlx.DB
}

func NewAddressRepo(db *sqlx.DB) *AddressRepo {
	return &AddressRepo{db: db}
}

type Address interface {
	Get(context.Context) ([]*models.Address, error)
	Create(context.Context, *models.Address) error
	Update(context.Context, *models.Address) error
	Delete(ctx context.Context, ip string) error
}

func (r *AddressRepo) Get(ctx context.Context) ([]*models.Address, error) {
	query := fmt.Sprintf(`SELECT id, ip, name, max_rtt, interval, count, timeout, not_count, period_start, period_end, enabled, created_at 
		FROM %s ORDER BY created_at`,
		AddressTable,
	)
	data := []*models.Address{}

	err := r.db.SelectContext(ctx, &data, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *AddressRepo) Create(ctx context.Context, address *models.Address) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, ip, name, max_rtt, interval, count, timeout, not_count, period_start, period_end, enabled) 
		VALUES (:id, :ip, :name, :max_rtt, :interval, :count, :timeout, :not_count, :period_start, :period_end, :enabled)`,
		AddressTable,
	)
	address.ID = uuid.NewString()

	_, err := r.db.NamedExecContext(ctx, query, address)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *AddressRepo) Update(ctx context.Context, address *models.Address) error {
	query := fmt.Sprintf(`UPDATE %s SET name = :name, max_rtt = :max_rtt, interval = :interval, count = :count, timeout = :timeout, 
		not_count = :not_count, period_start = :period_start, period_end = :period_end, enabled = :enabled WHERE ip = :ip`,
		AddressTable,
	)

	_, err := r.db.NamedExecContext(ctx, query, address)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *AddressRepo) Delete(ctx context.Context, ip string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE ip = $1`, AddressTable)

	_, err := r.db.NamedExecContext(ctx, query, ip)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
