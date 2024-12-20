package models

import "time"

type Statistic struct {
	ID        string        `json:"id" db:"id"`
	IP        string        `json:"ip" db:"ip"`
	Name      string        `json:"name" db:"name"`
	Time      time.Duration `json:"time" db:"time"`
	TimeStart time.Time     `json:"timeStart" db:"time_start"`
	TimeEnd   time.Time     `json:"timeEnd" db:"time_end"`
	Created   time.Time     `json:"created" db:"created_at"`
}

type GetStatisticDTO struct {
	PeriodStart time.Time `json:"periodStart" db:"period_start"`
	PeriodEnd   time.Time `json:"periodEnd" db:"period_end"`
}

type GetStatisticByIPDTO struct {
	IP          string    `json:"ip" db:"ip"`
	PeriodStart time.Time `json:"periodStart" db:"period_start"`
	PeriodEnd   time.Time `json:"periodEnd" db:"period_end"`
}

type GetUnavailableDTO struct{}

type StatisticDTO struct {
	ID        string    `json:"id" db:"id"`
	IP        string    `json:"ip" db:"ip"`
	Name      string    `json:"name" db:"name"`
	TimeStart time.Time `json:"timeStart" db:"time_start"`
	TimeEnd   time.Time `json:"timeEnd" db:"time_end"`
	Created   time.Time `json:"created" db:"created_at"`
}
