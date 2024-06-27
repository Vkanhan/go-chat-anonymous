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
	client := createClient(conn)
	if client == nil {
		return
	}

	// Notify that a new client has joined
	clientJoined(client)

	// Welcome message
	sendWelcomeMessage(client)

	// Read and handle messages from the client
	handleClientMessages(client)

	// Notify that the client has left
	clientLeft(client)
}

// createClient initializes a new Client object based on user input
func createClient(conn net.Conn) *Client {
	client := &Client{conn: conn}

	// Ask for client's name and read it
	client.name = readInput(conn, "Enter your anonymous name: ")
	if client.name == "" {
		return nil
	}

	// Ask and read the client's passkey
	client.passkey = readInput(conn, "Enter your secret passkey to join a chat room: ")
	if client.passkey == "" {
		return nil
	}

	return client
}

// readInput prompts the user for input and returns the trimmed result
func readInput(conn net.Conn, prompt string) string {
	fmt.Fprint(conn, prompt)
	input, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return ""
	}
	return strings.TrimSpace(input)
}

// clientJoined notifies the system that a new client has joined
func clientJoined(client *Client) {
	clientChan <- client
}

// sendWelcomeMessage sends a welcome message to the client
func sendWelcomeMessage(client *Client) {
	fmt.Fprintf(client.conn, "Welcome to the chat room, %s!\n", client.name)
}

// handleClientMessages reads and handles messages from the client
func handleClientMessages(client *Client) {
	input := bufio.NewScanner(client.conn)
	for input.Scan() {
		message := input.Text()
		if message == "leave" {
			break
		}
		// Send the message to the message channel
		sendMessageToChannel(client, message)
	}
}

// sendMessageToChannel formats and sends the message to the message channel
func sendMessageToChannel(client *Client, message string) {
	// The message is formatted as "passkey|name: message"
	// The '|' character is used as a delimiter to separate the passkey from the rest of the message
	// This allows the message handler to easily split the message string into the passkey and content parts
	messageChan <- fmt.Sprintf("%s|%s: %s", client.passkey, client.name, message)
}

// clientLeft notifies the system that a client has left
func clientLeft(client *Client) {
	leaveChan <- client
	fmt.Printf("Client disconnected: %s\n", client.name)
}
