// Copyright Raul Vera 2016

// Tests for package simage.

package simage

import (
//	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"
	)

var testDir = filepath.Join(os.Getenv("GOPATH"), "src", "github.com", 
							"Causticity", "sipp", "testdata")

func TestRead (t *testing.T) {
	// Read a file that doesn't exist
	_, err := Read("blahblah")
	if err == nil {
		t.Error("Error: Read of garbage succeeded!")
	} else {
		t.Log(err) // What did happen? TODO:Can we check that it's the right error?
	}
	
	// Read a file that exists but isn't a png
	_, err = Read(filepath.Join(testDir,"README"))
	if err == nil {
		t.Error("Error: Read of non-image file succeeded!")
	} else {
		t.Log(err)
	}

	// Read a file that isn't gray
	_, err = Read(filepath.Join(testDir,"mandrill.png"))
	if err == nil {
		t.Error("Error: Read of non-gray succeeded!")
	} else {
		t.Log(err)
	}
}

func TestGraySippImage (t *testing.T) {
	barb, err := Read(filepath.Join(testDir,"barbara.png"))
	if err != nil {
		t.Fatal("Fatal: Can't read 8-bit test image")
	}

	gray, ok := barb.(*SippGray)
	if !ok {
		t.Fatal("Fatal: Type of 8-bit image is not Gray")
	}

	if gray.Bpp() != 8 {
		t.Error("Error: 8-bit image reports bpp != 8")
	}

	if &(gray.Pix()[0]) != &(gray.Gray.Pix[0]) {
		t.Error("Error: Pix return incorrect for 8-bit image")
	}
	
	// Val
	b := gray.Bounds()
	x := b.Dx()/2
	y := b.Dy()/2
	val := gray.At(x, y).(color.Gray).Y
	if float64(val) != gray.Val(x, y) {
		t.Errorf("Error: Val failed for gray at pixel %d, %d", x, y)
	}
	
	// Write
		// write the image out
		// compare it to the original (might not be byte for byte, though
		// or
		// write the image out
		// read it back in
		// compare them with a deep copy
		// if success, delete the written copy
		// else save the bad output for debugging and notify
		
	// Thumbnail
}

func TestGray16SippImage (t *testing.T) {
	cc, err := Read(filepath.Join(testDir,"cosxcosy.png"))
	if err != nil {
		t.Fatal("Fatal: Can't read 16-bit test image")
	}

	gray16, ok := cc.(*SippGray16)
	if !ok {
		t.Fatal("Fatal: Type of 16-bit image is not Gray16")
	}

	if gray16.Bpp() != 16 {
		t.Error("Error: 16-bit image reports bpp != 16")
	}

	if &(gray16.Pix()[0]) != &(gray16.Gray16.Pix[0]) {
		t.Error("Error: Pix return incorrect for 16-bit image")
	}

	// Val
	b := gray16.Bounds()
	x := b.Dx()/2
	y := b.Dy()/2
	val := gray16.At(x, y).(color.Gray16).Y
	if float64(val) != gray16.Val(x, y) {
		t.Errorf("Error: Val failed for gray16 at pixel %d, %d", x, y)
	}
	
	// Write
	// Thumbnail
}