package controllers

import (
	"chatroom/app/room"
	"github.com/revel/revel"
)

type Room struct {
	*revel.Controller
}

// Run this when go to chat room
func (c Room) RoomPage(device, roomName string) revel.Result {
	return c.Render(device, roomName)
}

// Start this when a room created or user goes to an existing room
func (c Room) RoomWebSocket(device, roomName string, ws revel.ServerWebSocket) revel.Result {
	// Check valid ws
	if ws == nil {
		return nil
	}

	// If room name not exist, start a room
	chatroom := room.GetRoom(roomName)
	// If room name exist, allow this device to join
	room.Subscribe(device, roomName)
	defer room.UnSubscribe(device, roomName)

	// Send old Events of this room
	for event := chatroom.Events.Front(); event != nil; event.Next() {
		if ws.MessageSendJSON(&event) != nil {
			return nil
		}
	}

	// receive json from user
	// handle subscribtion TODO
	newEvents := make(chan string)

	// start new thread to receive messages from devices
	go receiveFromWs(ws, newEvents)


	for {
		select {
		case mes, ok := <-newEvents :
			if !ok {
				return nil
			}

			// otherwise publish message
			room.Message(device, mes, roomName)
		}
	}

	return nil
}

func receiveFromWs(ws revel.ServerWebSocket, newEvents chan string) {
	var msg string
	for{
		error := ws.MessageReceiveJSON(&msg)
		if error != nil {
			close(newEvents)
			return
		}
		newEvents <- msg
	}
}
