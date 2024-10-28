package main

import (
	"bufio"
	"fmt"
	"strings"
)

func (manager *ChatRoomManager) clientJoined(client *Client) {
	manager.clientChan <- client
}

func (manager *ChatRoomManager) handleClientMessages(client *Client) {
	input := bufio.NewScanner(client.conn)
	for input.Scan() {
		message := input.Text()
		if message == "leave" {
			break
		}

		formattedMessage := fmt.Sprintf("%s: %s", client.name, message)
		manager.sendMessageToChannel(client, formattedMessage)
	}
}

func (manager *ChatRoomManager) clientLeft(client *Client) {
	leaveMessage := fmt.Sprintf("%s has left the chat", client.name)
	manager.broadcastToChat(leaveMessage, client.passkey)

	// Remove the client from the chatroom
	manager.leaveChan <- client
}

func (manager *ChatRoomManager) handleNewClient(client *Client) {
	fmt.Printf("New client connected: %s\n", client.name)
	manager.addClientToChatRoom(client)
	manager.broadcastToChat(fmt.Sprintf("%s has joined", client.name), client.passkey)
}

func (manager *ChatRoomManager) handleClientLeave(client *Client) {
	manager.removeClientFromChatRoom(client)
}

func (manager *ChatRoomManager) handleMessage(message string) {
	// Split the message string in the format "passkey|name: message"
	//Only the first occurrence of '|' is used for splitting
	parts := strings.SplitN(message, "|", 2)
	if len(parts) < 2 {
		return
	}
	passkey := parts[0]
	content := parts[1]
	manager.broadcastMessage(passkey, content)
}