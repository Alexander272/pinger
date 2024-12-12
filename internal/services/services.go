package services

import (
	"github.com/Alexander272/Pinger/internal/repo"
	"github.com/Alexander272/Pinger/pkg/mattermost"
)

type Services struct {
	Post
	Address
	Ping
	Information
	Message
	Scheduler
}

type Deps struct {
	Repo      *repo.Repository
	Client    *mattermost.Client
	ChannelID string
}

func NewServices(deps *Deps) *Services {
	post := NewPostService(deps.Client.Http, deps.ChannelID)
	addresses := NewAddressService(deps.Repo.Address)
	ping := NewPingService(addresses, post)
	information := NewInformationService(post)
	message := NewMessageService(addresses, post)
	scheduler := NewSchedulerService(ping, deps.Client)

	return &Services{
		Post:        post,
		Address:     addresses,
		Ping:        ping,
		Information: information,
		Message:     message,
		Scheduler:   scheduler,
	}
}
