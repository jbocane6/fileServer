package protocol

import "net"

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "8000"
	SERVER_TYPE = "tcp"
)

type ActiveClient struct {
	channel int
	conn    net.Conn
}
