package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"os"
	"strconv"

	"image/png"

	"github.com/kbinani/screenshot"
	"golang.design/x/clipboard"
)

//complile: go build -ldflags "-H windowsgui"
//hiden console

func main() {

	p1x, p1y, p2x, p2y, err := parsConf(os.Args[1:])
	if err != nil {
		fmt.Println("parsConf:", err)
		return
	}

	img, err := screenshot.CaptureRect(image.Rect(p1x, p1y, p2x, p2y))
	if err != nil {
		panic(err)
	}

	imgB, err := imgToByte(img)
	if err != nil {
		fmt.Println("imgToByte:", err)
		return
	}
	clipboard.Write(clipboard.FmtImage, imgB)

}

func parsConf(args []string) (int, int, int, int, error) {

	const nArg int = 4

	if len(args) == 0 || len(args) == 1 && (args[0] == "help" || args[0] == "?") {
		s := "\ninvocare il programma con i seguenti argomenti:\nX1, Y1, X2, Y2\nX1, Y1: coordinate del primo angolo\nX2, Y2: coordinate dell'angolo opposto"
		return 0, 0, 0, 0, errors.New(s)
	}

	if len(args) != nArg {
		return 0, 0, 0, 0, errors.New("argomenti ansufficenti")
	}

	var argsi [nArg]int
	var err error

	for i := 0; i < nArg; i++ {
		argsi[i], err = strconv.Atoi(args[1])
		if err != nil {
			fmt.Println(err)
			return 0, 0, 0, 0, err
		}
	}

	return argsi[0], argsi[1], argsi[2], argsi[3], nil
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

	return reed[:len(reed)-1], nil
}
