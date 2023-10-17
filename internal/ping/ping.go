package ping

import (
	"fmt"
	"time"

	"github.com/Alexander272/Pinger/internal/bot"
	"github.com/Alexander272/Pinger/internal/config"
	"github.com/Alexander272/Pinger/pkg/logger"
	probing "github.com/prometheus-community/pro-bing"
)

type PingClient struct {
	conf *config.PingerConfig
	bot  bot.MostBot
}

func NewPingClient(conf *config.PingerConfig, bot bot.MostBot) *PingClient {
	return &PingClient{
		conf: conf,
		bot:  bot,
	}
}

func (c *PingClient) Ping(addresses []string) {
	for _, addr := range addresses {
		go c.checkPing(addr)
	}
}

func (c *PingClient) checkPing(addr string) {
	logger.Debug("addr ", addr)

	pinger, err := probing.NewPinger(addr)
	if err != nil {
		logger.Errorf("failed to create new pinger. error: %s", err.Error())

		message := fmt.Sprintf(`Возникла ошибка: %s`, err.Error())
		if err := c.bot.Send(addr, message); err != nil {
			logger.Errorf("failed to send message. error: %s", err.Error())
		}
	}
	pinger.Count = 5
	pinger.Timeout = time.Second

	err = pinger.Run() // Blocks until finished.
	if err != nil {
		logger.Errorf("failed to rung pinger. error: %s", err.Error())

		message := fmt.Sprintf(`Возникла ошибка: %s`, err.Error())
		if err := c.bot.Send(addr, message); err != nil {
			logger.Errorf("failed to send message. error: %s", err.Error())
		}
	}

	stats := pinger.Statistics()

	if stats.PacketLoss > 50 {
		// 		statistics := fmt.Sprintf(`
		// --- %s ping statistics ---
		// %d packets transmitted, %d packets received, %v%% packet loss
		// round-trip min/avg/max/stddev = %v/%v/%v/%v
		// `,
		// 			addr,
		// 			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss,
		// 			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt,
		// 		)

		statistics := fmt.Sprintf(`--- %s ping statistics ---
%d packets transmitted, %d packets received, %v%% packet loss`,
			addr, stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss,
		)

		if err := c.bot.Send(addr, statistics); err != nil {
			logger.Errorf("failed to send message. error: %s", err.Error())
		}

		// logger.Info(statistics)
	}

}
