package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/fatih/color"
)

var (
	clients     = make(map[string][]*Client) // Map of chat rooms and their clients
	clientChan  = make(chan *Client)         // Channel for new client connections
	leaveChan   = make(chan *Client)         // Channel for client disconnections
	messageChan = make(chan string)          // Channel for incoming messages
	clientMutex sync.Mutex                   //Mutex to protect access to the clients map

	systemMessage = color.New(color.FgYellow, color.Bold).SprintFunc()
	userMessage   = color.New(color.FgCyan).SprintFunc()
)

// handleMessages manages incoming clients, disconnecting clients, and messages
func handleMessages() {
	for {
		select {
		case client := <-clientChan:
			handleNewClient(client)
		case client := <-leaveChan:
			handleClientLeave(client)
		case message := <-messageChan:
			handleMessage(message)
		}
	}
}

func handleNewClient(client *Client) {
	fmt.Printf("New client connected: %s\n", client.name)
	addClientToChatRoom(client)
	broadcastToChat(fmt.Sprintf("%s has joined", client.name), client.passkey)
}

func handleClientLeave(client *Client) {
	removeClientFromChatRoom(client)
}

func handleMessage(message string) {
	// Split the message string in the format "passkey|name: message"
	//Only the first occurrence of '|' is used for splitting
	parts := strings.SplitN(message, "|", 2)
	if len(parts) < 2 {
		return
	}
	// Extract passkey and content from the parts array store
	passkey := parts[0] //passkey, which identifies the chat room
	content := parts[1] //message to be broadcasted
	broadcastMessage(passkey, content)
}

func addClientToChatRoom(client *Client) {
	clientMutex.Lock()
	defer clientMutex.Unlock()

	if _, ok := clients[client.passkey]; !ok {
		clients[client.passkey] = []*Client{client}
	} else {
		clients[client.passkey] = append(clients[client.passkey], client)
	}
}

func removeClientFromChatRoom(client *Client) bool {
	clientMutex.Lock()
	defer clientMutex.Unlock()

	if clientList, ok := clients[client.passkey]; ok {
		for i, c := range clientList {
			// Find the client that is leaving
			if c == client {
				// Remove the client from the chat room slice
				clients[client.passkey] = append(clientList[:i], clientList[i+1:]...)
				return true
			}
		}
	}
	return false
}

func broadcastMessage(passkey, message string) {
	clientMutex.Lock()
	defer clientMutex.Unlock()

	if clientList, ok := clients[passkey]; ok {
		// Send message to each client in the chat room
		for _, client := range clientList {
			fmt.Fprintf(client.conn, "%s\n", userMessage(message)) // Cyan for user messages
		}
	}
}

func broadcastToChat(message, passkey string) {
	formattedMessage := systemMessage(fmt.Sprintf("---- %s ----", message)) // Yellow for system messages
	broadcastMessage(passkey, formattedMessage)
}
