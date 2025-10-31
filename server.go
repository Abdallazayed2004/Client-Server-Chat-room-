package main

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
)

type Message struct {
	Sender  string
	Content string
}

type ChatService struct {
	mu       sync.Mutex
	messages []string
}

func (c *ChatService) SendMessage(msg Message, reply *[]string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	formatted := fmt.Sprintf("%s: %s", msg.Sender, msg.Content)
	c.messages = append(c.messages, formatted)
	*reply = append([]string(nil), c.messages...)
	fmt.Println(formatted)
	return nil
}

func (c *ChatService) GetMessages(_ string, reply *[]string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	*reply = append([]string(nil), c.messages...)
	return nil
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:1234")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Chat server running on port 1234...")
	rpc.Register(new(ChatService))

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
