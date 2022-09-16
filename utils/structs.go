package utils

import (
	"net"
	"time"
)

type Client struct {
	socket  net.Conn
	channel int
	data    chan string
}

type ClientManager struct {
	clients    map[*Client]bool
	file       chan string
	destiny    int
	register   chan *Client
	unregister chan *Client
}

func Now() string {
	return time.Now().Format("02/01/06 - 15:04")
}
