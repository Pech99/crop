package main

import (
	"fmt"
	"image"
	"os"

	"image/jpeg"
	"image/png"
)

func main() {

	img, err := readImage(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	img, err = cropImage(img, image.Rect(parsConf("cnf.txt")))
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := writeImage(img, os.Args[1]); err != nil {
		fmt.Println(err)
		return
	}
}

func parsConf(name string) (int, int, int, int) {
	/*
		imp, err := os.Open(name)
		if err != nil {
			panic(err)
		}

		var conf []byte = make([]byte, 10)
		var confs string

		for i := len(conf); i == len(conf); {
			if i, err = imp.Read(conf); err == nil {
				confs += string(conf)
			} else {
				panic(err)
			}
			conf = make([]byte, 10)
		}
		imp.Close()

		var dim [5]int
		confs = strings.ReplaceAll(confs, "\r", "\n")
		confs = strings.ReplaceAll(confs, "\n\n", "\n")
		dims := strings.Split(confs, "\n")

		if len(dims) < 4 {
			fmt.Println("parsConf:troppi pochi parametri")
			panic("")
		}

		for i := 0; i < 4; i++ {
			if dim[i], err = strconv.Atoi(dims[i]); err != nil {
				panic("Parametro non Valido")
			}
		}

		return dim[0], dim[1], dim[2], dim[3]
	*/

	return 185, 73, 1230, 720 //1280x720 presentazione da teams quadrata Errore

	//return 104, 74, 1279, 719 //1280x720 presentazione da teams lunga Errore

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

func writeImage(img image.Image, name string) error {
	fd, err := os.Create(name)
	if err != nil {
		return err
	}
	defer fd.Close()

	return jpeg.Encode(fd, img, nil)
}
