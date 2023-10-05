package mattermost

import (
	"github.com/mattermost/mattermost-server/v6/model"
)

type Config struct {
	Server string
	Token  string
}

func NewMattermostClient(conf Config) *model.Client4 {
	client := model.NewAPIv4Client(conf.Server)

	client.SetToken(conf.Token)

	return client
}
