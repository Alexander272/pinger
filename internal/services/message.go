package services

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/Alexander272/Pinger/internal/models"
	"github.com/Alexander272/Pinger/pkg/logger"
)

type MessageService struct {
	addresses Address
	post      Post
}

func NewMessageService(addresses Address, post Post) *MessageService {
	return &MessageService{
		addresses: addresses,
		post:      post,
	}
}

type Message interface {
	List(message string) error
	Create(message string) error
	Update(message string) error
	Delete(message string) error
	ToggleActive(message string, isEnable bool) error
}

func (s *MessageService) List(message string) error {
	addresses, err := s.addresses.GetAll(context.Background())
	if err != nil {
		logger.Error("failed to get addresses.", logger.ErrAttr(err))
		s.post.Send(&models.Post{Message: "Произошла ошибка при получении адресов."})
		return err
	}

	isAll := false
	parts := strings.Split(message, " ")
	if len(parts) > 1 && (parts[1] == "--all" || parts[1] == "-a") {
		isAll = true
	}
	table := []string{
		"| IP-адрес | Название | Статус |",
		"|:----|:----|:--|",
		// "|:-:|:-:|:-:|",
	}
	if isAll {
		table = []string{
			"| IP-адрес | Название | Допустимое время пинга | Количество уведомлений | Период | Интервал отправки пакетов | Таймаут до завершения ping | Количество пакетов | Статус |",
			"|:----|:----|:--|:--|:--|:--|:--|:--|:--|",
		}
	}

	for _, address := range addresses {
		isEnable := "Активен"
		if !address.Enabled {
			isEnable = "Не активен"
		}

		if !isAll {
			table = append(table, fmt.Sprintf("|%s|%s|%s|", address.IP, address.Name, isEnable))
		} else {
			start := time.Date(0, 1, 1, 0, int(address.PeriodStart.Minutes()), 0, 0, time.UTC)
			end := time.Date(0, 1, 1, 0, int(address.PeriodEnd.Minutes()), 0, 0, time.UTC)
			period := fmt.Sprintf("%s-%s", start.Format("15:04"), end.Format("15:04"))
			table = append(table, fmt.Sprintf("|%s|%s|%d|%d|%s|%d|%d|%d|%s|",
				address.IP, address.Name, address.MaxRTT.Milliseconds(), address.Count, period, address.Interval.Milliseconds(),
				address.Timeout.Milliseconds(), address.NotificationCount, isEnable,
			))
		}
	}

	s.post.Send(&models.Post{Message: strings.Join(table, "\n")})
	return nil
}

func (s *MessageService) Create(message string) error {
	logger.Info("create ip", logger.StringAttr("message", message))
	address := s.decode(message)
	if err := s.addresses.Create(context.Background(), address); err != nil {
		s.post.Send(&models.Post{Message: "Не удалось создать адрес."})
		logger.Error("failed to create address.", logger.ErrAttr(err))
		return err
	}
	return nil
}

func (s *MessageService) Update(message string) error {
	logger.Info("update ip", logger.StringAttr("message", message))
	address := s.decode(message)

	data, err := s.addresses.GetByIP(context.Background(), address.IP)
	if err != nil {
		if errors.Is(err, models.ErrNoRows) {
			s.post.Send(&models.Post{Message: "Не найден указанный IP адрес."})
			return nil
		}
		s.post.Send(&models.Post{Message: "При получении адреса произошла ошибка"})
		logger.Error("failed to get address by ip.", logger.ErrAttr(err))
		return err
	}

	if address.Name == nil {
		address.Name = &data.Name
	}
	if address.MaxRTT == nil {
		address.MaxRTT = &data.MaxRTT
	}
	if address.NotificationCount == nil {
		address.NotificationCount = &data.NotificationCount
	}
	if address.PeriodStart == nil {
		address.PeriodStart = &data.PeriodStart
	}
	if address.PeriodEnd == nil {
		address.PeriodEnd = &data.PeriodEnd
	}
	if address.Interval == nil {
		address.Interval = &data.Interval
	}
	if address.Count == nil {
		address.Count = &data.Count
	}
	if address.Timeout == nil {
		address.Timeout = &data.Timeout
	}
	if address.Enabled == nil {
		address.Enabled = &data.Enabled
	}

	if err := s.addresses.Update(context.Background(), address); err != nil {
		s.post.Send(&models.Post{Message: "Не удалось обновить адрес."})
		logger.Error("failed to update address.", logger.ErrAttr(err))
		return err
	}
	return nil
}

func (s *MessageService) ToggleActive(message string, isEnable bool) error {
	logger.Info("toggle active ip", logger.StringAttr("message", message), logger.BoolAttr("isEnable", isEnable))
	parts := strings.Split(message, " ")
	if net.ParseIP(parts[1]) == nil {
		s.post.Send(&models.Post{Message: "Не удалось распознать команду. Неправильный IP адрес."})
		return nil
	}

	//TODO наверное надо возвращать ошибку и отправлять ее в бота

	data, err := s.addresses.GetByIP(context.Background(), parts[1])
	if err != nil {
		if errors.Is(err, models.ErrNoRows) {
			s.post.Send(&models.Post{Message: "Не найден указанный IP адрес."})
			return nil
		}
		s.post.Send(&models.Post{Message: "При получении адреса произошла ошибка"})
		logger.Error("failed to get address by ip.", logger.ErrAttr(err))
		return err
	}

	address := &models.AddressDTO{
		ID:                data.ID,
		IP:                data.IP,
		Name:              &data.Name,
		MaxRTT:            &data.MaxRTT,
		Count:             &data.Count,
		Timeout:           &data.Timeout,
		PeriodStart:       &data.PeriodStart,
		PeriodEnd:         &data.PeriodEnd,
		Interval:          &data.Interval,
		NotificationCount: &data.NotificationCount,
		Enabled:           &isEnable,
	}

	if err := s.addresses.Update(context.Background(), address); err != nil {
		s.post.Send(&models.Post{Message: "Не удалось обновить адрес."})
		logger.Error("failed to update address.", logger.ErrAttr(err))
		return err
	}
	return nil
}

func (s *MessageService) Delete(message string) error {
	logger.Info("delete ip", logger.StringAttr("message", message))
	parts := strings.Split(message, " ")
	if net.ParseIP(parts[1]) == nil {
		s.post.Send(&models.Post{Message: "Не удалось распознать команду. Неправильный IP адрес."})
		return nil
	}

	if err := s.addresses.Delete(context.Background(), parts[1]); err != nil {
		s.post.Send(&models.Post{Message: "Не удалось удалить адрес."})
		logger.Error("failed to delete address.", logger.ErrAttr(err))
		return err
	}
	return nil
}

func (s *MessageService) decode(message string) *models.AddressDTO {
	address := &models.AddressDTO{}

	parts := strings.Split(message, " ")
	if net.ParseIP(parts[1]) == nil {
		s.post.Send(&models.Post{Message: "Не удалось распознать команду. Неправильный IP адрес."})
		return nil
	}
	address.IP = parts[1]

	parts = parts[2:]

	for i := 0; i < len(parts); i += 1 {
		if (parts[i] == "-n" || parts[i] == "--name") && i+1 < len(parts) {
			address.Name = &parts[i+1]
		}

		if (parts[i] == "-r" || parts[i] == "--rtt") && i+1 < len(parts) {
			rtt, err := time.ParseDuration(parts[i+1])
			if err != nil {
				s.post.Send(&models.Post{Message: "Не удалось распознать команду. Не удалось понять время пинга."})
				return nil
			}
			rtt = rtt * time.Millisecond
			address.MaxRTT = &rtt
		}

		if (parts[i] == "-nc" || parts[i] == "--notification") && i+1 < len(parts) {
			count, err := strconv.Atoi(parts[i+1])
			if err != nil {
				s.post.Send(&models.Post{Message: "Не удалось распознать команду. Не удалось понять количество уведомлений."})
				return nil
			}
			address.NotificationCount = &count
		}

		if (parts[i] == "-p" || parts[i] == "--period") && i+1 < len(parts) {
			dates := strings.Split(parts[i+1], "-")
			if len(dates) != 2 {
				s.post.Send(&models.Post{Message: "Не удалось распознать команду. Не удалось понять период."})
				return nil
			}

			start, err := time.Parse("15:04", dates[0])
			if err != nil {
				s.post.Send(&models.Post{Message: "Не удалось распознать команду. Не удалось понять период."})
				return nil
			}
			end, err := time.Parse("15:04", dates[1])
			if err != nil {
				s.post.Send(&models.Post{Message: "Не удалось распознать команду. Не удалось понять период."})
				return nil
			}
			sd := start.Sub(time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC))
			ed := end.Sub(time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC))

			address.PeriodStart = &sd
			address.PeriodEnd = &ed
		}

		if (parts[i] == "-i" || parts[i] == "--interval") && i+1 < len(parts) {
			interval, err := time.ParseDuration(parts[i+1])
			if err != nil {
				s.post.Send(&models.Post{Message: "Не удалось распознать команду. Не удалось понять интервал."})
				return nil
			}
			interval = interval * time.Millisecond
			address.Interval = &interval
		}
		if (parts[i] == "-t" || parts[i] == "--timeout") && i+1 < len(parts) {
			timeout, err := time.ParseDuration(parts[i+1])
			if err != nil {
				s.post.Send(&models.Post{Message: "Не удалось распознать команду. Не удалось понять таймаут."})
				return nil
			}
			timeout = timeout * time.Millisecond
			address.Timeout = &timeout
		}

		if (parts[i] == "-c" || parts[i] == "--count") && i+1 < len(parts) {
			count, err := strconv.Atoi(parts[i+1])
			if err != nil {
				s.post.Send(&models.Post{Message: "Не удалось распознать команду. Не удалось понять количество пакетов."})
				return nil
			}
			address.Count = &count
		}
	}

	return address
}
