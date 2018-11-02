package torrent

import (
	"time"

	"github.com/hekmon/transmissionrpc"
)

type Config struct {
	Username    string
	Password    string
	Host        string
	Port        uint16
	HTTPS       bool
	Timeout     time.Duration
	DownloadDir string
}

type Client struct {
	config Config
	tr     *transmissionrpc.Client
}

type Option func(client *Client)

func WithConfig(config Config) Option {
	return func(client *Client) {
		client.config = config
	}
}

func New(options ...Option) (Client, error) {
	var client = Client{
		config: Config{
			Timeout: 10 * time.Second,
			Host:    "localhost",
			Port:    9091,
		},
	}
	for _, option := range options {
		option(&client)
	}
	var transmission, errNewTransmission = transmissionrpc.New(
		client.config.Host,
		client.config.Username,
		client.config.Username,
		&transmissionrpc.AdvancedConfig{
			Port:        client.config.Port,
			HTTPTimeout: client.config.Timeout,
			HTTPS:       client.config.HTTPS,
		})
	if errNewTransmission != nil {
		return Client{}, errNewTransmission
	}
	client.tr = transmission
	return client, nil
}
