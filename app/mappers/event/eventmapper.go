package eventmapper

import (
	"chatroom/app/models"
	"github.com/jinzhu/gorm"
)

func Insert(event *models.Event, db *gorm.DB){
	db.Create(event)
}
