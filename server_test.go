package main

import (
	"testing"
)

func TestAddClientToChatRoom(t *testing.T) {
	manager := NewChatRoomManager()
	client := &ChatClient{name: "testUser", passkey: "room1"}

	manager.addClientToChatRoom(client)

	if len(manager.clients["room1"]) != 1 {
		t.Errorf("Expected 1 client in room1, got %d", len(manager.clients["room1"]))
	}

	// Check that the client in "room1" is the one we added
	addedClient := manager.clients["room1"][0]
	if addedClient != client {
		t.Errorf("Expected client %v in room1, got %v", client, addedClient)
	}

	if addedClient.name != "testUser" {
		t.Errorf("Expected client name 'testUser', got '%s'", addedClient.name)
	}
}

func TestRemoveClientFromChatRoom(t *testing.T) {
	manager := NewChatRoomManager()
	client := &ChatClient{name: "testUser", passkey: "room1"}

	manager.addClientToChatRoom(client)
	removed := manager.removeClientFromChatRoom(client)

	if !removed {
		t.Errorf("Expected client to be removed, but wasn't")
	}

	if len(manager.clients["room1"]) != 0 {
		t.Errorf("Expected 0 clients in room1, got %d", len(manager.clients["room1"]))
	}
}
