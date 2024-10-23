package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// Client represents a connected user
type Client struct {
	conn    net.Conn
	name    string
	passkey string
}

// handleConnection manages a single client's connection
func handleConnection(conn net.Conn) {
	defer conn.Close()

	client := createClient(conn)
	if client == nil {
		return
	}

	clientJoined(client)
	sendWelcomeMessage(client)
	handleClientMessages(client)
	clientLeft(client)
}

func createClient(conn net.Conn) *Client {
	client := &Client{conn: conn}

	client.name = readInput(conn, "Enter your anonymous name: ")
	if client.name == "" {
		return nil
	}

	client.passkey = readInput(conn, "Enter your secret passkey to join a chat room: ")
	if client.passkey == "" {
		return nil
	}

	return client
}

func readInput(conn net.Conn, prompt string) string {
	fmt.Fprint(conn, prompt)
	input, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return ""
	}
	return strings.TrimSpace(input)
}

func clientJoined(client *Client) {
	clientChan <- client
}

func sendWelcomeMessage(client *Client) {
	fmt.Fprintf(client.conn, "Welcome to the chat room, %s!\n", client.name)
}

func handleClientMessages(client *Client) {
	input := bufio.NewScanner(client.conn)
	for input.Scan() {
		message := input.Text()
		if message == "leave" {
			break
		}

		formattedMessage := fmt.Sprintf("%s: %s", client.name, message)
		sendMessageToChannel(client, formattedMessage)
	}
}

func sendMessageToChannel(client *Client, message string) {
	messageChan <- fmt.Sprintf("%s|%s: %s", client.passkey, client.name, message)
}

func clientLeft(client *Client) {
	leaveMessage := fmt.Sprintf("%s has left the chat", client.name)
	broadcastToChat(leaveMessage, client.passkey)

	// Remove the client from the chatroom
	leaveChan <- client
}
