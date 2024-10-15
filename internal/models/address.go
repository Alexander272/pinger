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

type Decode struct {
	Name              *string `short:"n" long:"name" description:"Название адреса"`
	MaxRTT            int     `short:"r" long:"rtt" description:"Допустимое время пинга в миллисекундах"`
	Interval          int     `short:"i" long:"interval" description:"Show verbose debug information"`
	Timeout           int     `short:"t" long:"timeout" description:"Show verbose debug information"`
	Count             *int    `short:"c" long:"count" description:"Show verbose debug information"`
	NotificationCount *int    `short:"N" long:"notification" description:"Show verbose debug information"`
	Period            string  `short:"p" long:"period" description:"Show verbose debug information"`
	// Enabled           *bool   `short:"v" long:"verbose" description:"Show verbose debug information"`
}
