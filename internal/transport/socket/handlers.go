package socket

import (
	"regexp"
	"strings"

	"github.com/Alexander272/Pinger/internal/models"
	"github.com/Alexander272/Pinger/internal/services"
	"github.com/Alexander272/Pinger/pkg/error_bot"
	"github.com/Alexander272/Pinger/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/mattermost/mattermost-server/v6/model"
)

type Handler struct {
	socket   *model.WebSocketClient
	user     *model.User
	services *services.Services
}

type Deps struct {
	Socket   *model.WebSocketClient
	User     *model.User
	Services *services.Services
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		socket:   deps.Socket,
		user:     deps.User,
		services: deps.Services,
	}
}

func (h *Handler) Listen() {
	logger.Info("listen socket")
	h.socket.Listen()

	// logger.Debug("socket is listening", logger.AnyAttr("error", h.socket.ListenError))

	for event := range h.socket.EventChannel {
		// Launch new goroutine for handling the actual event.
		// If required, you can limit the number of events beng processed at a time.
		go h.handleEvent(event)
	}
}

func (h *Handler) Close() {
	logger.Info("close socket")
	h.socket.Close()
}

func (h *Handler) handleEvent(event *model.WebSocketEvent) {
	// logger.Debug("event", logger.StringAttr("type", event.EventType()))
	// logger.Debug("event", logger.AnyAttr("data", event))

	if event.EventType() != model.WebsocketEventPosted {
		return
	}
	//TODO можно попробовать по этому параметру отсеивать лишние каналы (чтобы бот не читал сообщения в каналах)
	// channelType := event.GetData()["channel_type"].(string)

	// Since this event is a post, unmarshal it to (*model.Post)
	post := &model.Post{}
	err := json.Unmarshal([]byte(event.GetData()["post"].(string)), &post)
	if err != nil {
		logger.Error("Could not cast event to *model.Post", logger.ErrAttr(err))
		error_bot.Send(&gin.Context{}, err.Error(), post)
		return
	}

	// Ignore messages sent by this bot itself.
	if post.UserId == h.user.Id {
		return
	}
	post.Message = strings.TrimSpace(post.Message)

	matches := [...]struct {
		pattern string
		handler func(*models.Post) error
	}{
		{"^about|^обо мне", h.services.Information.AboutMe},
		{"^list|^список", h.services.Message.List},
		{"^add|^добавить", h.services.Message.Create},
		{"^update|^обновить", h.services.Message.Update},
		{"^disable|^отключить", func(p *models.Post) error { return h.services.Message.ToggleActive(p, false) }},
		{"^enable|^включить", func(p *models.Post) error { return h.services.Message.ToggleActive(p, true) }},
		{"^delete|^удалить", h.services.Message.Delete},
		{"^stats|^statistics|^стат", h.services.Message.Statistics},
		{"help|man|помощь|мануал", h.services.Information.Help},
	}

	for _, match := range matches {
		if ok, _ := regexp.MatchString(match.pattern, post.Message); ok {
			if err := match.handler(&models.Post{ChannelID: post.ChannelId, Message: post.Message}); err != nil {
				error_bot.Send(&gin.Context{}, err.Error(), post)
			}
			return
		}
	}

	if ok, _ := regexp.MatchString("^panic", post.Message); ok {
		panic("panic")
	}

	h.services.Information.Help(&models.Post{ChannelID: post.ChannelId, Message: post.Message})
}
