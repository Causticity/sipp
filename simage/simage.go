// Package simage implements a polymorphic wrapper around Gray and Gray16 images
// from the Go standard library (http://golang.org/pkg/image/), for use by the
// rest of the SIPP library.

package simage

import (
	"image"
	// Package image/png is not used explicitly in the code below,
	// but is imported for its initialization side-effect, which allows
	// image.Decode to understand PNG formatted images.
	"image/png"
	"fmt"
    "math"
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
	return thumbnail(img)
}

func (img *SippGray16) Thumbnail() (Sippimage) {
	return thumbnail(img)
}

func thumbnail(src Sippimage) (Sippimage) {
	thumb := new(SippGray)
	thumb.Gray = image.NewGray(image.Rect(0,0,100,100))
	// TODO: if the original is smaller than or equal to this, just center it
	scaleDown(src, thumb.Gray)
	return thumb
}

// Scale the source image down to the destination image.
// Preserves aspect ratio, leaving unused destination pixels untouched.
// At present it just uses a simple box filter.
// It might be possible to improve performance and clarity by making all
// pixel fractions 1/16 and using essentially fixed-point arithmetic.
func scaleDown(src Sippimage, dst *image.Gray) {		
	srcRect := src.Bounds()
	dstRect := dst.Bounds()
	fmt.Println("srcRect:<", srcRect, ">")
	fmt.Println("dstRect:<", dstRect, ">")
	
	srcWidth := srcRect.Dx()
	srcHeight := srcRect.Dy()
	dstWidth := dstRect.Dx()
	dstHeight := dstRect.Dy()
	fmt.Println("srcWidth:<", srcWidth, ">")
	fmt.Println("srcHeight:<", srcHeight, ">")
	fmt.Println("dstWidth:<", dstWidth, ">")
	fmt.Println("dstHeight:<", dstHeight, ">")
		
	srcAR := float64(srcWidth) / float64(srcHeight)
	dstAR := float64(dstWidth) / float64(dstHeight)
	fmt.Println("srcAR:<", srcAR, ">")
	fmt.Println("dstAR:<", dstAR, ">")
		
	var scale float64
	var outWidth int
	var outHeight int
	if srcAR < dstAR {
		// scale vertically and use a horizontal offset
		scale = float64(srcHeight) / float64(dstHeight)
		outWidth = int(float64(srcWidth)/scale)
		outHeight = dstHeight
		fmt.Println("Scaling vertically")
	} else {
		// scale horizontally and use a vertical offset
		scale = float64(srcWidth) / float64(dstWidth)
		outHeight = int(float64(srcHeight)/scale)
		outWidth = dstWidth
		fmt.Println("Scaling horizontally")
	}
	fmt.Println("scale:<", scale, ">")
	fmt.Println("outWidth:<", outWidth, ">")
	fmt.Println("outHeight:<", outHeight, ">")
	
	// One of the following will be 0.
	hoff := (dstWidth - outWidth) / 2
	voff := (dstHeight - outHeight) / 2
	
		
	hend := hoff + outWidth
	vend := voff + outHeight

	fmt.Println("hoff, hend:<", hoff, ", ", hend, ">")
	fmt.Println("voff, vend:<", voff, ", ", vend, ">")
	
	// Scale 16-bit images down to 8. We incour the cost spuriously for 8-bit
	// images so that we can access the source polymorphically.
	var scaleBpp float64 = 1.0
	if src.Bpp() == 16 {
		scaleBpp = 1.0 / 256.0
	}
	
	fmt.Println("scaleBpp:", scaleBpp)
		
	// The minimum worthwhile fraction of a pixel 
	const minFrac = 1.0/256.0
	
	_, frac := math.Modf(scale)
	fmt.Println("frac:<", frac, ">")
	
	if frac < minFrac {
		// Treat this as an integer scale
		panic("Integer scaling not yet implemented") // TODO: implement
	} else {
		// fractional part of scale is non-trivial
		// We process a left side, 0 or more middle pixels, and a right side
		// similarly for rows a top, a set of middle rows, and a bottom
		// Each of the border areas has a separate weight
		
		// Note that the vertical and horizontal averaging are actually
		// separable, and more efficient done separately, as you avoid doing
		// kernelx X kernely ops per pixel.
				
		var srcy, srcx int
		var weight float64
		
		intrm := image.NewGray(image.Rect(0,0,outWidth,srcHeight))
		
		// generate intermediate rows by applying the scale factor to
		// the source image horizontally
		for inty := 0; inty < srcHeight; inty++ {
			srcx = 0
			weight = 1.0
			srcUsed := 0.0
			for intx := 0; intx < outWidth; intx++ {
				var count float64 = 0.0
				var val float64 = 0.0
				for acc := weight; acc < scale; {
					val = val + src.Val(srcx, srcy) * weight
					srcUsed = srcUsed + weight
					if srcUsed >= 1.0 {
						srcUsed = 0.0
						srcx++
					}
					count = count + 1.0
					acc = acc + weight
					if acc < scale {
						rem := scale - acc
						if rem <1.0 {
							weight = rem
						} else {
							weight = 1.0
						}
					} else {
						// Set up for next set of source pixels
						weight = 1.0 - weight
					}
				}
			val = val/count * scaleBpp
			intrm.Pix[intrm.PixOffset(intx, inty)] = uint8(val)
			}
			srcy++
		}
		
		// generate dst image by applying the scale factor to
		// the intermediate image vertically
		var intrmx, intrmy int
		for outx := hoff; outx < hend; outx++ {
			intrmy = 0
			weight = 1.0
			intrmUsed := 0.0
			for outy := voff; outy < vend; outy++ {
				var count float64 = 0.0
				var val float64 = 0.0
				for acc := weight; acc < scale; {
				pixVal := float64(intrm.Pix[intrm.PixOffset(intrmx, intrmy)])
					val = val + (pixVal * weight)
					intrmUsed = intrmUsed + weight
					if intrmUsed >= 1.0 {
						intrmUsed = 0.0
						intrmy++
					}
					count = count + 1.0
					acc = acc + weight
					if acc < scale {
					rem := scale - acc
						if rem <1.0 {
							weight = rem
						} else {
							weight = 1.0
						}
					} else {
						// Set up for next set of source pixels
						weight = 1.0 - weight
					}
				}
				val = val/count // intermediate image already 8 bit
				dst.Pix[dst.PixOffset(outx, outy)] = uint8(val)
			}
			intrmx++
		}
	}
}
