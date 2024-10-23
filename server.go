package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/fatih/color"
)

var (
	activeClients  = make(map[string][]*ChatClient) // Map of chat rooms and their clients
	clientChannel  = make(chan *ChatClient)         // Channel for new client connections
	leaveChannel   = make(chan *ChatClient)         // Channel for client disconnections
	messageChannel = make(chan string)              // Channel for incoming messages
	clientMutex    sync.Mutex                       //Mutex to protect access to the clients map

	systemMessageFormatter = color.New(color.FgYellow, color.Bold).SprintFunc()
	userMessageFormatter   = color.New(color.FgCyan).SprintFunc()
)

// handleMessages manages incoming clients, disconnecting clients, and messages
func handleMessageProcessing() {
	for {
		select {
		case chatClient := <-clientChannel:
			handleNewClient(chatClient)
		case chatClient := <-leaveChannel:
			handleClientLeave(chatClient)
		case incomingMessage := <-messageChannel:
			processIncomingMessage(incomingMessage)
		}
	}
}

func handleNewClient(chatClient *ChatClient) {
	fmt.Printf("New client connected: %s\n", chatClient.name)
	addClientToChatRoom(chatClient)
	broadcastToChatRoom(fmt.Sprintf("%s has joined", chatClient.name), chatClient.passkey)
}

func handleClientLeave(client *ChatClient) {
	removeClientFromChatRoom(client)
}

func processIncomingMessage(incomingMessage string) {
	// Split the message string in the format "passkey|name: message"
	//Only the first occurrence of '|' is used for splitting
	parts := strings.SplitN(incomingMessage, "|", 2)
	if len(parts) < 2 {
		return
	}
	// Extract passkey and content from the parts array store
	passkey := parts[0]        //passkey, which identifies the chat room
	messageContent := parts[1] //message to be broadcasted
	broadcastToChatRoom(passkey, messageContent)
}

func addClientToChatRoom(chatClient *ChatClient) {
	clientMutex.Lock()
	defer clientMutex.Unlock()

	if _, ok := activeClients[chatClient.passkey]; !ok {
		activeClients[chatClient.passkey] = []*ChatClient{chatClient}
	} else {
		activeClients[chatClient.passkey] = append(activeClients[chatClient.passkey], chatClient)
	}
}

func removeClientFromChatRoom(chatClient *ChatClient) bool {
	clientMutex.Lock()
	defer clientMutex.Unlock()

	if clientList, ok := activeClients[chatClient.passkey]; ok {
		for i, client := range clientList {
			// Find the client that is leaving
			if client == chatClient {
				// Remove the client from the chat room slice
				activeClients[client.passkey] = append(clientList[:i], clientList[i+1:]...)
				return true
			}
		}
	}
	return false
}

func broadcastMessageToClients(passkey, message string) {
	clientMutex.Lock()
	defer clientMutex.Unlock()

	if clientList, ok := activeClients[passkey]; ok {
		// Send message to each client in the chat room
		for _, client := range clientList {
			fmt.Fprintf(client.connection, "%s\n", userMessageFormatter(message)) // Cyan for user messages
		}
	}
}

func broadcastToChatRoom(message, passkey string) {
	formattedMessage := systemMessageFormatter(fmt.Sprintf("---- %s ----", message)) // Yellow for system messages
	broadcastMessageToClients(passkey, formattedMessage)
}
