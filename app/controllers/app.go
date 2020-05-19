package controllers

import (
	"chatroom/app/qrcode"
	"github.com/revel/log15"
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	a := qrcode.DecodeQrCode("sample.png")
	log15.Debug(a)
	return c.Render()
}

func (c App) GotoRoom(device, roomName string) revel.Result {
	return c.Redirect("/room?device=%s&roomName=%s", device, roomName)
}