package qrcode
import (
	"bytes"
	commonconst "chatroom/app/constants"
	"fmt"
	"github.com/liyue201/goqr"
	"github.com/revel/revel"
	"github.com/skip2/go-qrcode"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"chatroom/app/cloud"
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

//qrcode.RecognizeFile(revel.BasePath + "\\public\\img\\sample.png")
func RecognizeFile(path string) {
	imgdata, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	img, _, err := image.Decode(bytes.NewReader(imgdata))
	if err != nil {
		fmt.Printf("image.Decode error: %v\n", err)
		return
	}
	qrCodes, err := goqr.Recognize(img)
	if err != nil {
		fmt.Printf("Recognize failed: %v\n", err)
		return
	}
	for _, qrCode := range qrCodes {
		fmt.Printf("qrCode text: %s\n", qrCode.Payload)
	}
}

