package models

import "time"

type Scheduler struct {
	Interval time.Duration `json:"interval" db:"interval"` // интервал запуска cron
	// MaxCount int           `json:"max_count" db:"max_count"` // количество одновременных пингов
}
