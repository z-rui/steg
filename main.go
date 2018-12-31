/*
steg is a program for image steganography.
It hides a black-and-white PNG image into
another grayscale one, or does the reverse.

Usage:
1. To mix two pictures:
    steg [-o output] <visible image> <hidden image>
2. To extract the hidden image:
    steg [-o output] <mixed image>
*/
package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

var outputPath = flag.String("o", "", "output file name")

func main() {
	var img image.Image
	flag.Parse()
	switch N := flag.NArg(); N {
	case 1:
		img = demix(openArg(0))
	case 2:
		img = mix(openArg(0), openArg(1))
	default:
		log.Fatalln("Need exactly 1 or 2 arguments, got", N)
	}
	var output *os.File
	if *outputPath != "" {
		var err error
		output, err = os.Create(*outputPath)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		output = os.Stdout
	}
	png.Encode(output, img)
}

func openArg(i int) image.Image {
	path := flag.Arg(i)
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	return img
}

func mix(shown, hidden image.Image) image.Image {
	bounds := shown.Bounds()
	img := image.NewGray(bounds)
	for y := bounds.Min.Y; y <= bounds.Max.Y; y++ {
		for x := bounds.Min.X; x <= bounds.Max.X; x++ {
			px1 := color.GrayModel.Convert(shown.At(x, y)).(color.Gray)
			px2 := color.GrayModel.Convert(hidden.At(x, y)).(color.Gray)
			px1.Y = px1.Y &^ 1
			if px2.Y != 0 {
				px1.Y |= px2.Y & 1
			}
			img.SetGray(x, y, px1)
		}
	}
	return img
}

func demix(mixed image.Image) image.Image {
	bounds := mixed.Bounds()
	img := image.NewGray(bounds)
	for y := bounds.Min.Y; y <= bounds.Max.Y; y++ {
		for x := bounds.Min.X; x <= bounds.Max.X; x++ {
			px := color.GrayModel.Convert(mixed.At(x, y)).(color.Gray)
			px.Y = -(px.Y & 1)
			img.SetGray(x, y, px)
		}
	}
	return img
}
