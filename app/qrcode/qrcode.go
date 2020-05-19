package qrcode
import (
	"bytes"
	"fmt"
	"github.com/liyue201/goqr"
	"github.com/revel/revel"
	"github.com/skip2/go-qrcode"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
)

func EncodeUrl(url, roomName string) string{
	filePath := fmt.Sprintf("%s\\public\\img\\%s.png", revel.BasePath ,roomName)
	err := qrcode.WriteFile(url, qrcode.Medium, 256, filePath )
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println()
	fileUrl := 	fmt.Sprintf("http://%s:%d/public/img/%s.png",revel.HTTPAddr, revel.HTTPPort, roomName)

	fmt.Println("fileUrl")
	fmt.Println(fileUrl)
	return fileUrl
}


//qrcode.RecognizeFile(revel.BasePath + "\\public\\img\\sample.png")
func RecognizeFile(path string) {
	fmt.Printf("recognize file: %v\n", path)
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

