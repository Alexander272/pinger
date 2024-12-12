package services

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Alexander272/Pinger/pkg/mattermost"
	"github.com/go-co-op/gocron/v2"
)

type SchedulerService struct {
	cron   gocron.Scheduler
	ping   Ping
	client *mattermost.Client
}

func NewSchedulerService(ping Ping, client *mattermost.Client) *SchedulerService {
	cron, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("failed to create new scheduler. error: %s", err.Error())
	}

	return &SchedulerService{
		cron:   cron,
		ping:   ping,
		client: client,
	}
}

type Scheduler interface {
	Start() error
	Restart() error
	Stop() error
}

func (s *SchedulerService) Start() error {
	// hostIP := utils.GetOutboundIP().String()
	// поскольку я запускаю бота через docker compose, выполняя команду выше я получаю ip контейнера, а не хоста. Поэтому приходится задавать ip через env
	hostIP := os.Getenv("HOST_IP")
	jobStart := time.Now().Add(1 * time.Minute)

	// job := gocron.DurationJob(conf.Interval)
	job := gocron.DurationJob(1 * time.Minute)
	task := gocron.NewTask(s.job, hostIP)
	jobStartAt := gocron.WithStartAt(gocron.WithStartDateTime(jobStart))

	_, err := s.cron.NewJob(job, task, jobStartAt)
	if err != nil {
		return fmt.Errorf("failed to create new job. error: %w", err)
	}

	//? запуск крона через интервал
	s.cron.Start()
	return nil
}

func (s *SchedulerService) Restart() error {
	if err := s.Stop(); err != nil {
		return err
	}
	if err := s.Start(); err != nil {
		return err
	}
	return nil
}

func (s *SchedulerService) Stop() error {
	if err := s.cron.Shutdown(); err != nil {
		return fmt.Errorf("failed to shutdown cron scheduler. error: %w", err)
	}
	return nil
}

func (s *SchedulerService) job(hostIP string) {
	s.ping.CheckPing(hostIP)

	if !s.client.IsConnected() {
		ok := s.client.Reconnect()
		if ok {
			s.client.Socket.Listen()
		}
	}
}
