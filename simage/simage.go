// Package simage implements a polymorphic wrapper around Gray and Gray16 images
// from the Go standard library (http://golang.org/pkg/image/), for use by the
// rest of the SIPP library.

package simage

import (
	"image"
	"image/draw"
	// Package image/png is not used explicitly in the code below,
	// but is imported for its initialization side-effect, which allows
	// image.Decode to understand PNG formatted images.
	"image/png"
	"os"
	"reflect"
	)

// Sippimage embeds the Image interface from the Go standard library and adds
// a few methods to enable polymorphism.
type Sippimage interface {
	// Embed the Go image interface to allow reading and writing using the Go 
	// PNG decoder and encoder.
	image.Image
	// Bounded Go images all implement this, though it's not part of the Image
	// interface.
	PixOffset(x, y int) int
	// Pix returns the slice of underlying image data, for efficient access.
	Pix() []uint8
	// Val returns the grayscale value at x, y to a float64.
	Val(x, y int) float64
	// Bpp returns the pixel depth of this image, either 8 or 16.
	Bpp() int
	// Write encodes the image into a PNG of the given name.
	Write(out *string) error
	// Thumbnail returns a 100x100 version of the image. If the original is
	// smaller than 100x100, the returned image will contain the original,
	// centered. Thumbnails are always 8-bit Gray images.
	Thumbnail() Sippimage
}

// A SippGray wraps a Go Gray image and implements the Sippimage interface.
type SippGray struct {
	*image.Gray
}

// Pix returns the slice of underlying image data, for efficient access.
func (sg *SippGray) Pix() []uint8 {
	return sg.Gray.Pix
}

// Val returns the grayscale value at x, y to a float64.
func (sg *SippGray) Val(x, y int) float64 {
	return float64(sg.Gray.Pix[sg.PixOffset(x,y)])
}

// Bpp returns the pixel depth of this image, i.e. 8
func (sg *SippGray) Bpp() int {
	return 8
}

// A SippGray16 wraps a Go Gray16 image and implements the Sippimage interface.
type SippGray16 struct {
	*image.Gray16
}

// Pix returns the slice of underlying image data, for efficient access.
func (sg16 *SippGray16) Pix() []uint8 {
	return sg16.Gray16.Pix
}

// Val returns the grayscale value at x, y to a float64.
func (sg16 *SippGray16) Val(x, y int) float64 {
	i := sg16.PixOffset(x,y)
	return float64(uint16(sg16.Gray16.Pix[i+0])<<8 | uint16(sg16.Gray16.Pix[i+1]))
}

// Bpp returns the pixel depth of this image, i.e. 16
func (sg *SippGray16) Bpp() int {
	return 16
}

var grayType = reflect.TypeOf(new(image.Gray))
var gray16Type = reflect.TypeOf(new(image.Gray16))

// Read decodes the file named by the given string, returning a Sippimage.
// Currently panics if the image is not grayscale, either 8 or 16 bit.
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
	
	// TODO: This should be an error return instead of a panic.
	panic("input image must be 8-bit or 16-bit grayscale!")
		
	return nil, err
}

// Write encodes the image into a PNG of the given name.
func (img *SippGray) Write(out *string) error {
	return sippWrite(img, out)
}

// Write encodes the image into a PNG of the given name.
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

func (img *SippGray) Thumbnail() (Sippimage) {
	return thumbnail(img.Gray)
}

func (img *SippGray16) Thumbnail() (Sippimage) {
	return thumbnail(img.Gray16)
}

func thumbnail(src draw.Image) (Sippimage) {
	thumb := new(SippGray)
	thumb.Gray = image.NewGray(image.Rect(0,0,100,100))
	resize(src, thumb.Gray)
	return thumb
}

func resize(src, dst draw.Image) {
	// HACK: just paint the corner for now
	draw.Draw(dst, dst.Bounds(), src, image.Point{0,0}, draw.Src)
}