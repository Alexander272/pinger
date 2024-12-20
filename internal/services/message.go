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
	"github.com/goodsign/monday"
	"github.com/google/shlex"
)

type MessageService struct {
	addresses Address
	stats     Statistic
	post      Post
}

type MessageDeps struct {
	Address Address
	Stats   Statistic
	Post    Post
}

func NewMessageService(deps *MessageDeps) *MessageService {
	return &MessageService{
		addresses: deps.Address,
		stats:     deps.Stats,
		post:      deps.Post,
	}
}

type Message interface {
	List(post *models.Post) error
	Create(post *models.Post) error
	Update(post *models.Post) error
	Delete(post *models.Post) error
	ToggleActive(post *models.Post, isEnable bool) error
	Statistics(post *models.Post) error
	Unavailable(post *models.Post) error
}

func (s *MessageService) List(post *models.Post) error {
	addresses, err := s.addresses.GetAll(context.Background())
	if err != nil {
		logger.Error("failed to get addresses.", logger.ErrAttr(err))
		s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "#### Ошибка.\nПроизошла ошибка при получении адресов."})
		return err
	}

	isAll := false
	parts := strings.Split(post.Message, " ")
	if len(parts) > 1 && (parts[1] == "--all" || parts[1] == "-a") {
		isAll = true
	}
	table := []string{
		"| № | IP-адрес | Название | Статус |",
		"|:--|:----|:----|:--|",
		// "|:-:|:-:|:-:|",
	}
	if isAll {
		table = []string{
			"| № | IP-адрес | Название | Допустимое время пинга | Количество уведомлений | Период | Интервал отправки пакетов | Таймаут до завершения ping | Количество пакетов | Статус |",
			"|:--|:----|:----|:--|:--|:--|:--|:--|:--|:--|",
		}
	}

	for i, address := range addresses {
		isEnable := "Активен"
		if !address.Enabled {
			isEnable = "Не активен"
		}

		if !isAll {
			table = append(table, fmt.Sprintf("|%d|%s|%s|%s|", i+1, address.IP, address.Name, isEnable))
		} else {
			start := time.Date(0, 1, 1, 0, int(address.PeriodStart.Minutes()), 0, 0, time.UTC)
			end := time.Date(0, 1, 1, 0, int(address.PeriodEnd.Minutes()), 0, 0, time.UTC)
			period := fmt.Sprintf("%s-%s", start.Format("15:04"), end.Format("15:04"))
			table = append(table, fmt.Sprintf("|%d|%s|%s|%d|%d|%s|%d|%d|%d|%s|",
				i+1, address.IP, address.Name, address.MaxRTT.Milliseconds(), address.NotificationCount, period, address.Interval.Milliseconds(),
				address.Timeout.Milliseconds(), address.Count, isEnable,
			))
		}
	}

	s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: strings.Join(table, "\n")})
	return nil
}

func (s *MessageService) Create(post *models.Post) error {
	logger.Info("create ip", logger.StringAttr("message", post.Message))
	address := s.decode(post)
	if address == nil {
		return nil
	}

	if err := s.addresses.Create(context.Background(), address); err != nil {
		if errors.Is(err, models.ErrExist) {
			s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "IP адрес уже добавлен."})
			return nil
		}
		s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "#### Ошибка.\nНе удалось добавить IP адрес."})
		logger.Error("failed to create address.", logger.ErrAttr(err))
		return err
	}

	s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "IP адрес добавлен."})
	return nil
}

func (s *MessageService) Update(post *models.Post) error {
	logger.Info("update ip", logger.StringAttr("message", post.Message))
	address := s.decode(post)
	if address == nil {
		return nil
	}

	data, err := s.addresses.GetByIP(context.Background(), address.IP)
	if err != nil {
		if errors.Is(err, models.ErrNoRows) {
			s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "#### Ошибка.\nНе найден указанный IP адрес."})
			return nil
		}
		s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "#### Ошибка.\nПри получении адреса произошла ошибка"})
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
		s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "#### Ошибка.\nНе удалось обновить IP адрес."})
		logger.Error("failed to update address.", logger.ErrAttr(err))
		return err
	}

	s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "IP адрес обновлен."})
	return nil
}

func (s *MessageService) ToggleActive(post *models.Post, isEnable bool) error {
	logger.Info("toggle active ip", logger.StringAttr("message", post.Message), logger.BoolAttr("isEnable", isEnable))
	parts := strings.Split(post.Message, " ")
	if net.ParseIP(parts[1]) == nil {
		s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "#### Ошибка.\nНе удалось распознать команду. Неправильный IP адрес."})
		return nil
	}

	data, err := s.addresses.GetByIP(context.Background(), parts[1])
	if err != nil {
		if errors.Is(err, models.ErrNoRows) {
			s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "#### Ошибка.\nНе найден указанный IP адрес."})
			return nil
		}
		s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "#### Ошибка.\nПри получении адреса произошла ошибка"})
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
		s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "#### Ошибка.\nНе удалось обновить IP адрес."})
		logger.Error("failed to update address.", logger.ErrAttr(err))
		return err
	}

	s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "IP адрес обновлен."})
	return nil
}

func (s *MessageService) Delete(post *models.Post) error {
	logger.Info("delete ip", logger.StringAttr("message", post.Message))
	parts := strings.Split(post.Message, " ")
	if net.ParseIP(parts[1]) == nil {
		s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "#### Ошибка.\nНе удалось распознать команду. Неправильный IP адрес."})
		return nil
	}

	if err := s.addresses.Delete(context.Background(), parts[1]); err != nil {
		s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "#### Ошибка.\nНе удалось удалить IP адрес."})
		logger.Error("failed to delete address.", logger.ErrAttr(err))
		return err
	}

	s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "IP адрес удален."})
	return nil
}

func (s *MessageService) Statistics(post *models.Post) error {
	logger.Info("statistics ip", logger.StringAttr("message", post.Message))

	parts, err := shlex.Split(post.Message)
	if err != nil {
		s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "#### Ошибка.\nНе удалось распознать команду."})
		logger.Error("failed to split message.", logger.ErrAttr(err))
		return fmt.Errorf("failed to split message. error: %w", err)
	}

	args := []string{"", ""}
	for i := 1; i < len(parts); i++ {
		_, ok := strings.CutPrefix(parts[i], "-p")
		if !ok {
			args[0] = parts[i]
			continue
		}

		if strings.Contains(parts[i], "=") {
			tmp := strings.Split(parts[i], "=")
			args[1] = tmp[1]
			continue
		}

		args[1] = parts[i+1]
		i++
	}

	if args[0] != "" && net.ParseIP(args[0]) == nil {
		s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "#### Ошибка.\nНе удалось распознать команду. Некорректный IP адрес."})
		return nil
	}

	now := time.Now()
	period := &models.GetStatisticDTO{
		PeriodStart: time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()),
		PeriodEnd:   time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location()),
	}
	if args[1] != "" {
		parts := strings.Split(args[1], "-")
		if len(parts) != 2 {
			s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "#### Ошибка.\nНе удалось распознать команду. Некорректный период."})
			return fmt.Errorf("period is not correct")
		}

		start := []int{now.Day(), int(now.Month()), now.Year()}
		startParts := strings.Split(parts[0], ".")
		for i, p := range startParts {
			tmp, err := strconv.Atoi(p)
			if err != nil {
				s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "#### Ошибка.\nНе удалось распознать команду. Некорректный период."})
				return err
			}
			if tmp != 0 {
				start[i] = tmp
			}
		}

		end := []int{now.Day(), int(now.Month()), now.Year()}
		endParts := strings.Split(parts[1], ".")
		for i, p := range endParts {
			end[i], err = strconv.Atoi(p)
			if err != nil {
				s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "#### Ошибка.\nНе удалось распознать команду. Некорректный период."})
				return err
			}
		}
		logger.Debug("stats", logger.AnyAttr("start", start))

		//TODO не работает
		period.PeriodStart = time.Date(start[2], time.Month(start[1]), start[0], 0, 0, 0, 0, now.Location())
		period.PeriodEnd = time.Date(end[2], time.Month(end[1]), end[0], 0, 0, 0, 0, now.Location())
	}

	logger.Debug("stats", logger.AnyAttr("period", period))

	var data []*models.Statistic
	if args[0] != "" {
		data, err = s.stats.GetByIP(context.Background(), &models.GetStatisticByIPDTO{
			IP:          args[0],
			PeriodStart: period.PeriodStart,
			PeriodEnd:   period.PeriodEnd,
		})
	} else {
		data, err = s.stats.Get(context.Background(), period)
	}
	if err != nil {
		s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "#### Ошибка.\nПри получении статистики произошла ошибка"})
		logger.Error("failed to get statistic.", logger.ErrAttr(err))
		return err
	}
	if len(data) == 0 {
		s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "Ничего не найдено"})
		return nil
	}

	table := []string{
		"| № | IP-адрес | Название | Не отвечал |",
		"|:--|:--|:--|:--|",
	}
	isFull := false
	if !data[0].TimeStart.IsZero() {
		table[0] += " С | По |"
		table[1] += ":--|:--|"
		isFull = true
	}

	format := "Mon 2 Jan 2006 15:04:05"
	for i, d := range data {
		hours := ""
		if d.Time.Hours() > 0 {
			hours = fmt.Sprintf("%dч. ", int(d.Time.Hours()))
		}
		row := fmt.Sprintf("|%d|%s|%s|%s|", i+1, d.IP, d.Name,
			fmt.Sprintf("%s%2dм. %02dс.", hours, int(d.Time.Minutes())%60, int(d.Time.Seconds())%60),
		)
		if isFull {
			row += fmt.Sprintf("%s|%s|", monday.Format(d.TimeStart, format, monday.LocaleRuRU), monday.Format(d.TimeEnd, format, monday.LocaleRuRU))
		}
		table = append(table, row)
	}

	s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: strings.Join(table, "\n")})
	return nil
}

func (s *MessageService) Unavailable(post *models.Post) error {
	logger.Info("unavailable ip", logger.StringAttr("message", post.Message))

	data, err := s.stats.GetUnavailable(context.Background(), &models.GetUnavailableDTO{})
	if err != nil {
		s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "#### Ошибка.\nПри получении статистики произошла ошибка"})
		logger.Error("failed to get statistic.", logger.ErrAttr(err))
		return err
	}

	table := []string{
		"| № | IP-адрес | Название | Не отвечает с |",
		"|:--|:--|:--|:--|",
	}
	format := "Mon 2 Jan 2006 15:04:05"
	for i, d := range data {
		row := fmt.Sprintf("|%d|%s|%s|%s|", i+1, d.IP, d.Name, monday.Format(d.TimeStart, format, monday.LocaleRuRU))
		table = append(table, row)
	}

	s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: strings.Join(table, "\n")})
	return nil
}

func (s *MessageService) decode(post *models.Post) *models.AddressDTO {
	address := &models.AddressDTO{}

	// parts := strings.Split(message, " ")
	// pattern := regexp.MustCompile(`\"[^\"]+\"|\S+`)
	// parts := pattern.FindAllString(post.Message, -1)
	parts, err := shlex.Split(post.Message)
	if err != nil {
		s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "#### Ошибка.\nНе удалось распознать команду."})
		logger.Error("failed to split message.", logger.ErrAttr(err))
		return nil
	}

	if net.ParseIP(parts[1]) == nil {
		s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "#### Ошибка.\nНе удалось распознать команду. Некорректный IP адрес."})
		return nil
	}
	address.IP = parts[1]

	parts = parts[2:]
	args := make(map[string]string, len(parts)/2)

	logger.Debug("decode", logger.AnyAttr("args", args))
	// Используя пакет github.com/jessevdk/go-flags я немного сокращу код, но этот пакет содержит сильно больше функций, чем мне нужно
	// плюсом я могу вернуть ошибку пользователю только на английском поэтому пока оставлю так

	for i := 0; i < len(parts); i += 2 {
		if strings.Contains(parts[i], "=") {
			tmp := strings.Split(parts[i], "=")
			args[tmp[0]] = tmp[1]
			i -= 1
			continue
		}
		args[parts[i]] = parts[i+1]
	}

	if name, ok := args["-n"]; ok || args["--name"] != "" {
		address.Name = &name
	}
	if rtt, ok := args["-r"]; ok || args["--rtt"] != "" {
		rttDur, err := time.ParseDuration(rtt + "ms")
		if err != nil {
			s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "#### Ошибка.\nНе удалось распознать команду. Не удалось понять время пинга."})
			return nil
		}
		address.MaxRTT = &rttDur
	}
	if nc, ok := args["-N"]; ok || args["--notification"] != "" {
		count, err := strconv.Atoi(nc)
		if err != nil {
			s.post.Send(&models.Post{
				ChannelID: post.ChannelID,
				Message:   "#### Ошибка.\nНе удалось распознать команду. Не удалось понять количество уведомлений.",
			})
			return nil
		}
		address.NotificationCount = &count
	}
	if period, ok := args["-p"]; ok || args["--period"] != "" {
		dates := strings.Split(period, "-")
		if len(dates) != 2 {
			s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "#### Ошибка.\nНе удалось распознать команду. Не удалось понять период."})
			return nil
		}
		start, err := time.Parse("15:04", dates[0])
		if err != nil {
			s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "#### Ошибка.\nНе удалось распознать команду. Не удалось понять период."})
			return nil
		}
		end, err := time.Parse("15:04", dates[1])
		if err != nil {
			s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "#### Ошибка.\nНе удалось распознать команду. Не удалось понять период."})
			return nil
		}
		sd := start.Sub(time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC))
		ed := end.Sub(time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC))

		address.PeriodStart = &sd
		address.PeriodEnd = &ed
	}
	if interval, ok := args["-i"]; ok || args["--interval"] != "" {
		intervalDur, err := time.ParseDuration(interval + "ms")
		if err != nil {
			s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "#### Ошибка.\nНе удалось распознать команду. Не удалось понять интервал."})
			return nil
		}
		address.Interval = &intervalDur
	}
	if timeout, ok := args["-t"]; ok || args["--timeout"] != "" {
		timeoutDur, err := time.ParseDuration(timeout + "ms")
		if err != nil {
			s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "#### Ошибка.\nНе удалось распознать команду. Не удалось понять таймаут."})
			return nil
		}
		address.Timeout = &timeoutDur
	}
	if count, ok := args["-c"]; ok || args["--count"] != "" {
		countInt, err := strconv.Atoi(count)
		if err != nil {
			s.post.Send(&models.Post{
				ChannelID: post.ChannelID,
				Message:   "#### Ошибка.\nНе удалось распознать команду. Не удалось понять количество пакетов.",
			})
			return nil
		}
		address.Count = &countInt
	}

	return address
}

// func (s *MessageService) decodeNew(post *models.Post) *models.AddressDTO {
// 	address := &models.AddressDTO{}

// 	// parts := strings.Split(message, " ")
// 	pattern := regexp.MustCompile(`\"[^\"]+\"|\S+`)
// 	args := pattern.FindAllString(post.Message, -1)
// 	if net.ParseIP(args[1]) == nil {
// 		s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "#### Ошибка.\nНе удалось распознать команду. Некорректный IP адрес."})
// 		return nil
// 	}
// 	address.IP = args[1]

// 	args = args[2:]

// 	opts := models.Decode{}

// 	_, err := flags.ParseArgs(&opts, args)
// 	if err != nil {

// 		logger.Debug("failed to parse args.", logger.ErrAttr(err))
// 		//TODO надо бы как-то указать пользователя какой флаг неправильно задан
// 		s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "#### Ошибка.\nНе удалось распознать команду."})
// 		return nil
// 	}
// 	// logger.Debug("decode", logger.AnyAttr("opts", opts))

// 	address.Name = opts.Name
// 	address.Count = opts.Count
// 	address.NotificationCount = opts.NotificationCount
// 	// address.Enabled = opts.Enabled
// 	times := [5]time.Duration{}

// 	if opts.MaxRTT != 0 {
// 		times[0] = time.Duration(opts.MaxRTT) * time.Millisecond
// 		address.MaxRTT = &times[0]
// 	}
// 	if opts.Interval != 0 {
// 		times[1] = time.Duration(opts.Interval) * time.Millisecond
// 		address.Interval = &times[1]
// 	}
// 	if opts.Timeout != 0 {
// 		times[2] = time.Duration(opts.Timeout) * time.Millisecond
// 		address.Timeout = &times[2]
// 	}
// 	if opts.Period != "" {
// 		dates := strings.Split(opts.Period, "-")
// 		if len(dates) != 2 {
// 			s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "#### Ошибка.\nНе удалось распознать команду. Не удалось понять период."})
// 			return nil
// 		}
// 		start, err := time.Parse("15:04", dates[0])
// 		if err != nil {
// 			s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "#### Ошибка.\nНе удалось распознать команду. Не удалось понять период."})
// 			return nil
// 		}
// 		end, err := time.Parse("15:04", dates[1])
// 		if err != nil {
// 			s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: "#### Ошибка.\nНе удалось распознать команду. Не удалось понять период."})
// 			return nil
// 		}
// 		times[3] = start.Sub(time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC))
// 		times[4] = end.Sub(time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC))

// 		address.PeriodStart = &times[3]
// 		address.PeriodEnd = &times[4]
// 	}

// 	return address
// }
