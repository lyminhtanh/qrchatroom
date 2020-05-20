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
	device = c.Request.RemoteAddr

	isRoomExist := room.CheckRoom(roomName)
	if !isRoomExist {
		return c.NotFound("Room does not exist")
	}

	chatroom := room.GetRoom(roomName)
	return c.Render(device, chatroom)
}

// Start this when a room created or user goes to an existing room
func (c Room) RoomWebSocket(device, roomName string, ws revel.ServerWebSocket) revel.Result {
	// Check valid ws
	if ws == nil {
		return nil
	}

	room.GetRoom(roomName)

	// If room name exist, allow this device to join
	subscriber := room.Subscribe(device, roomName)
	//defer room.UnSubscribe(device, subscriber)

	// Allow device join to room
	room.Join(device, roomName)
	defer room.Leave(device, roomName)

	// Send old Events of this room
	for event := range subscriber.Events {
		if ws.MessageSendJSON(&event) != nil {
			return nil
		}
	}
	// receive json from devices
	newEvents := make(chan string)

	// start new thread to receive messages from devices
	go receiveFromWs(ws, newEvents)

	for {
		select {
		case event := <-subscriber.NewEvent:
			if ws.MessageSendJSON(&event) != nil {
				// error happens
				return nil
			}

		case mes, ok := <-newEvents:
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
	for {
		error := ws.MessageReceiveJSON(&msg)
		if error != nil {
			//close(newEvents)
			return
		}
		newEvents <- msg
	}
}
