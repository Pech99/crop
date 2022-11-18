package main

import (
	"errors"
	"fmt"
	"image"
	"os"
	"strconv"
	"time"

	"image/png"

	"github.com/Pech99/crop/clipboard"
	"github.com/kbinani/screenshot"
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

	clipboard.WriteImmage(img)
	writeImage(img, getName())

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
		argsi[i], err = strconv.Atoi(args[i])
		if err != nil {
			fmt.Println(err)
			return 0, 0, 0, 0, err
		}
	}

	return argsi[0], argsi[1], argsi[2], argsi[3], nil
}

func writeImage(img *image.RGBA, name string) error {
	//"C:\\Users\\vitto\\Desktop\\Debug\\penna\\out\\"
	os.MkdirAll("out\\", 0333)
	fd, err := os.Create("out\\" + name)
	if err != nil {
		return err
	}
	defer fd.Close()

	return png.Encode(fd, img)
}

func getName() string { //2022-Nov-17_161548_634275200.png
	n := time.Now()
	n.Date()
	return fmt.Sprintf("%04d-%s-%02d_%02d%02d%02d_%09d.png", n.Year(), n.Month().String()[:3], n.Day(), n.Hour(), n.Minute(), n.Second(), n.Nanosecond())
}
