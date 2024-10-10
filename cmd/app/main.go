package main

import (
	"log"

	"github.com/Alexander272/Pinger/internal/config"
	"github.com/Alexander272/Pinger/pkg/logger"
	"github.com/Alexander272/Pinger/pkg/mattermost"
	"github.com/subosito/gotenv"
)

func main() {
	if err := gotenv.Load(".env"); err != nil {
		log.Fatalf("failed to load env variables. error: %s", err.Error())
	}

	conf, err := config.Init("configs/config.yaml")
	if err != nil {
		log.Fatalf("failed to init configs. error: %s", err.Error())
	}
	logger.NewLogger(logger.WithLevel(conf.LogLevel), logger.WithAddSource(conf.LogSource))

	//* Dependencies
	mattermostConf := mattermost.Config{
		ServerLink: conf.Bot.Server,
		Token:      conf.Bot.Token,
	}
	mostClient := mattermost.NewMattermostClient(mattermostConf)

	_, _, err = mostClient.Http.GetPing()
	if err != nil {
		log.Fatalf("failed to ping most. error: %s", err.Error())
	}

	bot, _, err := mostClient.Http.GetMe("")
	if err != nil {
		log.Fatalf("failed to get bot data. error: %s", err.Error())
	}
	logger.Debug("me", logger.AnyAttr("bot", bot))

	//* Services, Repos & API Handlers
	// servicesDeps := services.Deps{
	// 	MostClient: mostClient.Http,
	// 	BotName:    conf.Most.BotName,
	// }
	// services := services.NewServices(servicesDeps)
	// handlers := transport.NewHandler(services)

	// socHandler := socket.NewHandler(mostClient.Socket, bot)

	// socHandler.Listen()

	//* HTTP Server
	// srv := server.NewServer(conf, handlers.Init(conf))
	// go func() {
	// 	if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
	// 		logger.Fatalf("error occurred while running http server: %s\n", err.Error())
	// 	}
	// }()
	// logger.Infof("Application started on port: %s", conf.Http.Port)

	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	// <-quit

	// const timeout = 5 * time.Second
	// ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	// defer shutdown()

	// socHandler.Close()

	// if err := srv.Stop(ctx); err != nil {
	// 	logger.Errorf("failed to stop server. error: %s", err.Error())
	// }
}
