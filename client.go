package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strings"
)

type Message struct {
	Sender  string
	Content string
}

func main() {
	client, err := rpc.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("Connection error:", err)
	}
	defer client.Close()

	reader := bufio.NewReader(os.Stdin)

	// Ask for user's name
	var name string
	for {
		fmt.Print("Enter your name: ")
		name, _ = reader.ReadString('\n')
		name = strings.TrimSpace(name)
		if name != "" {
			break
		}
		fmt.Println("Name cannot be empty. Please try again.")
	}

	fmt.Printf("\nWelcome, %s! Type your messages below.\n(Type 'exit' to quit)\n\n", name)

	// Chat loop
	for {
		fmt.Print("You: ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		msg := Message{Sender: name, Content: text}
		var history []string

		err := client.Call("ChatService.SendMessage", msg, &history)
		if err != nil {
			fmt.Println("Error sending message:", err)
			break
		}

		fmt.Println("\n📜 Chat History:")
		for _, line := range history {
			fmt.Println(" ", line)
		}
		fmt.Println()
	}
}
