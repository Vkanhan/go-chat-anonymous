package main

import (
	"fmt"
	"strings"
)

var (
	clients     = make(map[string][]*Client) // Map of chat rooms and their clients
	clientChan  = make(chan *Client)         // Channel for new client connections
	leaveChan   = make(chan *Client)         // Channel for client disconnections
	messageChan = make(chan string)          // Channel for incoming messages
)

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
			// The message string is expected to be in the format "passkey|name: message"
			// Splitting the string with a limit of 2 ensures that only the first occurrence of '|' is used for splitting
			parts := strings.SplitN(message, "|", 2)
			if len(parts) < 2 {
				continue
			}
			// Extract passkey and content from the parts array store
			passkey := parts[0] // parts[0] contains the passkey, which identifies the chat room
			content := parts[1] // parts[1] contains the rest of the message to be broadcasted
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
