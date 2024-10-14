package models

import "time"

type Address struct {
	ID                string        `json:"id" db:"id"`
	IP                string        `json:"ip" db:"ip"`
	Name              string        `json:"name" db:"name"`
	MaxRTT            time.Duration `json:"maxRtt" db:"max_rtt"`
	Interval          time.Duration `json:"interval" db:"interval"`           // Интервал - время ожидания между отправкой каждого пакета.
	Count             int           `json:"count" db:"count"`                 // Count указывает pinger на остановку после отправки (и получения) Count эхо-пакетов
	Timeout           time.Duration `json:"timeout" db:"timeout"`             // Timeout задает таймаут до завершения ping
	NotificationCount int           `json:"notificationCount" db:"not_count"` // Количество уведомлений
	PeriodStart       time.Duration `json:"periodStart" db:"period_start"`
	PeriodEnd         time.Duration `json:"periodEnd" db:"period_end"`
	Enabled           bool          `json:"enabled" db:"enabled"`
	Created           time.Time     `json:"created" db:"created_at"`
}

type AddressDTO struct {
	ID                string         `json:"id" db:"id"`
	IP                string         `json:"ip" db:"ip"`
	Name              *string        `json:"name" db:"name"`
	MaxRTT            *time.Duration `json:"maxRtt" db:"max_rtt"`
	Interval          *time.Duration `json:"interval" db:"interval"`
	Count             *int           `json:"count" db:"count"`
	Timeout           *time.Duration `json:"timeout" db:"timeout"`
	NotificationCount *int           `json:"notificationCount" db:"not_count"`
	PeriodStart       *time.Duration `json:"periodStart" db:"period_start"`
	PeriodEnd         *time.Duration `json:"periodEnd" db:"period_end"`
	Enabled           *bool          `json:"enabled" db:"enabled"`
}

type Statistic struct {
	IP              string `json:"ip"`
	IsLong          bool   `json:"isLong"`
	IsFailed        bool   `json:"isFailed"`
	MaxNotification int    `json:"maxNotification"`
}
