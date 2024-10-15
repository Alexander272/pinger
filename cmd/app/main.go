package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Alexander272/Pinger/internal/config"
	"github.com/Alexander272/Pinger/internal/repo"
	"github.com/Alexander272/Pinger/internal/services"
	"github.com/Alexander272/Pinger/internal/transport/socket"
	"github.com/Alexander272/Pinger/pkg/database/postgres"
	"github.com/Alexander272/Pinger/pkg/logger"
	"github.com/Alexander272/Pinger/pkg/mattermost"
	_ "github.com/lib/pq"
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
	db, err := postgres.NewPostgresDB(postgres.Config{
		Host:     conf.Postgres.Host,
		Port:     conf.Postgres.Port,
		Username: conf.Postgres.Username,
		Password: conf.Postgres.Password,
		DBName:   conf.Postgres.DbName,
		SSLMode:  conf.Postgres.SSLMode,
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

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
	repos := repo.NewRepository(db)

	servicesDeps := &services.Deps{
		Repo:      repos,
		Client:    mostClient.Http,
		ChannelID: conf.Bot.ChannelId,
	}
	services := services.NewServices(servicesDeps)
	// handlers := transport.NewHandler(services)
	socHandler := socket.NewHandler(&socket.Deps{Socket: mostClient.Socket, User: bot, Services: services})

	if err := services.Scheduler.Start(); err != nil {
		log.Fatalf("failed to start scheduler. error: %s\n", err.Error())
	}

	//* HTTP Server
	// srv := server.NewServer(conf, handlers.Init(conf))
	// go func() {
	// 	if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
	// 		logger.Fatalf("error occurred while running http server: %s\n", err.Error())
	// 	}
	// }()
	// logger.Infof("Application started on port: %s", conf.Http.Port)

	go func() {
		// TODO при ошибке приложение падает
		socHandler.Listen()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	if err := services.Scheduler.Stop(); err != nil {
		logger.Error("failed to stop sending notification.", logger.ErrAttr(err))
	}

	// const timeout = 5 * time.Second
	// ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	// defer shutdown()

	socHandler.Close()

	// if err := srv.Stop(ctx); err != nil {
	// 	logger.Errorf("failed to stop server. error: %s", err.Error())
	// }
}
