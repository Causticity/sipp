// Copyright Raul Vera 2015-2016

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
	//"fmt"
    "math"
	"os"
	"reflect"
	)

// SippImage embeds the Image interface from the Go standard library and adds
// a few methods to enable polymorphism.
type SippImage interface {
	// Embed the Go image interface to allow reading and writing using the Go 
	// PNG decoder and encoder.
	image.Image
	// Bounded Go images all implement this, though it's not part of the Image
	// interface.
	PixOffset(x, y int) int
	// Pix returns the slice of underlying image data, for efficient access.
	Pix() []uint8
	// Val returns the grayscale value at x, y as a float64.
	Val(x, y int) float64
	// Bpp returns the pixel depth of this image, either 8 or 16.
	Bpp() int
	// Write encodes the image into a PNG of the given name.
	Write(out *string) error
	// Thumbnail returns a small version of the image. If the original is
	// smaller than the thumbnail, the returned image will contain the original,
	// centered. Thumbnails are always 8-bit Gray images.
	Thumbnail() SippImage
}

// A SippGray wraps a Go Gray image and implements the SippImage interface.
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

// A SippGray16 wraps a Go Gray16 image and implements the SippImage interface.
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

// Read decodes the file named by the given string, returning a SippImage.
// Currently panics if the image is not grayscale, either 8 or 16 bit.
func Read(in *string) (SippImage, error) {
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

func (img *SippGray) Thumbnail() (SippImage) {
	return thumbnail(img)
}

func (img *SippGray16) Thumbnail() (SippImage) {
	return thumbnail(img)
}

// Thumbnails are square, this many pixels on a side, padded with black if
// original isn't square.
const thumbSide = 150

func thumbnail(src SippImage) (SippImage) {
	thumb := new(SippGray)
	thumb.Gray = image.NewGray(image.Rect(0,0,thumbSide,thumbSide))
	// TODO: if the original is smaller than or equal to this, just center it
	scaleDown(src, thumb.Gray)
	return thumb
}

// Scale the source image down to the destination image.
// Preserves aspect ratio, leaving unused destination pixels untouched.
// At present it just uses a simple box filter.
// It might be possible to improve performance and clarity by making all
// pixel fractions 1/16 and using essentially fixed-point arithmetic.
func scaleDown(src SippImage, dst *image.Gray) {		
	srcRect := src.Bounds()
	dstRect := dst.Bounds()
	//fmt.Println("srcRect:<", srcRect, ">")
	//fmt.Println("dstRect:<", dstRect, ">")
	
	srcWidth := srcRect.Dx()
	srcHeight := srcRect.Dy()
	dstWidth := dstRect.Dx()
	dstHeight := dstRect.Dy()
	//fmt.Println("srcWidth:<", srcWidth, ">")
	//fmt.Println("srcHeight:<", srcHeight, ">")
	//fmt.Println("dstWidth:<", dstWidth, ">")
	//fmt.Println("dstHeight:<", dstHeight, ">")
		
	srcAR := float64(srcWidth) / float64(srcHeight)
	dstAR := float64(dstWidth) / float64(dstHeight)
	//fmt.Println("srcAR:<", srcAR, ">")
	//fmt.Println("dstAR:<", dstAR, ">")
		
	var scale float64
	var outWidth int
	var outHeight int
	if srcAR < dstAR {
		// scale vertically and use a horizontal offset
		scale = float64(srcHeight) / float64(dstHeight)
		outWidth = int(float64(srcWidth)/scale)
		outHeight = dstHeight
		//fmt.Println("Scaling vertically")
	} else {
		// scale horizontally and use a vertical offset
		scale = float64(srcWidth) / float64(dstWidth)
		outHeight = int(float64(srcHeight)/scale)
		outWidth = dstWidth
		//fmt.Println("Scaling horizontally")
	}
	//fmt.Println("scale:<", scale, ">")
	//fmt.Println("outWidth:<", outWidth, ">")
	//fmt.Println("outHeight:<", outHeight, ">")
	
	// One of the following will be 0.
	hoff := (dstWidth - outWidth) / 2
	voff := (dstHeight - outHeight) / 2
	
		
	// Scale 16-bit images down to 8. We incour the cost spuriously for 8-bit
	// images so that we can access the source polymorphically.
	var scaleBpp float64 = 1.0
	if src.Bpp() == 16 {
		scaleBpp = 1.0 / 256.0
	}
	
	//fmt.Println("scaleBpp:", scaleBpp)
	
	hfilter := preComputeFilter(scale, outWidth, srcWidth, scaleBpp)
	
	intrm := image.NewGray(image.Rect(0,0,outWidth,srcHeight))

	for inty := 0; inty < srcHeight; inty++ {
		// Apply the filter to the source row, generating an intermediate row 
		for intx := 0; intx < outWidth; intx++ {
			var val float64
			for i := 0; i < hfilter[intx].n; i++ {
				val = val + src.Val(hfilter[intx].idx+i, inty) * hfilter[intx].weights[i]
			}
			intrm.Pix[intrm.PixOffset(intx, inty)] = uint8(math.Floor(val+0.5))
		}
	}
	
	vfilter := preComputeFilter(scale, outHeight, srcHeight, 1.0)
	
	for outx := 0; outx < outWidth; outx++ {
		// Apply the filter to the intermediate column, generating an output column
		for outy := 0; outy < outHeight; outy++ {
			var val float64
			for i := 0; i < vfilter[outy].n; i++ {
				index := intrm.PixOffset(outx, vfilter[outy].idx + i)
				val = val + float64(intrm.Pix[index]) * vfilter[outy].weights[i]
			}
			dst.Pix[dst.PixOffset(outx+hoff, outy+voff)] = uint8(math.Floor(val+0.5))
		}
	}
}

// Set of weights to use for an output pixel
type filter struct {
	// Index into the source row/columb where these weights start
    idx int
    // Number of pixels that contribute
	n int
	// Weight for each input pixel. There are n of these.
	weights []float64
}

func preComputeFilter(scale float64, 
					  outSize, srcSize int, 
					  scaleBpp float64) []filter {
	
	ret := make([]filter, outSize)
	
	// The minimum worthwhile fraction of a pixel. This value is also used
	// to avoid direct floating-point comparisons; instead of comparing two
	// values for equality, we test if their difference is smaller than this
	// value.
	const minFrac = 1.0/256.0

	for i:=0;i<outSize;i++ {
		// compute the address and first weight
		addr, invw := math.Modf(float64(i)*scale)
		ret[i].idx = int(addr)
		frstw := 1.0 - invw
		if (frstw < minFrac) {
			ret[i].idx++
			frstw = 0.0
		} else {
			ret[i].n = 1
		}
		// compute the number of pixels
		count, frac := math.Modf(scale - frstw)
		ret[i].n = ret[i].n + int(count)
		if frac >= minFrac {
			ret[i].n++
		} else {
			frac = 0.0
		}
		// allocate the slice of weights
		ret[i].weights = make([]float64,ret[i].n)
		var windx int
		if (frstw >0.0) {
			ret[i].weights[windx] = frstw / scale * scaleBpp
			windx++
		}
		for j:=0; j < int(count); j++ {
			ret[i].weights[windx] = 1.0 / scale * scaleBpp
			windx++
		}
		if frac > 0.0 {
			ret[i].weights[windx] = frac / scale * scaleBpp
		}
	}
	
	return ret
}
