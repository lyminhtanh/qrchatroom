package qrcode
import (
	"fmt"
	"github.com/revel/log15"
	"image"
	_ "image/jpeg"
	"os"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

func EncodeUrlToQrCode(url string){

}

func DecodeQrCode(fileName string) string{
	file, error := os.Open("")
	if error != nil {
		log15.Debug("ok")
		fmt.Println(error)
	}
	img, _, _ := image.Decode(file)

	// prepare BinaryBitmap
	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)

	// decode image
	qrReader := qrcode.NewQRCodeReader()
	result, _ := qrReader.Decode(bmp, nil)

	fmt.Println(result)
	return result.GetText()
}

func main() {
	// open and decode image file
	file, _ := os.Open("qrcode.jpg")
	img, _, _ := image.Decode(file)

	// prepare BinaryBitmap
	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)

	// decode image
	qrReader := qrcode.NewQRCodeReader()
	result, _ := qrReader.Decode(bmp, nil)

	fmt.Println(result)
}
