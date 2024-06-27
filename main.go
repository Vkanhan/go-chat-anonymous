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

var (
	clients     = make(map[string][]*Client) // Map of chat rooms and their clients
	clientChan  = make(chan *Client)         // Channel for new client connections
	leaveChan   = make(chan *Client)         // Channel for client disconnections
	messageChan = make(chan string)          // Channel for incoming messages
)

// handleConnection manages a single client's connection
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Create a new client
	client := &Client{conn: conn}

	//Ask for client's name and read it
	fmt.Fprint(conn, "Enter your name: ")
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
		messageChan <- fmt.Sprintf("%s|%s: %s", client.passkey, client.name, message)
	}

	// Notify that the client has left
	leaveChan <- client
	fmt.Printf("Client disconnected: %s\n", client.name)
}

// handleMessages manages incoming clients, departing clients, and messages
func handleMessages() {
	for {
		select {
		case client := <-clientChan:
			fmt.Printf("New client connected: %s\n", client.name)
			// Check if the chat room (identified by client.passkey) already exists
			if _, ok := clients[client.passkey]; !ok {
				// If the chat room does not exist create a new slice and add the client
				clients[client.passkey] = []*Client{client}
			} else {
				// If the chat room already exists, append the client to the existing slice
				clients[client.passkey] = append(clients[client.passkey], client)
			}
			// Broadcast join message to the chat room
			go broadcastToChat(fmt.Sprintf("%s has joined", client.name), client.passkey)
		case client := <-leaveChan:
			// Check if there are clients associated with the client's passkey
			if clientList, ok := clients[client.passkey]; ok {
				// Remove the client from the chat room
				for i, c := range clientList {
					// Find the client that is leaving
					if c == client {
						// Remove the client from the chat room slice
						clients[client.passkey] = append(clientList[:i], clientList[i+1:]...)
						break
					}
				}
				// Broadcast leave message to the chat room
				go broadcastToChat(fmt.Sprintf("%s has left", client.name), client.passkey)
			}
		case message := <-messageChan:
			// Split the message into passkey and content
			parts := strings.SplitN(message, "|", 2)
			if len(parts) < 2 {
				continue
			}
			// Extract passkey and content from the parts array store
			passkey := parts[0]
			content := parts[1]
			// Broadcast message to the specific chat room
			go broadcastMessage(passkey, content)
		}
	}
}

// broadcastMessage sends a message to all clients in the specified chat room
func broadcastMessage(passkey, message string) {
	if clientList, ok := clients[passkey]; ok {
		// Send message to each client in the chat room
		for _, client := range clientList {
			fmt.Fprintf(client.conn, "%s\n", message)
		}
	}
}

// broadcastToChat is a helper function to broadcast a message to a specific chat room
func broadcastToChat(message, passkey string) {
	broadcastMessage(passkey, message)
}

func main() {
	// Start the server and listen on port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	// Goroutine to handle messages
	go handleMessages()

	// Accept incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}
