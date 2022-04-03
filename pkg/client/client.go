package client

import (
	"sync"
	"time"
)

type Client struct {
	mu          sync.RWMutex
	mainAddress string
}

func NewClient(address string) (*Client, error) {
	return &Client{sync.RWMutex{}, address}, nil
}

func (c *Client) Run() error {
	select {
	case <-time.Tick(5 * time.Second):
	}
	return nil
}
