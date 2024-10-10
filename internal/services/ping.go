package services

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Alexander272/Pinger/internal/models"
	"github.com/Alexander272/Pinger/pkg/error_bot"
	"github.com/Alexander272/Pinger/pkg/logger"
	"github.com/gin-gonic/gin"
	probing "github.com/prometheus-community/pro-bing"
)

type PingService struct {
	addresses Address
	post      Post

	failed *models.Counters
	long   *models.Counters
}

func NewPingService(addresses Address, post Post) *PingService {
	return &PingService{
		addresses: addresses,
		post:      post,

		failed: models.NewCounters(),
		long:   models.NewCounters(),
	}
}

type Ping interface {
	Ping(addr *models.Address) (*models.Statistic, error)
}

func (s *PingService) Ping(addr *models.Address) (*models.Statistic, error) {
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
	statistic := &models.Statistic{
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
	logger.Debug("ping", logger.AnyAttr("addr", addr))
	ginCtx := &gin.Context{
		Request: &http.Request{
			Method: "Get",
		},
	}

	pinger, err := probing.NewPinger(addr.IP)
	if err != nil {
		logger.Error("failed to create new pinger.", logger.ErrAttr(err))
		ginCtx.Request.URL = &url.URL{Host: "ping-bot", Path: "create pinger"}
		error_bot.Send(ginCtx, err.Error(), nil)
		s.post.Send(&models.Post{Message: "Произошла ошибка при создании экземпляра pinger."})
		return
	}

	pinger.Count = addr.Count
	pinger.Interval = addr.Interval
	pinger.Timeout = addr.Timeout

	err = pinger.Run() // Blocks until finished.
	if err != nil {
		logger.Error("failed to run pinger.", logger.ErrAttr(err))
		ginCtx.Request.URL = &url.URL{Host: "ping-bot", Path: "run pinger"}
		error_bot.Send(ginCtx, err.Error(), nil)

		s.post.Send(&models.Post{Message: "Произошла ошибка при запуске pinger."})
		return
	}

	stats := pinger.Statistics()

	if stats.PacketLoss > 50 {
		count, ok := s.failed.Load(addr.IP)
		if !ok || count < addr.NotificationCount {
			s.failed.Inc(addr.IP)

			statistics := fmt.Sprintf(`--- ping statistics. from %s to %s ---\n%d packets transmitted, %d packets received, %v%% packet loss`,
				hostIP, addr.IP, stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss,
			)
			message := fmt.Sprintf("Пинг по адресу **%s (%s)** не прошел.\n```\n%s\n```", addr.IP, addr.Name, statistics)
			s.post.Send(&models.Post{Message: message})
			return
		}
	}

	count, ok := s.failed.Load(addr.IP)
	if ok && count != 0 {
		message := fmt.Sprintf("Пинг по адресу **%s (%s)** прошел.", addr.IP, addr.Name)
		s.post.Send(&models.Post{Message: message})
		s.failed.Store(addr.IP, 0)
	}

	if addr.MaxRTT != 0 && stats.AvgRtt >= addr.MaxRTT {
		count, ok := s.long.Load(addr.IP)
		if !ok || count < addr.NotificationCount {
			s.long.Inc(addr.IP)

			message := fmt.Sprintf("Превышено допустимое время пинга **(%s)** для IP **%s (%s)**", stats.AvgRtt.String(), addr.IP, addr.Name)
			s.post.Send(&models.Post{Message: message})
		}
	} else if addr.MaxRTT != 0 {
		count, ok := s.long.Load(addr.IP)
		if ok && count != 0 {
			message := fmt.Sprintf("Время пинга **(%s)** для IP **%s (%s)** в норме", stats.AvgRtt.String(), addr.IP, addr.Name)
			s.post.Send(&models.Post{Message: message})
			s.long.Store(addr.IP, 0)
		}
	}
}

func (s *PingService) CheckPing(hostIP string) {
	ginCtx := &gin.Context{
		Request: &http.Request{
			Method: "Get",
			// URL:    &url.URL{Host: "cron", Path: "delete-old-reagent"},
		},
	}

	addresses, err := s.addresses.Get(context.Background())
	if err != nil {
		logger.Error("failed to get addresses.", logger.ErrAttr(err))
		ginCtx.Request.URL = &url.URL{Host: "ping-bot", Path: "get addresses"}
		error_bot.Send(ginCtx, err.Error(), nil)
		s.post.Send(&models.Post{Message: "Произошла ошибка при получении адресов."})
		return
	}

	now := time.Now()
	for _, address := range addresses {
		isAfter := now.After(time.Date(now.Year(), now.Month(), now.Day(), 0, int(address.PeriodStart.Minutes()), 0, 0, now.Location()))
		isBefore := now.Before(time.Date(now.Year(), now.Month(), now.Day(), 0, int(address.PeriodEnd.Minutes()), 0, 0, now.Location()))

		if !isAfter || !isBefore {
			continue
		}

		go s.SendPing(address, hostIP)

		// go (func() {
		// 	stats, err := s.Ping(address)
		// 	if err != nil {
		// 		logger.Error("failed to ping address.", logger.ErrAttr(err))
		// 		ginCtx.Request.URL = &url.URL{Host: "ping-bot", Path: "ping address"}
		// 		error_bot.Send(ginCtx, err.Error(), nil)
		// 		return
		// 	}

		// 	// if stats.IsFailed {
		// 	// 	s.post.Send(&models.Post{Message: fmt.Sprintf("Пинг по адресу **%s (%s)** не прошел.\n```\n%s\n```", stats.IP, stats.Name, message)})
		// 	// }
		// })()

		// statistic, err := s.Ping(address)
		// if err != nil {
		// 	logger.Error("failed to ping address.", logger.ErrAttr(err))
		// 	ginCtx.Request.URL = &url.URL{Host: "ping-bot", Path: "ping address"}
		// 	error_bot.Send(ginCtx, err.Error(), nil)
		// 	continue
		// }
	}
}
