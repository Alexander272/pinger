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
	AboutMe()
	Help()
}

func (s *InformationService) AboutMe() {
	message := "Бот для проверки пинга."
	s.post.Send(&models.Post{Message: message})
}

func (s *InformationService) Help() {
	list := []string{
		"#####Список IP-адресов",
		"`list` или `список`",
		"с параметрами:",
		"```",
		"-a --all - вывести полную информацию о IP-адресах",
		"```",
	}
	add := []string{
		"#####Добавление нового IP-адреса в список",
		"`add <ip>` или `добавить <ip>`",
		"с параметрами:",
		"```",
		"-n, --name - название IP-адреса",
		"-r, --rtt - допустимое время пинга в миллисекундах",
		"-nc --notification - количество уведомлений",
		"-p, --period - время в течении которого отправляется запросы (формат: <часы>:<минуты>-<часы>:<минуты>)",
		"-i, --interval - время ожидания между отправкой каждого пакета в миллисекундах",
		"-t, --timeout - задает таймаут до завершения ping",
		"-c, --count - количество пакетов",
		"```",
	}
	update := []string{
		"#####Изменение параметров IP-адреса",
		"`update <ip>` или `изменить <ip>` c параметрами аналогичными добавлению",
	}
	disable := []string{
		"#####Отключение IP-адреса",
		"`disable <ip>` или `отключить <ip>`",
	}
	enable := []string{
		"#####Включение IP-адреса",
		"`enable <ip>` или `включить <ip>`",
	}
	delete := []string{
		"#####Удаление IP-адреса из списка",
		"`delete <ip>` или `удалить <ip>`",
	}
	about := []string{
		"#####Информация о боте",
		"`about` или `информация`",
	}
	// restart := []string{
	// 	"#####Перезапуск бота",
	// 	"`restart` или `перезапустить`",
	// }

	message := []string{
		"Доступные команды:",
		strings.Join(list, "\n"),
		strings.Join(add, "\n"),
		strings.Join(update, "\n"),
		strings.Join(disable, "\n"),
		strings.Join(enable, "\n"),
		strings.Join(delete, "\n"),
		strings.Join(about, "\n"),
		// strings.Join(restart, "\n"),
	}

	s.post.Send(&models.Post{Message: strings.Join(message, "\n")})
}

// func (s *InformationService)
