package qrcode

import (
	"chatroom/app/cloud"
	commonconst "chatroom/app/constants"
	"fmt"
	"github.com/revel/revel"
	"github.com/skip2/go-qrcode"
	_ "image/jpeg"
	_ "image/png"
)

func EncodeUrl(url, roomName string) string{
	filePath := fmt.Sprintf(commonconst.BASE_QR_FILE_PATH, revel.BasePath ,roomName)
	err := qrcode.WriteFile(url, qrcode.Medium, 256, filePath )
	if err != nil {
		fmt.Println(err)
	}

	// send to GS
	cloudClient := cloud.Client()

	if err := cloudClient.Write(roomName, filePath); err != nil {
		panic(err)
	}

	fileUrl, err := cloudClient.MakePublic(roomName);

	if err != nil {
		panic(err)
	}
	//fileUrl := 	fmt.Sprintf(commonconst.BASE_QR_FILE_URL, revel.HTTPAddr, revel.HTTPPort, roomName)
	return fileUrl
}

