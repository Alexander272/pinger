package bot

import (
	"fmt"

	"github.com/Alexander272/Pinger/internal/config"
	"github.com/mattermost/mattermost-server/v6/model"
)

type MostBotClient struct {
	conf   *config.BotConfig
	client *model.Client4
}

func NewMostBotClient(conf *config.BotConfig, client *model.Client4) *MostBotClient {
	return &MostBotClient{
		conf:   conf,
		client: client,
	}
}

type MostBot interface {
	Send(addr *config.Address) error
	SendErr(addr *config.Address, message string) error
	LongPing(message string) error
}

func (c *MostBotClient) SendErr(addr *config.Address, message string) error {
	post := &model.Post{
		ChannelId: c.conf.ChannelId,
		Message:   fmt.Sprintf("Пинг по адресу **%s (%s)** не прошел.\n```\n%s\n```", addr.Ip, addr.Name, message),
	}

	_, _, err := c.client.CreatePost(post)
	if err != nil {
		return fmt.Errorf("failed to send message. error: %w", err)
	}
	return nil
}

func (c *MostBotClient) Send(addr *config.Address) error {
	post := &model.Post{
		ChannelId: c.conf.ChannelId,
		Message:   fmt.Sprintf("Пинг по адресу **%s (%s)** прошел.", addr.Ip, addr.Name),
	}

	_, _, err := c.client.CreatePost(post)
	if err != nil {
		return fmt.Errorf("failed to send message. error: %w", err)
	}
	return nil
}

func (c *MostBotClient) LongPing(message string) error {
	post := &model.Post{
		ChannelId: c.conf.ChannelId,
		Message:   message,
	}

	_, _, err := c.client.CreatePost(post)
	if err != nil {
		return fmt.Errorf("failed to send message. error: %w", err)
	}
	return nil
}
