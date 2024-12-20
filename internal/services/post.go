package services

import (
	"fmt"

	"github.com/Alexander272/Pinger/internal/models"
	"github.com/Alexander272/Pinger/pkg/error_bot"
	"github.com/gin-gonic/gin"
	"github.com/mattermost/mattermost-server/v6/model"
)

type PostService struct {
	channelID string
	client    *model.Client4
}

func NewPostService(client *model.Client4, channelID string) *PostService {
	return &PostService{
		channelID: channelID,
		client:    client,
	}
}

type Post interface {
	Send(post *models.Post) error
}

func (s *PostService) Send(data *models.Post) error {
	post := &model.Post{
		ChannelId: data.ChannelID,
		Message:   data.Message,
	}
	if data.ChannelID == "" {
		post.ChannelId = s.channelID
	}

	_, _, err := s.client.CreatePost(post)
	if err != nil {
		error_bot.Send(&gin.Context{}, err.Error(), data)
		return fmt.Errorf("failed to send message. error: %w", err)
	}
	return nil
}
