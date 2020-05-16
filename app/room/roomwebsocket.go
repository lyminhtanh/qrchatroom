package room

import (
	"container/list"
	"fmt"
	"time"
)

// 1. build a map of rooms
var (
	chatrooms = make(map[string]ChatRoom)
)

// 2. create new room and add to map
// 3. Each room has a list of subscibers (devices)
type ChatRoom struct {
	RoomName string
	Subscribers *list.List
	Events *list.List
	subscribeChan chan Event
	unsubscribeChan chan Event
	messageChan chan Event
}



func createRoom(roomName string) ChatRoom {
	// Add new room to map
	room := ChatRoom{
		RoomName:    roomName,
		Subscribers: list.New(),
		Events: list.New(),
		subscribeChan: make(chan Event),
		unsubscribeChan: make(chan Event),
		messageChan: make(chan Event),
	}
	chatrooms[room.RoomName] = room
	fmt.Printf("createRoom %s", room.RoomName)



	fmt.Println(room.Events)
	// Start room as new thread
	go startRoom(&room)

	return room
}

func GetRoom(roomName string) ChatRoom {
	if room, ok := chatrooms[roomName]; ok {
		return room
	}
	fmt.Printf("GetRoom")

	return createRoom(roomName)
}

// 7. start a room, loop until all subscribers leaves
func startRoom(room *ChatRoom){
	fmt.Printf("startRoom")

	defer endRoom(room)
	for {
		select {
			// handle new subscriber
			case subscriber := <-room.subscribeChan:
				// 1. push to subsribers of this room
				fmt.Println("subscriber")
				fmt.Println(subscriber)
				fmt.Println(room.Events)

				room.Subscribers.PushBack(subscriber)
				room.Events.PushBack(subscriber)
			case unsubscribe := <-room.unsubscribeChan:
				// 1. remove from subscribers
				for subscriber := room.Subscribers.Front(); subscriber != nil; subscriber = subscriber.Next() {
					if subscriber.Value == unsubscribe.Device {
						room.Subscribers.Remove(subscriber)
					}
				}
				room.Events.PushBack(unsubscribe)

				// Check to close room
				if room.Subscribers.Len() == 0 {
					break
				}
			case mes := <-room.messageChan:
				room.Events.PushBack(mes)
		}
	}
}

// 7. end a room, when all subscribers leaves
func endRoom(room *ChatRoom){
	// stop all channels ?
	delete(chatrooms, room.RoomName)
	// TODO delete from DB
}

// 4. Event : join, leave, message
type Event struct {
	Type string
	Device string
	Timestamp time.Time
	Message string
}

// 5. Action Join, Leave, Message
// 5.1 Subscribe
func Subscribe(device string, roomName string){
	// Get room
	room, ok := chatrooms[roomName]
	if !ok {
		fmt.Printf("room %s not found", device)
	}
	room.subscribeChan <- Event{
		Type:      "SUBSCRIBE",
		Device:    device,
		Timestamp: time.Now(),
		Message:   fmt.Sprintf("%s has joined ", device),
	}
}

// 5.2 Leave remove device from room's subscribers
func UnSubscribe(device string, roomName string){
	// Get room
	room, ok := chatrooms[roomName]
	if !ok {
		fmt.Printf("room %s not found", device)
	}
	room.unsubscribeChan <- Event{
		Type:      "UNSUBSCRIBE",
		Device:    device,
		Timestamp: time.Now(),
		Message:   fmt.Sprintf("%s has left ", device),
	}
}

// 5.3 Message send mes from a user to all subscribers
func Message(device string, mes string, roomName string){
	// Get room
	room, ok := chatrooms[roomName]
	if !ok {
		fmt.Printf("room %s not found", device)
	}
	room.subscribeChan <- Event{
		Type:      "MESSAGE",
		Device:    device,
		Timestamp: time.Now(),
		Message:   mes,
	}
}



