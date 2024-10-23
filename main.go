package main

import (
	"log"
	"net"
)

const serverPort = "8080"

func main() {
	listener, err := net.Listen("tcp", ":"+serverPort)
	if err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
	defer listener.Close()

	go handleMessageProcessing()

	log.Println("Serve started, waiting for the clients to join..")
	for {
		clientConnection, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go handleClientConnection(clientConnection)
	}
}
