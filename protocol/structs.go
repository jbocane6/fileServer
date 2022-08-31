package protocol

import "net"

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "8080"
	SERVER_TYPE = "tcp"
)

// The ClientManager structure will hold all of the available clients, received data, and potential incoming or terminating clients.

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	destiny    chan int
	register   chan *Client
	unregister chan *Client
}

// The Client structure will hold information about the socket connection and data to be sent.

type Client struct {
	socket  net.Conn
	channel int
	data    chan []byte
}

type ActiveClient struct {
	Conn net.Conn
}
