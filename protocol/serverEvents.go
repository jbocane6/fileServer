package protocol

import (
	"fmt"
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
		case message := <-manager.broadcast:
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
		length, err := client.socket.Read(message)
		if err != nil {
			manager.unregister <- client
			client.socket.Close()
			break
		}
		if length > 0 {
			fmt.Printf("%v RECEIVED: %v from sender\n", Now(), string(message))
			manager.broadcast <- message
			manager.destiny <- client.channel
		}
	}
}

func (manager *ClientManager) send(client *Client) {
	for {
		select {
		case message, ok := <-client.data:
			if !ok {
				return
			} else if client.channel == <-manager.destiny {
				fmt.Printf("%v SENDING: %v to channel: %d\n", Now(), string(message), client.channel)
				client.socket.Write(message)
			}
		}
	}
}

func StartServerMode() {
	listener := GetServer()
	manager := ClientManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
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
		if c == 1 {
			ch, _ := strconv.Atoi(string(message[:c]))
			connection.Write([]byte("Server accepted connection"))
			client := &Client{socket: connection, channel: ch, data: make(chan []byte)}
			manager.register <- client
			go manager.send(client)
		} else {
			ch, _ := strconv.Atoi(string(message[:c-1]))
			connection.Write([]byte("Server accepted connection"))
			client := &Client{socket: connection, channel: ch, data: make(chan []byte)}
			manager.register <- client
			go manager.receive(client)
		}
	}
}
