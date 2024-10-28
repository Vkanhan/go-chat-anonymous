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

func sendWelcomeMessage(client *ChatClient) {
	fmt.Fprintf(client.connection, "Welcome to the chat room, %s!\n", client.name)
}
