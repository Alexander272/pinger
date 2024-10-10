package services

import (
	"fmt"
	"log"
	"time"

	"github.com/Alexander272/Pinger/internal/utils"
	"github.com/go-co-op/gocron/v2"
)

type SchedulerService struct {
	cron gocron.Scheduler
}

func NewSchedulerService() *SchedulerService {
	cron, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("failed to create new scheduler. error: %s", err.Error())
	}

	return &SchedulerService{
		cron: cron,
	}
}

type Scheduler interface {
	Start() error
	Restart() error
	Stop() error
}

func (s *SchedulerService) Start() error {
	// now := time.Now()
	// jobStart := time.Date(now.Year(), now.Month(), now.Day(), conf.StartTime, 0, 0, 0, now.Location())
	// if now.Hour() >= conf.StartTime {
	// 	jobStart = jobStart.Add(24 * time.Hour)
	// }
	// // jobStart := now.Add(1 * time.Minute)
	// logger.Info("start time of job " + jobStart.Format("02.01.2006 15:04:05"))

	hostIP := utils.GetOutboundIP()

	// job := gocron.DurationJob(conf.Interval)
	job := gocron.DurationJob(1 * time.Minute)
	task := gocron.NewTask(s.job, hostIP.String())
	// jobStartAt := gocron.WithStartAt(gocron.WithStartDateTime(jobStart))
	jobStartAt := gocron.WithStartAt(gocron.WithStartDateTime(time.Now()))

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
	// TODO
}
