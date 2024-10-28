package main

import (
	"net"
	"sync"
)

type ChatRoomManager struct {
	clients     map[string][]*Client // Map of chat rooms and their clients
	clientChan  chan *Client         // Channel for new client connections
	leaveChan   chan *Client         // Channel for client disconnections
	messageChan chan string          // Channel for incoming messages
	clientMutex sync.Mutex           //Mutex to protect access to the clients map
}

func NewChatRoomManager() *ChatRoomManager {
	return &ChatRoomManager{
		clients:     make(map[string][]*Client),
		clientChan:  make(chan *Client),
		leaveChan:   make(chan *Client),
		messageChan: make(chan string),
	}
}

// handleConnection manages a single client's connection
func (manager *ChatRoomManager) handleConnection(conn net.Conn) {
	defer conn.Close()

	client := createClient(conn)
	if client == nil {
		return
	}

	manager.clientJoined(client)
	sendWelcomeMessage(client)
	manager.handleClientMessages(client)
	manager.clientLeft(client)
}

// handleMessages manages incoming clients, disconnecting clients, and messages
func (manager *ChatRoomManager) handleMessages() {
	for {
		select {
		case client := <-manager.clientChan:
			manager.handleNewClient(client)
		case client := <-manager.leaveChan:
			manager.handleClientLeave(client)
		case message := <-manager.messageChan:
			manager.handleMessage(message)
		}
	}
}
