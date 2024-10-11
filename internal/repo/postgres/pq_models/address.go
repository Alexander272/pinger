package pq_models

import "time"

type Address struct {
	ID                string    `db:"id"`
	IP                string    `db:"ip"`
	Name              string    `db:"name"`
	MaxRTT            int64     `db:"max_rtt"`
	Interval          int64     `db:"interval"`
	Count             int       `db:"count"`
	Timeout           int64     `db:"timeout"`
	NotificationCount int       `db:"not_count"`
	PeriodStart       int64     `db:"period_start"`
	PeriodEnd         int64     `db:"period_end"`
	Enabled           bool      `db:"enabled"`
	Created           time.Time `json:"created" db:"created_at"`
}

type AddressDTO struct {
	ID                string  `db:"id"`
	IP                string  `db:"ip"`
	Name              *string `db:"name"`
	MaxRTT            *int64  `db:"max_rtt"`
	Interval          *int64  `db:"interval"`
	Count             *int    `db:"count"`
	Timeout           *int64  `db:"timeout"`
	NotificationCount *int    `db:"not_count"`
	PeriodStart       *int64  `db:"period_start"`
	PeriodEnd         *int64  `db:"period_end"`
	Enabled           *bool   `db:"enabled"`
}
