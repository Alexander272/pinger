package services

import (
	"github.com/Alexander272/Pinger/internal/repo"
	"github.com/mattermost/mattermost-server/v6/model"
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
	Client    *model.Client4
	ChannelID string
}

func NewServices(deps *Deps) *Services {
	post := NewPostService(deps.Client, deps.ChannelID)
	addresses := NewAddressService(deps.Repo.Address)
	ping := NewPingService(addresses, post)
	information := NewInformationService(post)
	message := NewMessageService(addresses, post)
	scheduler := NewSchedulerService(ping)

	return &Services{
		Post:        post,
		Address:     addresses,
		Ping:        ping,
		Information: information,
		Message:     message,
		Scheduler:   scheduler,
	}
}
