package services

import (
	"strings"

	"github.com/Alexander272/Pinger/internal/models"
)

type InformationService struct {
	post Post
}

func NewInformationService(post Post) *InformationService {
	return &InformationService{
		post: post,
	}
}

type Information interface {
	AboutMe(post *models.Post) error
	Help(post *models.Post) error
}

func (s *InformationService) AboutMe(post *models.Post) error {
	message := "Бот для проверки пинга."
	s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: message})
	return nil
}

func (s *InformationService) Help(post *models.Post) error {
	list := []string{
		"##### Список IP-адресов",
		"`list` или `список`",
		"с параметрами:",
		"```",
		"-a --all - вывести полную информацию о IP-адресах",
		"```",
	}
	add := []string{
		"##### Добавление нового IP-адреса в список",
		"`add <ip>` или `добавить <ip>`",
		"с параметрами:",
		"```",
		"-n, --name - название IP-адреса",
		"-r, --rtt - допустимое время пинга в миллисекундах",
		"-N --notification - количество уведомлений",
		"-p, --period - время в течении которого отправляется запросы (формат: <часы>:<минуты>-<часы>:<минуты>)",
		"-i, --interval - время ожидания между отправкой каждого пакета в миллисекундах",
		"-t, --timeout - задает таймаут до завершения ping в миллисекундах",
		"-c, --count - количество пакетов",
		"```",
		"Пример:",
		"```",
		"добавить 8.8.8.8 -n \"Google\"",
		"add 8.8.8.8 -r 100 -N 3 -p \"10:00-20:25\"",
		"```",
	}
	update := []string{
		"##### Изменение параметров IP-адреса",
		"`update <ip>` или `изменить <ip>` c параметрами аналогичными добавлению",
		"Пример:",
		"```",
		"изменить 8.8.8.8 -n \"Google\"",
		"update 8.8.8.8 -r 100 -N 3 -p \"10:00-20:25\"",
		"```",
	}
	disable := []string{
		"##### Отключение IP-адреса",
		"`disable <ip>` или `отключить <ip>`",
		"Пример:",
		"```",
		"отключить 8.8.8.8",
		"disable 8.8.8.8",
		"```",
	}
	enable := []string{
		"##### Включение IP-адреса",
		"`enable <ip>` или `включить <ip>`",
		"Пример:",
		"```",
		"включить 8.8.8.8",
		"enable 8.8.8.8",
		"```",
	}
	delete := []string{
		"##### Удаление IP-адреса из списка",
		"`delete <ip>` или `удалить <ip>`",
		"Пример:",
		"```",
		"удалить 8.8.8.8",
		"delete 8.8.8.8",
		"```",
	}
	stats := []string{
		"##### Статистика",
		"`stats`, `statistics` или `стат`, `статистика`",
		"По умолчанию бот выведет статистику с 1 числа текущего месяца по последнее число текущего месяца.",
		"Для получения статистики по конкретному IP-адресу пропишите его после команды",
		"с параметрами:",
		"```",
		// "-ip - IP-адрес, статистику которого нужно вывести",
		"-p, --period - диапазон времени за который нужно вывести статистику (формат: <день>[.<месяц>[.<год>]]-<день>[.<месяц>[.<год>]])",
		"```",
		"Пример:",
		"```",
		"стат -p \"01.11-1.12\"",
		"stats 8.8.8.8",
		"```",
	}
	unavailable := []string{
		"##### Список недоступных IP-адресов",
		"`unavailable` или `недоступные`",
		"Выводит список недоступных в данный момент IP-адресов.",
	}
	about := []string{
		"##### Информация о боте",
		"`about` или `информация`",
	}
	// restart := []string{
	// 	"#####Перезапуск бота",
	// 	"`restart` или `перезапустить`",
	// }

	message := []string{
		"### Доступные команды:",
		strings.Join(list, "\n"),
		strings.Join(add, "\n"),
		strings.Join(update, "\n"),
		strings.Join(disable, "\n"),
		strings.Join(enable, "\n"),
		strings.Join(delete, "\n"),
		strings.Join(stats, "\n"),
		strings.Join(unavailable, "\n"),
		strings.Join(about, "\n"),
		// strings.Join(restart, "\n"),
	}

	s.post.Send(&models.Post{ChannelID: post.ChannelID, Message: strings.Join(message, "\n")})
	return nil
}

// func (s *InformationService)
