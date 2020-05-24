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

	isRoomExist, err := room.CheckRoomExist(roomName)
	if err != nil {
		c.RenderError(err)
	}

	if !isRoomExist {
		return c.NotFound("Room does not exist")
	}

	chatroom, err1 := room.GetRoom(roomName)

	if err1 != nil {
		c.RenderError(err1)
	}

	wsProtocol := revel.Config.StringDefault("ws.type", "ws")

	return c.Render(device, chatroom, wsProtocol)
}

// Start this when a room created or user goes to an existing room
func (c Room) RoomWebSocket(device, roomName string, ws revel.ServerWebSocket) revel.Result {
	// Check valid ws
	if ws == nil {
		return nil
	}

	chatroom , err := room.GetRoom(roomName)

	if err != nil {
		c.RenderError(err)
	}

	// If room name exist, allow this device to join
	subscriber := room.Subscribe(device, chatroom)
	//defer room.UnSubscribe(roomName, subscriber)


	// Allow device join to room
	room.Join(device, chatroom)
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
			// otherwise publish message
			room.Message(subscriber, mes, chatroom)
		}
	}

	return nil
}

func receiveFromWs(ws revel.ServerWebSocket, newEvents chan room.Event) {
	var msg room.Event
	for {
		error := ws.MessageReceiveJSON(&msg)
		if error != nil {
			fmt.Println("msg")
			fmt.Println(msg)
			newEvents <- room.Event{
				Type:      "QUIT",
				Timestamp: time.Time{},
			}
			return
		}
		newEvents <- msg
	}
}
