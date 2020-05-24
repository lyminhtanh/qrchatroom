package app

import (
	"chatroom/app/db"
	roommapper "chatroom/app/mappers/room"
	"chatroom/app/models"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
	"os"
)

var (
	// AppVersion revel app version (ldflags)
	AppVersion string

	// BuildTime revel app build-time (ldflags)
	BuildTime string
)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.BeforeAfterFilter,       // Call the before and after filter functions
		revel.ActionInvoker,           // Invoke the action.
	}

	// Register startup functions with OnAppStart
	// revel.DevMode and revel.RunMode only work inside of OnAppStart. See Example Startup Script
	// ( order dependent )
	revel.OnAppStart(InitDbScript)
	//revel.OnAppStart(SeedDbScript)
	revel.OnAppStart(ExampleStartupScript)
	// revel.OnAppStart(InitDB)
	// revel.OnAppStart(FillCache)
}

// HeaderFilter adds common security headers
// There is a full implementation of a CSRF filter in
// https://github.com/revel/modules/tree/master/csrf
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")
	c.Response.Out.Header().Add("Referrer-Policy", "strict-origin-when-cross-origin")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

func SeedDbScript() {
	db, err := db.Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// new room
	device := models.Device{
		Model:       gorm.Model{},
		Nickname:    "iphone",
		FullAddress: "1.1.1.1",
	}

	event := models.Event{
		Model:     gorm.Model{},
		Type:      "JOIN",
		Device:  &device,
		Message:   "has joined",
	}

	room := models.Room{
		Model:          gorm.Model{},
		Name:           "AAAAA",
		Address:        "",
		QrCodeUrl:      "",
		QrCodeFilePath: "",
		Events:         []*models.Event{&event},
		Devices:        []*models.Device{&device},
	}

	roommapper.Insert(&room, db)
	//
	//fmt.Println("room.Events")
	//fmt.Println(room.Events)
	//
	// add event to new room
	event = models.Event{
		Model:     gorm.Model{},
		Type:      "MESSAGE",
		Device:  &device,
		Message:   "has joined",
	}

	roommapper.InsertRoomEvent(&room, &event, db)
	//
	device = models.Device{
		Model:       gorm.Model{},
		Nickname:    "android",
		FullAddress: "1.1.1.2",
	}

	event = models.Event{
		Model:     gorm.Model{},
		Type:      "JOIN",
		Device:  &device,
		Message:   "has joined",
	}

	roommapper.InsertRoomEvent(&room, &event, db)
	event.ID = 0
	roommapper.InsertRoomEvent(&room, &event, db)




}
func InitDbScript() {
	conn, err := db.Connect()
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	conn.DropTableIfExists(models.Room{})
	conn.DropTableIfExists(models.Event{})
	conn.DropTableIfExists(models.Device{})
	conn.DropTableIfExists("room_devices")

	//if !conn.HasTable(&models.Room{}) {
	//	conn.CreateTable(models.Room{})
	//}
	//if !conn.HasTable(&models.Device{}) {
	//	conn.CreateTable(models.Device{})
	//}
	//if !conn.HasTable(&models.Event{}) {
	//	conn.CreateTable(models.Event{})
	//}

	conn.AutoMigrate(models.Room{})
	conn.AutoMigrate(models.Event{})
	conn.AutoMigrate(models.Device{})
}

func ExampleStartupScript() {
	// revel.DevMod and revel.RunMode work here
	// Use this script to check for dev mode and set dev/prod startup scripts here!
	if revel.DevMode {
		// Dev mode
		keypath, kok := revel.Config.String("google.keypath")
		if !kok {
			fmt.Println("--- keypath not found")
			return
		}
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", keypath)
	}

}
