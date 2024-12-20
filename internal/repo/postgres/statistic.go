package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Alexander272/Pinger/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type StatisticRepo struct {
	db *sqlx.DB
}

func NewStatisticRepo(db *sqlx.DB) *StatisticRepo {
	return &StatisticRepo{db: db}
}

type Statistic interface {
	Get(ctx context.Context, req *models.GetStatisticDTO) ([]*models.Statistic, error)
	GetByIP(ctx context.Context, req *models.GetStatisticByIPDTO) ([]*models.Statistic, error)
	GetUnavailable(ctx context.Context, req *models.GetUnavailableDTO) ([]*models.Statistic, error)
	GetLast(ctx context.Context, req *models.GetStatisticByIPDTO) (*models.Statistic, error)
	Create(ctx context.Context, dto *models.StatisticDTO) error
	Update(ctx context.Context, dto *models.StatisticDTO) error
}

func (r *StatisticRepo) Get(ctx context.Context, req *models.GetStatisticDTO) ([]*models.Statistic, error) {
	// по умолчанию я хочу получать суммарное количество времени за месяц по каждому IP
	// но думаю, нужно еще предусмотреть возможность указания периода
	query := fmt.Sprintf(`SELECT ip, name, ROUND(SUM(extract (epoch from time_end - time_start))) AS time FROM %s 
		WHERE time_end IS NOT NULL AND time_start >= $1 AND time_start <= $2 GROUP BY ip, name ORDER BY ip`,
		StatisticTable,
	)
	data := []*models.Statistic{}

	err := r.db.SelectContext(ctx, &data, query, req.PeriodStart, req.PeriodEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for i := range data {
		data[i].Time = data[i].Time * time.Second
	}

	return data, nil
}

func (r *StatisticRepo) GetByIP(ctx context.Context, req *models.GetStatisticByIPDTO) ([]*models.Statistic, error) {
	// а еще вывести все даты простоя по одному IP, по умолчанию за месяц
	query := fmt.Sprintf(`SELECT id, ip, name, ROUND(extract (epoch from time_end - time_start)) AS time, time_start, time_end FROM %s 
		WHERE ip = $1 AND time_end IS NOT NULL AND time_start >= $2 AND time_start <= $3 ORDER BY time_start`,
		StatisticTable,
	)
	data := []*models.Statistic{}

	err := r.db.SelectContext(ctx, &data, query, req.IP, req.PeriodStart, req.PeriodEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for i := range data {
		data[i].Time = data[i].Time * time.Second
	}

	return data, nil
}

func (r *StatisticRepo) GetUnavailable(ctx context.Context, req *models.GetUnavailableDTO) ([]*models.Statistic, error) {
	query := fmt.Sprintf(`SELECT id, ip, name, time_start FROM %s WHERE time_end IS NULL ORDER BY time_start`,
		StatisticTable,
	)
	data := []*models.Statistic{}

	err := r.db.SelectContext(ctx, &data, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *StatisticRepo) GetLast(ctx context.Context, req *models.GetStatisticByIPDTO) (*models.Statistic, error) {
	query := fmt.Sprintf(`SELECT id, ip, name, time_start FROM %s 
		WHERE ip = $1 AND time_end IS NULL ORDER BY time_start DESC LIMIT 1`,
		StatisticTable,
	)
	data := &models.Statistic{}

	err := r.db.GetContext(ctx, data, query, req.IP)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRows
		}
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *StatisticRepo) Create(ctx context.Context, dto *models.StatisticDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, ip, name, time_start) VALUES (:id, :ip, :name, :time_start)`, StatisticTable)
	dto.ID = uuid.NewString()

	_, err := r.db.NamedExecContext(ctx, query, dto)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *StatisticRepo) Update(ctx context.Context, dto *models.StatisticDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET time_end=:time_end WHERE ip=:ip AND time_end IS	NULL`, StatisticTable)

	_, err := r.db.NamedExecContext(ctx, query, dto)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
