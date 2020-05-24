package models

import (
	"github.com/jinzhu/gorm"
)

type Event struct {
	gorm.Model
	RoomID uint
	Type      string
	DeviceID    uint
	Device    *Device
	Message   string
}
