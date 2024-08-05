package ping

import (
	"fmt"

	"github.com/Alexander272/Pinger/internal/bot"
	"github.com/Alexander272/Pinger/internal/config"
	"github.com/Alexander272/Pinger/pkg/logger"
	probing "github.com/prometheus-community/pro-bing"
)

type PingClient struct {
	conf *config.PingerConfig
	bot  bot.MostBot
	lost map[string]int
	long map[string]int
}

func NewPingClient(conf *config.PingerConfig, bot bot.MostBot) *PingClient {
	lost := make(map[string]int, 0)
	long := make(map[string]int, 0)

	return &PingClient{
		conf: conf,
		bot:  bot,
		lost: lost,
		long: long,
	}
}

func (c *PingClient) Ping(addresses []*config.Address) {

	for _, addr := range addresses {
		go c.checkPing(addr)
	}
}

func (c *PingClient) checkPing(addr *config.Address) {
	logger.Debug("addr ", addr)

	pinger, err := probing.NewPinger(addr.Ip)
	if err != nil {
		logger.Errorf("failed to create new pinger. error: %s", err.Error())

		message := fmt.Sprintf(`Возникла ошибка: %s`, err.Error())
		if err := c.bot.SendErr(addr, message); err != nil {
			logger.Errorf("failed to send message. error: %s", err.Error())
		}
	}

	pinger.Count = c.conf.Count
	pinger.Interval = c.conf.Interval
	pinger.Timeout = c.conf.Timeout

	pinger.OnRecv = func(pkt *probing.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}

	pinger.OnDuplicateRecv = func(pkt *probing.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v (DUP!)\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.TTL)
	}

	pinger.OnFinish = func(stats *probing.Statistics) {
		fmt.Printf("\n--- ping statistics. from %s to %s ---\n", c.conf.IP, stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}

	err = pinger.Run() // Blocks until finished.
	if err != nil {
		logger.Errorf("failed to rung pinger. error: %s", err.Error())

		message := fmt.Sprintf(`Возникла ошибка: %s`, err.Error())
		if err := c.bot.SendErr(addr, message); err != nil {
			logger.Errorf("failed to send message. error: %s", err.Error())
		}
	}

	stats := pinger.Statistics()

	// 	logger.Debug(
	// 		fmt.Sprintf(`
	// --- %s ping statistics ---
	// %d packets transmitted, %d packets received, %v%% packet loss
	// round-trip min/avg/max/stddev = %v/%v/%v/%v`,
	// 			addr,
	// 			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss,
	// 			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt,
	// 		))

	// stats.AvgRtt
	if addr.Rtt != 0 {
		logger.Debug("rtt ", addr.Ip, " ", addr.Name)
		if stats.AvgRtt > addr.Rtt {
			if c.long[addr.Ip] < 3 {
				c.long[addr.Ip] += 1
				message := fmt.Sprintf("Превышено допустимое время пинга **(%s)** для IP **%s (%s)**", stats.AvgRtt.String(), addr.Ip, addr.Name)

				if err := c.bot.LongPing(message); err != nil {
					logger.Errorf("failed to send message. error: %s", err.Error())
				}
			}
		} else {
			if c.long[addr.Ip] != 0 {
				message := fmt.Sprintf("Время пинга **(%s)** для IP **%s (%s)** в норме", stats.AvgRtt.String(), addr.Ip, addr.Name)

				if err := c.bot.LongPing(message); err != nil {
					logger.Errorf("failed to send message. error: %s", err.Error())
				}
			}
			c.long[addr.Ip] = 0
		}
	}

	// logger.Debug("lost ", addr.Ip, " ", c.lost[addr.Ip])
	if stats.PacketLoss > 50 {
		if c.lost[addr.Ip] < 3 {
			c.lost[addr.Ip] += 1
			// 		statistics := fmt.Sprintf(`
			// --- %s ping statistics ---
			// %d packets transmitted, %d packets received, %v%% packet loss
			// round-trip min/avg/max/stddev = %v/%v/%v/%v
			// `,
			// 			addr,
			// 			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss,
			// 			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt,
			// 		)

			statistics := fmt.Sprintf(`--- ping statistics. from %s to %s ---
%d packets transmitted, %d packets received, %v%% packet loss`,
				c.conf.IP, addr.Ip, stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss,
			)

			if err := c.bot.SendErr(addr, statistics); err != nil {
				logger.Errorf("failed to send message. error: %s", err.Error())
			}

			// logger.Info(statistics)
		}
	} else {
		if c.lost[addr.Ip] != 0 {
			if err := c.bot.Send(addr); err != nil {
				logger.Errorf("failed to send message. error: %s", err.Error())
			}
		}
		c.lost[addr.Ip] = 0
	}

}
