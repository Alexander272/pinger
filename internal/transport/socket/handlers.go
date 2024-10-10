package socket

import (
	"github.com/Alexander272/Pinger/internal/services"
	"github.com/Alexander272/Pinger/pkg/logger"
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

	// for event := range h.socket.EventChannel {
	// 	// Launch new goroutine for handling the actual event.
	// 	// If required, you can limit the number of events beng processed at a time.
	// 	go h.handleEvent(event)
	// }
}

func (h *Handler) Close() {
	logger.Info("close socket")
	h.socket.Close()
}
