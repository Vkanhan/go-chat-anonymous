package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// Client represents a connected user
type Client struct {
	conn    net.Conn // Connection object for the client
	name    string   // Name of the client
	passkey string   // Passkey for joining chat rooms
}

// handleConnection manages a single client's connection
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Create a new client
	client := &Client{conn: conn}

	// Ask for client's name and read it
	fmt.Fprint(conn, "Enter your anonymous name: ")
	name, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("error reading name:", err)
		return
	}
	client.name = strings.TrimSpace(name)

	// Ask and read the client's passkey
	fmt.Fprint(conn, "Enter your secret passkey to join a chat room: ")
	passkey, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("error reading passkey:", err)
		return
	}
	client.passkey = strings.TrimSpace(passkey)

	// Notify that a new client has joined
	clientChan <- client

	// Welcome message
	fmt.Fprintf(conn, "Welcome to the chat room, %s!\n", client.name)

	// Read and handle messages from the client
	input := bufio.NewScanner(conn)
	for input.Scan() {
		message := input.Text()
		if message == "leave" {
			break
		}
		// Send the message to the message channel
		// The message is formatted as "passkey|name: message"
		// The '|' character is used as a delimiter to separate the passkey from the rest of the message
		// This allows the message handler to easily split the message string into the passkey and content parts
		messageChan <- fmt.Sprintf("%s|%s: %s", client.passkey, client.name, message)
	}

	// Notify that the client has left
	leaveChan <- client
	fmt.Printf("Client disconnected: %s\n", client.name)
}
