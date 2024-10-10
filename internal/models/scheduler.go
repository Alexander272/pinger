package models

import "time"

type Scheduler struct {
	Interval time.Duration `json:"interval" db:"interval"`
	MaxCount int           `json:"max_count" db:"max_count"`
}
