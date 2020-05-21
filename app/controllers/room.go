package controllers

import (
	"chatroom/app/room"
	"fmt"
	"github.com/revel/revel"
	"time"
)

type Room struct {
	*revel.Controller
}

// Run this when go to chat room
func (c Room) RoomPage(device, roomName string) revel.Result {

	device = c.Request.Header.Get("X-Forwarded-For")

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
	//defer room.UnSubscribe(roomName, subscriber)


	// Allow device join to room
	room.Join(device, roomName)
	//defer room.Leave(device, roomName)

	for _, event := range subscriber.Events {
		if ws.MessageSendJSON(&event) != nil {
			// error happens
			return nil
		}
	}
	// receive json from devices
	newEvents := make(chan room.Event)

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
			fmt.Println("<-newEvents")
			fmt.Println(mes)
			// otherwise publish message
			room.Message(subscriber, mes, roomName)
		}
	}

	return nil
}

func receiveFromWs(ws revel.ServerWebSocket, newEvents chan room.Event) {
	var msg room.Event
	for {
		error := ws.MessageReceiveJSON(&msg)
		//println(error)
		if error != nil {
			//close(newEvents)

			newEvents <- room.Event{
				Type:      "QUIT",
				Device:    "",
				Timestamp: time.Time{},
				Message:   "",
			}
			return
		}
		newEvents <- msg
	}
}
