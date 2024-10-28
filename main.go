package main

import (
	"log"
	"net"
)

const Serverport = "8080"

func main() {
	listener, err := net.Listen("tcp", ":"+Serverport)
	if err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
	defer listener.Close()

	manager := NewChatRoomManager()

	go manager.handleMessages()

	log.Println("Serve started, waiting for the clients to join..")
	for {
		clientConnection, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go manager.handleConnection(clientConnection)
	}
}
