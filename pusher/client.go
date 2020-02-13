package pusher

import (
	"fmt"
	"time"

	pusher "github.com/pusher/pusher-http-go"
	log "github.com/sirupsen/logrus"
)

// Status client
type Status int

const (
	statusInvalid Status = iota
	statusHealthy
)

// Client pusher
type Client struct {
	pusher.Client
	Status Status
}

// Trigger ...
func (c *Client) Trigger(channel string, eventName string, data interface{}) error {
	if c.Status != statusHealthy {
		return fmt.Errorf("pusher is not healthy")
	}
	return c.Client.Trigger(channel, eventName, data)
}

// NewClient Wrapper for getting pusher client
func NewClient(c *Config) *Client {
	client := &Client{
		Client: pusher.Client{
			AppID:   c.AppID,
			Key:     c.Key,
			Secret:  c.Secret,
			Cluster: c.Cluster,
			Secure:  c.Secure,
		},
		Status: statusInvalid,
	}
	go func() {
		for {
			err := client.Trigger("ping", "ping", map[string]string{})
			if err != nil {
				log.Errorf("Failed to ping pusher: %v", err)
				client.Status = statusInvalid
			}

			client.Status = statusHealthy
			time.Sleep(time.Minute * 5)
		}
	}()

	return client
}
