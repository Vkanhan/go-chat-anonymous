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

func sendWelcomeMessage(client *Client) {
	fmt.Fprintf(client.conn, "Welcome to the chat room, %s!\n", client.name)
}
