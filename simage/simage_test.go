package simage

import (
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
	// Read a Gray one
	_, err := Read(filepath.Join(testDir,"barbara.png"))
	if err != nil {
		t.Fatal("Can't read 8-bit test image")
	}
	// Pix
	// Val
	// Bpp
	// Write
	// Thumbnail
}

func TestGray16SippImage (t *testing.T) {
	// Read a Gray16 one
	_, err := Read(filepath.Join(testDir,"cosxcosy.png"))
	if err != nil {
		t.Fatal("Can't read 16-bit test image")
	}
	// Test every method
	// Pix
	// Val
	// Bpp
	// Write
	// Thumbnail
}