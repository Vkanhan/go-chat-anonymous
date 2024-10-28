package main

import (
	"fmt"
	"net"
)

const Serverport = "8080"

func main() {
	listener, err := net.Listen("tcp", ":"+Serverport)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	manager := NewChatRoomManager()

	go manager.handleMessages()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go manager.handleConnection(conn)
	}
}
