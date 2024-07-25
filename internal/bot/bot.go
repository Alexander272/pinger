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
	Send(addr string) error
	SendErr(addr, message string) error
}

func (c *MostBotClient) SendErr(addr, message string) error {
	post := &model.Post{
		ChannelId: c.conf.ChannelId,
		Message:   fmt.Sprintf("Пинг по адресу %s не прошел.\n```\n%s\n```", addr, message),
	}

	_, _, err := c.client.CreatePost(post)
	if err != nil {
		return fmt.Errorf("failed to send message. error: %w", err)
	}
	return nil
}

func (c *MostBotClient) Send(addr string) error {
	post := &model.Post{
		ChannelId: c.conf.ChannelId,
		Message:   fmt.Sprintf("Пинг по адресу %s прошел.", addr),
	}

	_, _, err := c.client.CreatePost(post)
	if err != nil {
		return fmt.Errorf("failed to send message. error: %w", err)
	}
	return nil
}
