package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// Client represents a connected user
type ChatClient struct {
	connection net.Conn
	name       string
	passkey    string
}

// handleConnection manages a single client's connection
func handleClientConnection(connection net.Conn) {
	defer connection.Close()

	chatClient := createClient(connection)
	if chatClient == nil {
		return
	}

	notifyClientJoined(chatClient)
	sendWelcomeMessage(chatClient)
	handleClientMessages(chatClient)
	notifyClientLeft(chatClient)
}

func createClient(connection net.Conn) *ChatClient {
	chatClient := &ChatClient{connection: connection}

	chatClient.name = readInput(connection, "Enter your anonymous name: ")
	if chatClient.name == "" {
		return nil
	}

	chatClient.passkey = readInput(connection, "Enter your secret passkey to join a chat room: ")
	if chatClient.passkey == "" {
		return nil
	}

	return chatClient
}

func readInput(connection net.Conn, prompt string) string {
	fmt.Fprint(connection, prompt)
	input, err := bufio.NewReader(connection).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return ""
	}
	return strings.TrimSpace(input)
}

func notifyClientJoined(client *ChatClient) {
	clientChannel <- client
}

func sendWelcomeMessage(client *ChatClient) {
	fmt.Fprintf(client.connection, "Welcome to the chat room, %s!\n", client.name)
}

func handleClientMessages(client *ChatClient) {
	input := bufio.NewScanner(client.connection)
	for input.Scan() {
		message := input.Text()
		if message == "leave" {
			break
		}

		formattedMessage := fmt.Sprintf("%s: %s", client.name, message)
		sendMessageToChannel(client, formattedMessage)
	}
}

func sendMessageToChannel(chatClient *ChatClient, message string) {
	messageChannel <- fmt.Sprintf("%s|%s: %s", chatClient.passkey, chatClient.name, message)
}

func notifyClientLeft(chatClient *ChatClient) {

	leaveMessage := fmt.Sprintf("%s has left the chat", chatClient.name)
	broadcastToChatRoom(leaveMessage, chatClient.passkey)

	// Remove the client from the chatroom
	leaveChannel <- chatClient
}
