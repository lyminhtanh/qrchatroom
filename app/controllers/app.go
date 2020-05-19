package controllers

import (
	"chatroom/app/room"
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) GotoRoom(room *room.ChatRoom) revel.Result {
	return c.Redirect("/room/%s", room.RoomName)
}

func (c App) RequestNewRoom() revel.Result {
	roomName := "newRoom"
	chatroom := room.GetRoom(roomName)
	return c.GotoRoom(chatroom)
}