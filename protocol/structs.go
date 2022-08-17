package protocol

import "net"

type ActiveClient struct {
	channel int
	conn    net.Conn
}
