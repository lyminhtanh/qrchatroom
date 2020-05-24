package devicemapper

import (
	"chatroom/app/models"
	"github.com/jinzhu/gorm"
)

func SelectByName(nickname string, db *gorm.DB) *models.Device{
	var device models.Device
	db.Where("nickname = ?", nickname).First(&device)
	return &device
}