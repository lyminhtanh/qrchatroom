package controllers

import (
	. "chatroom/app/constants"
	"chatroom/app/room"
	. "chatroom/app/utils"
	"fmt"
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) GotoRoom(roomName string) revel.Result {
	return c.Redirect("/room/%s", roomName)
}

func (c App) RequestNewRoom() revel.Result {
	roomName := RandString(ROOM_LEN)
	for {
		if !room.CheckRoom(roomName) {
			break
		}
		roomName = RandString(ROOM_LEN)
		fmt.Println(roomName)
	}
	// create new room
	chatroom := room.GetRoom(roomName)
fmt.Println("chatroom info")
fmt.Println(chatroom)
	return c.GotoRoom(chatroom.RoomName)
}