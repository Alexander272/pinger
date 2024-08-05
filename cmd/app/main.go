package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Alexander272/Pinger/internal/bot"
	"github.com/Alexander272/Pinger/internal/config"
	"github.com/Alexander272/Pinger/internal/ping"
	"github.com/Alexander272/Pinger/pkg/logger"
	"github.com/Alexander272/Pinger/pkg/mattermost"
	"github.com/go-co-op/gocron"
)

func main() {
	// if err := gotenv.Load(".env"); err != nil {
	// 	logger.Fatalf("error loading env variables: %s", err.Error())
	// }

	conf, err := config.Init("configs/config.yaml")
	if err != nil {
		logger.Fatalf("error initializing configs: %s", err.Error())
	}
	logger.Init(os.Stdout, conf.Environment)

	mostClient := mattermost.NewMattermostClient(mattermost.Config{Server: conf.Bot.Server, Token: conf.Bot.Token})

	// mostClient.CreateDirectChannel()
	botClient := bot.NewMostBotClient(&conf.Bot, mostClient)

	pingClient := ping.NewPingClient(&conf.Pinger, botClient)

	cron := gocron.NewScheduler(time.UTC)

	for _, ac := range conf.Pinger.Addresses {
		logger.Debug(ac)
		logger.Infof("started cron. interval: %s", ac.Interval.String())
		_, err = cron.Every(ac.Interval).Do(func(ac *config.AddressesConfig) {
			logger.Debug("started ping")
			pingClient.Ping(ac.List)
		}, ac)
		if err != nil {
			logger.Fatalf("failed to start cron. err: %s", err.Error())
		}
	}

	logger.Debug("jobs ", len(cron.Jobs()))
	cron.StartAsync()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	cron.Clear()
	cron.Stop()

	const timeout = 5 * time.Second

	_, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()
}
