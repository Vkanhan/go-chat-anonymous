package main

func (manager *ChatRoomManager) addClientToChatRoom(client *Client) {
	manager.clientMutex.Lock()
	defer manager.clientMutex.Unlock()

	if _, ok := manager.clients[client.passkey]; !ok {
		manager.clients[client.passkey] = []*Client{client}
	} else {
		manager.clients[client.passkey] = append(manager.clients[client.passkey], client)
	}
}

func (manager *ChatRoomManager) removeClientFromChatRoom(client *Client) {
	manager.clientMutex.Lock()
	defer manager.clientMutex.Unlock()

	if clientList, ok := manager.clients[client.passkey]; ok {
		for i, c := range clientList {
			// Find the client that is leaving
			if c == client {
				// Remove the client from the chat room slice
				manager.clients[client.passkey] = append(clientList[:i], clientList[i+1:]...)
				break
			}
		}
	}
}
