package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Alexander272/Pinger/internal/models"
	"github.com/Alexander272/Pinger/internal/repo/postgres/pq_models"
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
	Create(context.Context, *models.AddressDTO) error
	Update(context.Context, *models.AddressDTO) error
	Delete(ctx context.Context, ip string) error
}

func (r *AddressRepo) Get(ctx context.Context) ([]*models.Address, error) {
	query := fmt.Sprintf(`SELECT id, ip, name, max_rtt, interval, count, timeout, not_count, period_start, period_end, enabled, created_at 
		FROM %s WHERE enabled=true ORDER BY created_at`,
		AddressTable,
	)
	tmp := []*pq_models.Address{}
	data := []*models.Address{}

	//TODO если я буду хранить данные не в ns, тогда придется создать структуру в которую будут записываться данные из базы, а затем нужно будет преобразовывать их в time.Duration
	/* Если хранить это все в int
	*	max_rtt, interval, timeout число в миллисекундах
	*	period_start, period_end в минутах
	 */

	err := r.db.SelectContext(ctx, &tmp, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, v := range tmp {
		data = append(data, &models.Address{
			ID:                v.ID,
			IP:                v.IP,
			Name:              v.Name,
			MaxRTT:            time.Duration(v.MaxRTT) * time.Millisecond,
			Interval:          time.Duration(v.Interval) * time.Millisecond,
			Count:             v.Count,
			Timeout:           time.Duration(v.Timeout) * time.Millisecond,
			NotificationCount: v.NotificationCount,
			PeriodStart:       time.Duration(v.PeriodStart) * time.Minute,
			PeriodEnd:         time.Duration(v.PeriodEnd) * time.Minute,
			Enabled:           v.Enabled,
			Created:           v.Created,
		})
	}

	return data, nil
}

func (r *AddressRepo) Create(ctx context.Context, dto *models.AddressDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, ip, name, max_rtt, interval, count, timeout, not_count, period_start, period_end, enabled) 
		VALUES (:id, :ip, :name, :max_rtt, :interval, :count, :timeout, :not_count, :period_start, :period_end, :enabled)`,
		AddressTable,
	)
	dto.ID = uuid.NewString()

	data := pq_models.AddressDTO{
		ID:                dto.ID,
		IP:                dto.IP,
		Name:              dto.Name,
		Count:             dto.Count,
		NotificationCount: dto.NotificationCount,
		Enabled:           dto.Enabled,
	}
	times := [5]int64{}
	if dto.MaxRTT != nil {
		times[0] = dto.MaxRTT.Milliseconds()
		data.MaxRTT = &times[0]
	}
	if dto.Interval != nil {
		times[1] = dto.Interval.Milliseconds()
		data.Interval = &times[1]
	}
	if dto.Timeout != nil {
		times[2] = dto.Timeout.Milliseconds()
		data.Timeout = &times[2]
	}
	if dto.PeriodStart != nil {
		times[3] = int64(dto.PeriodStart.Minutes())
		data.PeriodStart = &times[3]
	}
	if dto.PeriodEnd != nil {
		times[4] = int64(dto.PeriodEnd.Minutes())
		data.PeriodEnd = &times[4]
	}

	_, err := r.db.NamedExecContext(ctx, query, data)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *AddressRepo) Update(ctx context.Context, dto *models.AddressDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET name = :name, max_rtt = :max_rtt, interval = :interval, count = :count, timeout = :timeout, 
		not_count = :not_count, period_start = :period_start, period_end = :period_end, enabled = :enabled WHERE ip = :ip`,
		AddressTable,
	)

	data := pq_models.AddressDTO{
		ID:                dto.ID,
		IP:                dto.IP,
		Name:              dto.Name,
		Count:             dto.Count,
		NotificationCount: dto.NotificationCount,
		Enabled:           dto.Enabled,
	}
	times := [5]int64{}
	if dto.MaxRTT != nil {
		times[0] = dto.MaxRTT.Milliseconds()
		data.MaxRTT = &times[0]
	}
	if dto.Interval != nil {
		times[1] = dto.Interval.Milliseconds()
		data.Interval = &times[1]
	}
	if dto.Timeout != nil {
		times[2] = dto.Timeout.Milliseconds()
		data.Timeout = &times[2]
	}
	if dto.PeriodStart != nil {
		times[3] = int64(dto.PeriodStart.Minutes())
		data.PeriodStart = &times[3]
	}
	if dto.PeriodEnd != nil {
		times[4] = int64(dto.PeriodEnd.Minutes())
		data.PeriodEnd = &times[4]
	}

	_, err := r.db.NamedExecContext(ctx, query, data)
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
