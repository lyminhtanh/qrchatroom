package models

import (
	"github.com/jinzhu/gorm"
)

type Device struct {
	gorm.Model
	Nickname string `gorm:"type:varchar(100);unique_index"`
	FullAddress string `gorm:"type:varchar(100);unique_index"`
}
