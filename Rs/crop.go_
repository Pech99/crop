package main

import (
	"fmt"
	"image"
	"log"
	"os"

	"image/png"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	img, err := readImage("pic.png")
	if err != nil {
		return err
	}

	// I've hard-coded a crop rectangle, start (0,0), end (100, 100).
	img, err = cropImage(img, image.Rect(0, 0, 100, 100))
	if err != nil {
		return err
	}

	return writeImage(img, "pic-cropped.png")
}

// readImage reads a image file from disk. We're assuming the file will be png
// format.
func readImage(name string) (image.Image, error) {
	fd, err := os.Open("pic.png")
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	// image.Decode requires that you import the right image package. We've
	// imported "image/png", so Decode will work for png files. If we needed to
	// decode jpeg files then we would need to import "image/jpeg".
	//
	// Ignored return value is image format name.
	img, _, err := image.Decode(fd)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// cropImage takes an image and crops it to the specified rectangle.
func cropImage(img image.Image, crop image.Rectangle) (image.Image, error) {
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}

	// img is an Image interface. This checks if the underlying value has a
	// method called SubImage. If it does, then we can use SubImage to crop the
	// image.
	simg, ok := img.(subImager)
	if !ok {
		return nil, fmt.Errorf("image does not support cropping")
	}

	return simg.SubImage(crop), nil
}

// writeImage writes an Image back to the disk.
func writeImage(img image.Image, name string) error {
	fd, err := os.Create(name)
	if err != nil {
		return err
	}
	defer fd.Close()

	return png.Encode(fd, img)
}
