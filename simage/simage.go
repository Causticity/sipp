package simage

import (
	"image"
	// Package image/png is not used explicitly in the code below,
	// but is imported for its initialization side-effect, which allows
	// image.Decode to understand PNG formatted images.
	"image/png"
	"os"
	"reflect"
	)

type Sippimage interface {
	image.Image
	PixOffset(x, y int) int
	Pix() []uint8
	Val(x, y int) float64
	Bpp() int
	Write(out *string) error
}

type SippGray struct {
	*image.Gray
}

func (sg *SippGray) Pix() []uint8 {
	return sg.Gray.Pix
}

func (sg *SippGray) Val(x, y int) float64 {
	return float64(sg.Gray.Pix[sg.PixOffset(x,y)])
}

func (sg *SippGray) Bpp() int {
	return 8
}

type SippGray16 struct {
	*image.Gray16
}

func (sg16 *SippGray16) Pix() []uint8 {
	return sg16.Gray16.Pix
}

func (sg16 *SippGray16) Val(x, y int) float64 {
	i := sg16.PixOffset(x,y)
	return float64(uint16(sg16.Gray16.Pix[i+0])<<8 | uint16(sg16.Gray16.Pix[i+1]))
}

func (sg *SippGray16) Bpp() int {
	return 16
}

var grayType = reflect.TypeOf(new(image.Gray))
var gray16Type = reflect.TypeOf(new(image.Gray16))

func Read(in *string) (Sippimage, error) {
	reader, err := os.Open(*in)
	if err != nil {
		return nil, err
	}

	defer reader.Close()
	im, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	
	t := reflect.TypeOf(im)
	
	if t == grayType {
		g := new(SippGray)
		g.Gray = im.(*image.Gray)
		return g, err
	}
	
	if t == gray16Type {
		g16 := new(SippGray16)
		g16.Gray16 = im.(*image.Gray16)
		return g16, err
	}
	
	panic("input image must be 8-bit or 16-bit grayscale!")
		
	return nil, err
}

func (img *SippGray) Write(out *string) error {
	return sippWrite(img, out)
}

func (img *SippGray16) Write(out *string) error {
	return sippWrite(img, out)
}

func sippWrite(img image.Image, out *string) error { 
	writer, err := os.Create(*out) 
	if err != nil {
		return err
	}
	return png.Encode(writer, img)
}