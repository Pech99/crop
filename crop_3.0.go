package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"os"
	"strconv"

	"image/jpeg"
	"image/png"

	"github.com/kbinani/screenshot"
	"golang.design/x/clipboard"
)

//complile: go build -ldflags "-H windowsgui"
//hiden console

func main() {

	disp, X1, Y1, X2, Y2, err := parsConf(os.Args)
	if err != nil {
		fmt.Println("parsConf:", err)
		return
	}

	bounds := screenshot.GetDisplayBounds(disp)

	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		panic(err)
	}

	img2, err := cropImage(img, image.Rect(X1, Y1, X2, Y2))
	if err != nil {
		fmt.Println("cropImage:", err)
		return
	}
	fmt.Print(" Crp")

	imgB, err := imgToByte(img2)
	if err != nil {
		fmt.Println("imgToByte:", err)
		return
	}
	clipboard.Write(clipboard.FmtImage, imgB)
	fmt.Print(" Wrt\n")

}

func parsConf(args []string) (int, int, int, int, int, error) {

	if len(args) == 1 || len(args) == 2 && (args[1] == "help" || args[1] == "?") {
		s := "\ninvocare il programma con i seguenti argomenti:\ndisp, X1, Y1, X2, Y2\ndisp: numero del display di cui efetuare lo screen\nX1, Y1: coordinate del primo angolo\nX2, Y2: coordinate dell'angolo opposto"
		return 0, 0, 0, 0, 0, errors.New(s)
	}

	if len(args) != 6 {
		return 0, 0, 0, 0, 0, errors.New("Argomenti Insufficenti")
	}

	disp, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println(err)
		return 0, 0, 0, 0, 0, err
	}

	X1, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println(err)
		return 0, 0, 0, 0, 0, err
	}

	Y1, err := strconv.Atoi(args[3])
	if err != nil {
		fmt.Println(err)
		return 0, 0, 0, 0, 0, err
	}

	X2, err := strconv.Atoi(args[4])
	if err != nil {
		fmt.Println(err)
		return 0, 0, 0, 0, 0, err
	}

	Y2, err := strconv.Atoi(args[5])
	if err != nil {
		fmt.Println(err)
		return 0, 0, 0, 0, 0, err
	}

	return disp, X1, Y1, X2, Y2, nil
}

func readImage(name string) (image.Image, error) {
	fd, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	img, err := png.Decode(fd)
	//img, _, err := image.Decode(fd)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func cropImage(img image.Image, crop image.Rectangle) (image.Image, error) {
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}

	simg, ok := img.(subImager)
	if !ok {
		return nil, fmt.Errorf("image does not support cropping")
	}

	return simg.SubImage(crop), nil
}

func imgToByte(img image.Image) ([]byte, error) {

	buff := new(bytes.Buffer)
	err := png.Encode(buff, img)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(buff.Bytes())
	var reed []byte = make([]byte, reader.Len())
	_, err = reader.Read(reed)
	if err != nil {
		return nil, err
	}

	return reed, nil
}

func writeImage(img image.Image, name string) error {
	os.MkdirAll("out\\", 0333)
	fd, err := os.Create("out\\" + name)
	if err != nil {
		return err
	}
	defer fd.Close()

	return jpeg.Encode(fd, img, nil)
}
