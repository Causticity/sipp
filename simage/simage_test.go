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
	// Create a Gray one
	// Test every method
}

func TestGray16SippImage (t *testing.T) {
	// Create a Gray16 one
	// Test every method
}