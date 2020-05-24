package controllers

import (
	. "chatroom/app/constants"
	"chatroom/app/room"
	. "chatroom/app/utils"
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
	// Generate unique room name
	for {
		isRoomExist, err := room.CheckRoomExist(roomName)
		if err != nil {
			return c.RenderError(err)
		}
		if !isRoomExist {
			break
		}
		roomName = RandString(ROOM_LEN)
	}
	// Create new room
	_, err := room.GetRoom(roomName)
	if err != nil {
		c.RenderError(err)
	}
	return c.GotoRoom(roomName)
}