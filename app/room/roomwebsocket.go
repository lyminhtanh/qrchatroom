package room

import (
	commonconst "chatroom/app/constants"
	"chatroom/app/qrcode"
	"container/list"
	"fmt"
	"os"
	"time"

	"github.com/revel/revel"
)

// 1. build a map of rooms
var (
	chatrooms = make(map[string]ChatRoom)
)

// 2. create new room and add to map
// 3. Each room has a list of subscibers (devices)
type ChatRoom struct {
	// public
	RoomName, RoomAddr, QrCodeUrl, QrCodeFilePath string

	// private
	events          list.List
	subscribers     *list.List // list of subscriber{chan Event}
	subscribeChan   chan (chan<- Subscription)
	unsubscribeChan chan (<-chan Event)
	messageChan     chan Event
}
type Subscriber struct {
	Events   *list.List
	NewEvent chan<- Event
}

func createRoom(roomName string) *ChatRoom {
	// Add new room to map
	roomAddr := fmt.Sprintf(commonconst.BASE_ROOM_ADDRESS, revel.HTTPAddr, revel.HTTPPort, roomName)
	room := ChatRoom{
		RoomName:       roomName,
		RoomAddr:       roomAddr,
		QrCodeUrl:      qrcode.EncodeUrl(roomAddr, roomName),
		QrCodeFilePath: qrcode.GetQrFilePathByRoomName(roomName),

		subscribers:     list.New(),
		subscribeChan:   make(chan (chan<- Subscription), 10),
		unsubscribeChan: make(chan (<-chan Event), 10),
		messageChan:     make(chan Event),
	}

	chatrooms[room.RoomName] = room

	// Start room as new thread
	go startRoom(&room)
	return &room
}
func CheckRoom(roomName string) bool {
	fmt.Println("CheckRoom chatrooms")
	fmt.Println(chatrooms)
	if _, ok := chatrooms[roomName]; !ok {
		return false
	}
	return true
}
func GetRoom(roomName string) *ChatRoom {

	if room, ok := chatrooms[roomName]; ok {
		return &room
	}
	return createRoom(roomName)
}

func LeaveUnexpected() {
	fmt.Println("LeaveUnexpected")
}

// 7. start a room, loop until all subscribers leaves
func startRoom(room *ChatRoom) {
	defer LeaveUnexpected()
	for {
		select {

		// handle new subscriber
		case subscriptionChan := <-room.subscribeChan:
			// 1. push to subsribers of this room

			// send all events in room into events of current subscriber
			var events []Event
			for event := room.events.Front(); event != nil; event = event.Next() {
				events = append(events, event.Value.(Event))
			}
			subscriber := make(chan Event)

			room.subscribers.PushBack(subscriber)

			subscriptionChan <- Subscription{
				Events:   events,
				NewEvent: subscriber,
			}

		case unsubscribeChan := <-room.unsubscribeChan:

			// 1. remove from subscribers
			for subscriber := room.subscribers.Front(); subscriber != nil; subscriber = subscriber.Next() {
				if subscriber.Value.(chan Event) == unsubscribeChan {
					room.subscribers.Remove(subscriber)
				}
			}

			// Check to close room
			if room.subscribers.Len() == 0 {
				fmt.Println("\t\t\t\troom.endRoom()\n")
				room.endRoom()
				break
			}
		case mes := <-room.messageChan:
			// mes is an event of Join, Leave or Message
			// add to room event
			fmt.Println("mes := <-room.messageChan:")
			fmt.Println(mes)
			room.events.PushBack(mes)

			// send mes to all subscribers, this is also the chan that link to subscription coresponsing device
			for sub := room.subscribers.Front(); sub != nil; sub = sub.Next() {
				sub.Value.(chan Event) <- mes
			}
		}

	}
}

// 7. end a room, when all subscribers leaves
func (room ChatRoom) endRoom() {
	// stop all channels ?
	//close(room.subscribeChan)
	//close(room.unsubscribeChan)
	//for subscriber := room.subscribers.Front(); subscriber != nil; subscriber = subscriber.Next() {
	//		close(subscriber.Value.(chan Event))
	//		room.subscribers.Remove(subscriber)
	//}
	// remove QR image
	err := os.Remove(room.QrCodeFilePath)

	if err != nil {
		fmt.Println(err)
	}

	delete(chatrooms, room.RoomName)

	fmt.Println("endRoom chatrooms")
	fmt.Println(chatrooms)
	//delete QR code
}

// 4. Event : join, leave, message
type Event struct {
	Type      string
	Device    string
	Timestamp time.Time
	Message   string
}

// define Subscription
type Subscription struct {
	Events   []Event
	NewEvent <-chan Event // avoid sending directly through subsription object
	// but allow through subscribers list
}

// 5. Action Join, Leave, Message
// 5.1 once Subscribe, device will receive Subscription {All events (limit ~20), and NewEvent chan Event}
func Subscribe(device string, roomName string) Subscription {
	// Get room
	room, ok := chatrooms[roomName]
	if !ok {
		panic("Subscribe failed, chat room not found")
	}

	subscriber := make(chan Subscription)
	room.subscribeChan <- subscriber
	return <-subscriber // Subsciption will be sent in the room loop at case <-subscribeChan

}

// 5.2 Leave remove device from room's subscribers
func UnSubscribe(roomName string, subscription Subscription) {
	// Get room
	room, ok := chatrooms[roomName]
	if !ok {
		panic("UnSubscribe failed, room not found")
	}
	room.unsubscribeChan <- subscription.NewEvent
}

// 5.3 Message send mes from a user to all subscribers
func Message(subscriber Subscription, mes Event, roomName string) {
	// Get room
	room, ok := chatrooms[roomName]

	if !ok {
		panic("Message failed, room not found")
	}
	if mes.Type == "QUIT" {
		room.unsubscribeChan <- subscriber.NewEvent
		return
	}

	mes.Timestamp = time.Now()
	room.messageChan <- mes
}

func Join(device string, roomName string) {
	// Get room
	room, ok := chatrooms[roomName]

	if !ok {
		panic("Join failed, room not found")
	}

	room.messageChan <- Event{
		Type:      "JOIN",
		Device:    device,
		Timestamp: time.Now(),
		Message:   fmt.Sprintf("%s has joined", device),
	}
}

func Leave(device string, roomName string) {
	// Get room
	room, ok := chatrooms[roomName]

	if !ok {
		panic("Leave failed, room not found")
	}

	room.messageChan <- Event{
		Type:      "LEAVE",
		Device:    device,
		Timestamp: time.Now(),
		Message:   fmt.Sprintf("%s has left", device),
	}
}
