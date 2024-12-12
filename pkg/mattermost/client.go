package mattermost

import (
	"log"

	"github.com/mattermost/mattermost-server/v6/model"
)

type Config struct {
	ServerLink string
	Token      string
}

type Client struct {
	Http   *model.Client4
	Socket *model.WebSocketClient
	apiUrl string
	token  string
	// eventChan chan Event
}

// может разделить клиенты и отдельно подключать http и websocket
func NewMattermostClient(conf Config) *Client {
	httpClient := model.NewAPIv4Client("https://" + conf.ServerLink)
	httpClient.SetToken(conf.Token)

	socketClient, err := model.NewWebSocketClient("wss://"+conf.ServerLink, conf.Token)
	if err != nil {
		log.Fatalf("failed to websocket connect to mattermost. error: %s", err.Error())
	}
	// socketClient.Listen()

	return &Client{
		Http:   httpClient,
		Socket: socketClient,
		apiUrl: conf.ServerLink,
		token:  conf.Token,
		// eventChan: make(chan Event, 10),
	}
}

// TODO подумать может действительно вынести подключение по websocket в функцию
func (m *Client) Connect() bool {
	socket, err := model.NewWebSocketClient("wss://"+m.apiUrl, m.token)
	if err != nil {
		log.Printf("[!] Error connecting to the Mattermost WS: %s\n", err.Error())
		return false
	}
	m.Socket = socket
	m.Socket.Listen()
	log.Println("[+] Mattermost Websocket connection established")

	return true
}

func (m *Client) Reconnect() bool {
	log.Println("[!] Reconnecting to the Mattermost WS")
	return m.Connect()
}

func (m *Client) IsConnected() bool {
	if m.Socket.ListenError != nil {
		log.Printf("[!] Error: Lost connect to the Mattermost WS: %s\n", m.Socket.ListenError.Error())
		return false
	}
	return true
}
