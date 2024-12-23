package main

func (manager *ChatRoomManager) addClientToChatRoom(client *ChatClient) {
	manager.clientMutex.Lock()
	defer manager.clientMutex.Unlock()

	if _, ok := manager.clients[client.passkey]; !ok {
		manager.clients[client.passkey] = []*ChatClient{client}
	} else {
		manager.clients[client.passkey] = append(manager.clients[client.passkey], client)
	}
}

func (manager *ChatRoomManager) removeClientFromChatRoom(client *ChatClient) bool {
	manager.clientMutex.Lock()
	defer manager.clientMutex.Unlock()

	if clientList, ok := manager.clients[client.passkey]; ok {
		for i, c := range clientList {
			// Find the client that is leaving
			if c == client {
				// Remove the client from the chat room slice
				manager.clients[client.passkey] = append(clientList[:i], clientList[i+1:]...)
				// If the room is empty after removal, delete the room entry
				if len(manager.clients[client.passkey]) == 0 {
					delete(manager.clients, client.passkey)
				}
				return true 
			}
		}
	}
	return false 
}
