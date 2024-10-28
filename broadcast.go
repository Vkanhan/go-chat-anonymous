package main

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	systemMessage = color.New(color.FgYellow, color.Bold).SprintFunc()
	userMessage   = color.New(color.FgCyan).SprintFunc()
)

func (manager *ChatRoomManager) sendMessageToChannel(client *Client, message string) {
	manager.messageChan <- fmt.Sprintf("%s|%s: %s", client.passkey, client.name, message)
}

func (manager *ChatRoomManager) broadcastMessage(passkey, message string) {
	manager.clientMutex.Lock()
	defer manager.clientMutex.Unlock()

	if clientList, ok := manager.clients[passkey]; ok {
		// Send message to each client in the chat room
		for _, client := range clientList {
			fmt.Fprintf(client.conn, "%s\n", userMessage(message)) // Cyan for user messages
		}
	}
}

func (manager *ChatRoomManager) broadcastToChat(message, passkey string) {
	formattedMessage := systemMessage(fmt.Sprintf("---- %s ----", message)) // Yellow for system messages
	manager.broadcastMessage(passkey, formattedMessage)
}
