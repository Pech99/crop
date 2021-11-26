package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"os"
	"strconv"

	"image/jpeg"
	"image/png"

	"golang.design/x/clipboard"
)

func main() {

	var n int = 0
	var img, imgold image.Image
	var err error
	imgCh := clipboard.Watch(context.TODO(), clipboard.FmtImage)
	X1, Y1, X2, Y2, err := parsConf(os.Args)
	if err != nil {
		fmt.Println("parsConf:", err)
		return
	}

	for {

		img, _, err = image.Decode(bytes.NewReader(<-imgCh))
		if err != nil {
			fmt.Println("image.Decode:", err)
			return
		}

		if imgold == nil || img.Bounds().Size() != imgold.Bounds().Size() {
			fmt.Print("--> ", n)

			img, err = cropImage(img, image.Rect(X1, Y1, X2, Y2))
			if err != nil {
				fmt.Println("cropImage:", err)
				return
			}
			fmt.Print(" Crp")

			imgB, err := imgToByte(img)
			if err != nil {
				fmt.Println("imgToByte:", err)
				return
			}
			clipboard.Write(clipboard.FmtImage, imgB)
			fmt.Print(" Wrt\n")

			go writeImage(img, fmt.Sprint(n, ".png"))
			n++
			imgold = img
		}

	}
}

func parsConf(args []string) (int, int, int, int, error) {

	if len(os.Args) < 5 {
		return 0, 0, 0, 0, errors.New("Argomenti Insufficenti")
	}

	X1, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println(err)
		return 0, 0, 0, 0, err
	}

	Y1, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println(err)
		return 0, 0, 0, 0, err
	}

	X2, err := strconv.Atoi(args[3])
	if err != nil {
		fmt.Println(err)
		return 0, 0, 0, 0, err
	}

	Y2, err := strconv.Atoi(args[4])
	if err != nil {
		fmt.Println(err)
		return 0, 0, 0, 0, err
	}

	return X1, Y1, X2, Y2, nil
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
