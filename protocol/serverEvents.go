package protocol

import (
	"fmt"
	"io"
	"strconv"
)

func (manager *ClientManager) start() {
	for {
		select {
		case connection := <-manager.register:
			manager.clients[connection] = true
			fmt.Printf("%v Added new connection!, channel: %v\n", Now(), connection.channel)
		case connection := <-manager.unregister:
			if _, ok := manager.clients[connection]; ok {
				close(connection.data)
				fmt.Println("A connection has terminated!")
				delete(manager.clients, connection)
			}
		case message := <-manager.file:
			for connection := range manager.clients {
				select {
				case connection.data <- message:
				default:
					close(connection.data)
					delete(manager.clients, connection)
				}
			}
		}
	}
}

func (manager *ClientManager) receive(client *Client) {
	for {
		message := make([]byte, 4096)
		/*
			//el error es aquí, está leyendo por partes
			length, err := client.socket.Read(message)
			if err != nil {
				manager.unregister <- client
				client.socket.Close()
				break
			}
			if length > 0 {
				fmt.Printf("%v RECEIVED: file from sender\n", Now())
				manager.file <- message
				manager.destiny <- client.channel
			} */
		fl := []byte{}
		for {
			_, err := client.socket.Read(message)
			if err != nil {
				if err != io.EOF {
					client.socket.Close()
					break
				}
				break
			}
			fl = append(fl, message...)
		}
		manager.file <- fl
		manager.destiny <- client.channel
	}
}

func (manager *ClientManager) send(client *Client) {
	defer client.socket.Close()
	for {
		select {
		case message, ok := <-client.data:
			if !ok {
				return
			} else if client.channel == <-manager.destiny {
				fmt.Printf("%v SENDING: file to channel: %d\n", Now(), client.channel)
				client.socket.Write(message)
			}
		}
	}
}

func StartServerMode() {
	listener := GetServer()
	manager := ClientManager{
		clients:    make(map[*Client]bool),
		file:       make(chan []byte),
		destiny:    make(chan int),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	go manager.start()
	for {
		connection, error := listener.Accept()
		if error != nil {
			fmt.Println(error)
		}
		message := make([]byte, 4096)
		c, _ := connection.Read(message)
		var l int
		if c == 1 {
			l = c
		} else {
			l = c - 1
		}
		ch, _ := strconv.Atoi(string(message[:l]))
		connection.Write([]byte("Server accepted connection"))
		client := &Client{socket: connection, channel: ch, data: make(chan []byte)}
		manager.register <- client
		if c == 1 {
			go manager.send(client)
		} else {
			go manager.receive(client)
		}
	}
}
