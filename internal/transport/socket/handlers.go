package socket

import (
	"net/http"
	"net/url"
	"regexp"

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
	services services.Services
}

type Deps struct {
	Socket   *model.WebSocketClient
	User     *model.User
	Services services.Services
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
	logger.Debug("event", logger.StringAttr("type", event.EventType()))
	logger.Debug("event", logger.AnyAttr("data", event))

	gCtx := &gin.Context{
		Request: &http.Request{},
	}

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
		gCtx.Request = &http.Request{Method: "Get", URL: &url.URL{Host: "ping-bot", Path: "cast event to post"}}
		error_bot.Send(gCtx, err.Error(), post)
		return
	}

	// Ignore messages sent by this bot itself.
	if post.UserId == h.user.Id {
		return
	}

	match, err := regexp.MatchString("about|обо мне", post.Message)
	if err != nil {
		logger.Error("failed to match string", logger.ErrAttr(err))
	}
	if match {
		h.services.Information.AboutMe()
		return
	}

	match, err = regexp.MatchString("list|список", post.Message)
	if err != nil {
		logger.Error("failed to match string", logger.ErrAttr(err))
	}
	if match {
		if err := h.services.Message.List(post.Message); err != nil {
			gCtx.Request = &http.Request{Method: "Get", URL: &url.URL{Host: "ping-bot", Path: "list ip"}}
			error_bot.Send(gCtx, err.Error(), post)
		}
		return
	}

	match, err = regexp.MatchString("add|добавить", post.Message)
	if err != nil {
		logger.Error("failed to match string", logger.ErrAttr(err))
	}
	if match {
		if err := h.services.Message.Create(post.Message); err != nil {
			gCtx.Request = &http.Request{Method: "Post", URL: &url.URL{Host: "ping-bot", Path: "create ip"}}
			error_bot.Send(gCtx, err.Error(), post)
		}
		return
	}

	match, err = regexp.MatchString("update|обновить", post.Message)
	if err != nil {
		logger.Error("failed to match string", logger.ErrAttr(err))
	}
	if match {
		if err := h.services.Message.Update(post.Message); err != nil {
			gCtx.Request = &http.Request{Method: "Put", URL: &url.URL{Host: "ping-bot", Path: "update ip"}}
			error_bot.Send(gCtx, err.Error(), post)
		}
		return
	}

	match, err = regexp.MatchString("disable|отключить", post.Message)
	if err != nil {
		logger.Error("failed to match string", logger.ErrAttr(err))
	}
	if match {
		if err := h.services.Message.ToggleActive(post.Message, false); err != nil {
			gCtx.Request = &http.Request{Method: "Put", URL: &url.URL{Host: "ping-bot", Path: "disable ip"}}
			error_bot.Send(gCtx, err.Error(), post)
		}
		return
	}
	match, err = regexp.MatchString("enable|включить", post.Message)
	if err != nil {
		logger.Error("failed to match string", logger.ErrAttr(err))
	}
	if match {
		if err := h.services.Message.ToggleActive(post.Message, true); err != nil {
			gCtx.Request = &http.Request{Method: "Put", URL: &url.URL{Host: "ping-bot", Path: "enable ip"}}
			error_bot.Send(gCtx, err.Error(), post)
		}
		return
	}

	match, err = regexp.MatchString("delete|удалить", post.Message)
	if err != nil {
		logger.Error("failed to match string", logger.ErrAttr(err))
	}
	if match {
		if err := h.services.Message.Delete(post.Message); err != nil {
			gCtx.Request = &http.Request{Method: "Delete", URL: &url.URL{Host: "ping-bot", Path: "delete ip"}}
			error_bot.Send(gCtx, err.Error(), post)
		}
		return
	}

	// match, err = regexp.MatchString("help|man|помощь|мануал", post.Message)
	// if err != nil {
	// 	logger.Error("failed to match string", logger.ErrAttr(err))
	// }
	// if match {
	h.services.Information.Help()
	// return
	// }
}
