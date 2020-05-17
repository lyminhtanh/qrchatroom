package controllers

import (
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) GotoRoom(device, roomName string) revel.Result {
	return c.Redirect("/room?device=%s&roomName=%s", device, roomName)
}