package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Alexander272/Pinger/internal/models"
	"github.com/Alexander272/Pinger/pkg/error_bot"
	"github.com/Alexander272/Pinger/pkg/logger"
	"github.com/gin-gonic/gin"
	probing "github.com/prometheus-community/pro-bing"
)

type PingService struct {
	addresses Address
	stats     Statistic
	post      Post

	failed *models.Counters
	long   *models.Counters
}

type PingDeps struct {
	Address Address
	Stats   Statistic
	Post    Post
}

func NewPingService(deps *PingDeps) *PingService {
	return &PingService{
		addresses: deps.Address,
		stats:     deps.Stats,
		post:      deps.Post,

		failed: models.NewCounters(),
		long:   models.NewCounters(),
	}
}

type Ping interface {
	Ping(addr *models.Address) (*models.PingStatistic, error)
	CheckPing(hostIP string)
}

func (s *PingService) Ping(addr *models.Address) (*models.PingStatistic, error) {
	logger.Debug("ping", logger.AnyAttr("addr", addr))

	pinger, err := probing.NewPinger(addr.IP)
	if err != nil {
		logger.Error("failed to create new pinger.", logger.ErrAttr(err))
		return nil, fmt.Errorf("failed to create new pinger. error: %w", err)
	}

	pinger.Count = addr.Count
	pinger.Interval = addr.Interval
	pinger.Timeout = addr.Timeout

	err = pinger.Run() // Blocks until finished.
	if err != nil {
		logger.Error("failed to run pinger.", logger.ErrAttr(err))
		return nil, fmt.Errorf("failed to run pinger. error: %w", err)
	}

	stats := pinger.Statistics()
	statistic := &models.PingStatistic{
		IP:              addr.IP,
		IsFailed:        stats.PacketLoss > 50,
		MaxNotification: addr.NotificationCount,
	}
	if addr.MaxRTT != 0 {
		statistic.IsLong = stats.AvgRtt >= addr.MaxRTT
	}

	return statistic, nil
}

func (s *PingService) SendPing(addr *models.Address, hostIP string) {
	pinger, err := probing.NewPinger(addr.IP)
	if err != nil {
		logger.Error("failed to create new pinger.", logger.ErrAttr(err))
		error_bot.Send(&gin.Context{}, err.Error(), nil)
		s.post.Send(&models.Post{Message: "Произошла ошибка при создании экземпляра pinger."})
		return
	}

	pinger.Count = addr.Count
	pinger.Interval = addr.Interval
	pinger.Timeout = addr.Timeout

	err = pinger.Run() // Blocks until finished.
	if err != nil {
		logger.Error("failed to run pinger.", logger.ErrAttr(err))
		error_bot.Send(&gin.Context{}, err.Error(), nil)

		s.post.Send(&models.Post{Message: "Произошла ошибка при запуске pinger."})
		return
	}

	stats := pinger.Statistics()

	if stats.PacketLoss > 50 {
		count, ok := s.failed.Load(addr.IP)
		if count == 0 {
			stats := &models.StatisticDTO{
				IP:        addr.IP,
				Name:      addr.Name,
				TimeStart: time.Now(),
			}
			if err := s.stats.Create(context.Background(), stats); err != nil {
				error_bot.Send(&gin.Context{}, err.Error(), stats)
			}
		}

		if addr.NotificationCount == 0 || !ok || count < addr.NotificationCount {
			s.failed.Inc(addr.IP)

			statistics := fmt.Sprintf("--- ping statistics. from %s to %s ---\n%d packets transmitted, %d packets received, %v%% packet loss",
				hostIP, addr.IP, stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss,
			)
			message := fmt.Sprintf("Пинг по адресу **%s (%s)** не прошел.\n```\n%s\n```", addr.IP, addr.Name, statistics)
			s.post.Send(&models.Post{Message: message})
		}
		return
	}

	count, ok := s.failed.Load(addr.IP)
	if ok && count != 0 {
		message := fmt.Sprintf("Пинг по адресу **%s (%s)** прошел.", addr.IP, addr.Name)
		s.post.Send(&models.Post{Message: message})
		s.failed.Store(addr.IP, 0)

		stats := &models.StatisticDTO{IP: addr.IP, TimeEnd: time.Now()}
		if err := s.stats.Update(context.Background(), stats); err != nil {
			error_bot.Send(&gin.Context{}, err.Error(), stats)
		}
	}

	if addr.MaxRTT == 0 {
		return
	}

	if stats.AvgRtt >= addr.MaxRTT {
		count, ok := s.long.Load(addr.IP)
		if addr.NotificationCount == 0 || !ok || count < addr.NotificationCount {
			s.long.Inc(addr.IP)

			message := fmt.Sprintf("Превышено допустимое время пинга **(%s)** для IP **%s (%s)**", stats.AvgRtt.String(), addr.IP, addr.Name)
			s.post.Send(&models.Post{Message: message})
		}
	} else {
		count, ok := s.long.Load(addr.IP)
		if ok && count != 0 {
			message := fmt.Sprintf("Время пинга **(%s)** для IP **%s (%s)** в норме", stats.AvgRtt.String(), addr.IP, addr.Name)
			s.post.Send(&models.Post{Message: message})
			s.long.Store(addr.IP, 0)
		}
	}
}

func (s *PingService) CheckPing(hostIP string) {
	addresses, err := s.addresses.Get(context.Background())
	if err != nil {
		logger.Error("failed to get addresses.", logger.ErrAttr(err))
		error_bot.Send(&gin.Context{}, err.Error(), nil)
		s.post.Send(&models.Post{Message: "Произошла ошибка при получении адресов."})
		return
	}

	// urls := make(chan struct{}, 20)

	now := time.Now()
	for _, address := range addresses {
		isAfter, isBefore := true, true
		if address.PeriodStart != 0 && address.PeriodEnd != 0 {
			isAfter = now.After(time.Date(now.Year(), now.Month(), now.Day(), 0, int(address.PeriodStart.Minutes()), 0, 0, now.Location()))
			isBefore = now.Before(time.Date(now.Year(), now.Month(), now.Day(), 0, int(address.PeriodEnd.Minutes()), 0, 0, now.Location()))
		}
		if !isAfter || !isBefore {
			continue
		}

		logger.Debug("ping", logger.AnyAttr("addr", address))
		go s.SendPing(address, hostIP)
	}
}
