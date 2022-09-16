package utils

import (
	"net"
	"time"
)

// Client struct
type Client struct {
	socket  net.Conn
	channel int
	data    chan string
}

// Manager struct
type ClientManager struct {
	clients    map[*Client]bool
	file       chan string
	destiny    int
	register   chan *Client
	unregister chan *Client
}

// Return actual date
func Now() string {
	return time.Now().Format("02/01/06 - 15:04")
}
