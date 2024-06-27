## Go chat anonymous

This is a simple chat server implemented in Go, allowing multiple clients to connect to different chat rooms and exchange messages.

Chat with your friends anonymously by typing the same secret key!


### Features:

- **Multiple Chat Rooms:** Clients can join different chat rooms using a passkey.
- **Broadcast Messaging:** Messages sent by one client are broadcasted to all clients in the same chat room.
- **Basic Command Handling:** Clients can send messages and use the command "leave" to disconnect from the server.

### Components:

- **`main.go`:**   Entry point of the server, responsible for listening to incoming connections and managing client connections.
- **`server.go`:** Manages the overall server operation, including handling new connections, client management, and message broadcasting.
- **`client.go`:** Represents a connected client, handles client-specific operations such as message sending and disconnection.
  
### Usage:

1. **Clone the repository:**
    ```sh
   git clone https://github.com/Vkanhan/go-chat-anonymous.git         
   cd go-chat-anonymous

2. **Starting the Server:**
    - Run `go run main.go` to start the server. By default, it listens on port `8080`.

3. **Connecting Clients:**
    - Clients can connect to the server using any TCP client software (e.g., telnet, netcat).
    - Connect to the server's IP address and port (e.g., `telnet localhost 8080`).
    - Upon connection, clients are prompted to enter an anonymous name and a passkey to join a chat room.

4. **Sending Messages:**
   - Clients can send messages to the chat room by typing their message. Messages are formatted as `passkey|name: message`.

5. **Leaving the Chat:**
   - Clients can type `leave` to disconnect from the server and leave the chat room.

### License:
This project is licensed under the  Apache License Version 2.0 [LICENSE](LICENSE). See the LICENSE file for details.

