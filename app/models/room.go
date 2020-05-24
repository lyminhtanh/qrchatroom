package models

import (
	"github.com/jinzhu/gorm"
)

type Room struct {
	gorm.Model
	Name    string `gorm:"type:varchar(100);unique_index"`
	Address, QrCodeUrl, QrCodeFilePath     string
	Events []*Event
	Devices []*Device `gorm:"many2many:room_devices;"`
}
