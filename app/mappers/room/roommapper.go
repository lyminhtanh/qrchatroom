package roommapper

import (
	devicemapper "chatroom/app/mappers/device"
	"chatroom/app/models"
	"fmt"
	"github.com/jinzhu/gorm"
)

func LoadRoomEvents(roomName string, db *gorm.DB) []models.Event {
	var (
		room   models.Room
		events []models.Event
	)
	db.Where("name = ?", roomName).First(&room)

	db.Model(&room).Related(&events)

	return events
}

func Insert(room *models.Room, db *gorm.DB){
	fmt.Println("room.Name")
	fmt.Println(room.Name)
	fmt.Println(room.ID)
	db.Create(room)
}

func SelectByName(roomName string, db *gorm.DB) *models.Room{
	var room models.Room
	db.Where("name = ?", roomName).First(&room)
	return &room
}

func DeleteRoom(room *models.Room, db *gorm.DB) {
	// delete all room events
	ass := db.Model(room).Association("Events").Clear()
	if ass.Error != nil {
		panic(ass.Error)
	}
	ass = db.Model(room).Association("Devices").Clear()
	if ass.Error != nil {
		panic(ass.Error)
	}
	db.Unscoped().Where("room_id IS NULL").Delete(&[]models.Event{})
	db.Unscoped().Delete(room)
	db.Unscoped().Where("devices.id NOT IN (SELECT device_id FROM room_devices)").Delete(&[]models.Device{})
}

func RemoveDeviceFromRoom(room *models.Room, deviceName string, db *gorm.DB) {
	deviceModel := devicemapper.SelectByName(deviceName, db)

	if deviceModel.ID == 0 {
		return
	}

	ass := db.Model(room).Association("Devices").Delete(deviceModel)
	if ass.Error != nil {
		db.Rollback()
		panic(ass.Error)
	}
}
func InsertRoomEvent(room *models.Room, event *models.Event, db *gorm.DB) {
	// create new Event
	db = db.Create(event)

	if db.Error != nil {
		db.Rollback()
		panic(db.Error)
	}

	// link new event to room
	ass := db.Model(room).Association("Events").Append(event)

	if ass.Error != nil {
		db.Rollback()
		panic(ass.Error)
	}

	// link device to room if not linked yet
	if !CheckDeviceInRoom(room, event.Device, db) {
		ass = db.Model(room).Association("Devices").Append(event.Device)
		if ass.Error != nil {
			db.Rollback()
			panic(ass.Error)
		}
	}
}

func CheckDeviceInRoom(room *models.Room, device *models.Device, db *gorm.DB) bool {
	row := db.Table("room_devices").Where("room_id = ? and device_id = ?", room.ID, device.ID).Select("count(device_id) AS count").Row() // (*sql.Row)
	var count uint
	row.Scan(&count)
	return count > 0
}
