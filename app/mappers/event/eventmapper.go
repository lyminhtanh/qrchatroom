package eventmapper

import (
	"chatroom/app/mappers"
	"chatroom/app/models"
	"github.com/jinzhu/gorm"
)

type EventMapper struct {
	mappers.Mapper
}

func Insert(event models.Event, db *gorm.DB){
	db.Create(&event)
}
