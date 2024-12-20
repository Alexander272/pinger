package services

import (
	"github.com/Alexander272/Pinger/internal/repo"
	"github.com/Alexander272/Pinger/pkg/mattermost"
)

type Services struct {
	Post
	Address
	Statistic
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
	statistic := NewStatisticService(deps.Repo.Statistic)
	ping := NewPingService(&PingDeps{Address: addresses, Stats: statistic, Post: post})
	information := NewInformationService(post)
	message := NewMessageService(&MessageDeps{Address: addresses, Stats: statistic, Post: post})
	scheduler := NewSchedulerService(ping, deps.Client)

	return &Services{
		Post:        post,
		Address:     addresses,
		Statistic:   statistic,
		Ping:        ping,
		Information: information,
		Message:     message,
		Scheduler:   scheduler,
	}
}
