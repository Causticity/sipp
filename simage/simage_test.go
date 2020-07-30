// Copyright Raul Vera 2016

// Tests for package simage.

package simage

import (
	"image/color"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

import (
	. "github.com/Causticity/sipp/sipptesting/sipptestcore"
)

// TODO. The coverage tool shows that we aren't testing 3 code paths:
// - an error return from file Create
// - when the src aspect ratio is smaller than the target thumbnail aspect ratio
// - when the first weight when precomputing the scaling filter is less than the
//   minimum fraction.
// These should be corrected at some point, but the first is trivial and the
// other two only apply to thumbnail generation, so this is low priority.
func TestRead(t *testing.T) {
	// Read a file that doesn't exist
	_, err := Read("blahblah")
	if err == nil {
		t.Error("Error: Read of garbage succeeded!")
	} else {
		t.Log(err) // What did happen? Can we check that it's the right error?
	}

	// Read a file that exists but isn't a png
	_, err = Read(filepath.Join(TestDir, "README"))
	if err == nil {
		t.Error("Error: Read of non-image file succeeded!")
	} else {
		t.Log(err)
	}

	// Read a file that isn't gray
	_, err = Read(filepath.Join(TestDir, "mandrill.png"))
	if err == nil {
		t.Error("Error: Read of non-gray succeeded!")
	} else {
		t.Log(err)
	}
}

func TestGraySippImage(t *testing.T) {
	barb, err := Read(filepath.Join(TestDir, "barbara.png"))
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

	// Val and IntVal
	b := gray.Bounds()
	x := b.Dx() / 2
	y := b.Dy() / 2
	val := gray.At(x, y).(color.Gray).Y
	if int32(val) != gray.IntVal(x, y) {
		t.Errorf("Error: IntVal failed for gray at pixel %d, %d", x, y)
	}
	if float64(val) != gray.Val(x, y) {
		t.Errorf("Error: Val failed for gray at pixel %d, %d", x, y)
	}

	// Write
	name := filepath.Join(TestDir, "test.png")
	err = barb.Write(&name)
	if err != nil {
		t.Fatal("Error writing gray: " + err.Error())
	}
	// read it back in
	comp, err := Read(name)
	if err != nil {
		t.Fatal("Error reading written gray")
	}
	if !reflect.DeepEqual(barb, comp) {
		t.Error("Error: written gray and read differ; written saved as " + name)
	} else {
		err = os.Remove(name)
	}

	// Thumbnail
	thm := barb.Thumbnail()
	name = filepath.Join(TestDir, "barb_thumb.png")
	gold, err := Read(name)
	if err != nil {
		t.Fatal("Error reading golden gray thumb")
	}
	if !reflect.DeepEqual(thm, gold) {
		t.Error("Error: golden gray thumbnail and generated differ")
	}
}

func TestGray16SippImage(t *testing.T) {
	cc, err := Read(filepath.Join(TestDir, "cosxcosy.png"))
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
	x := b.Dx() / 2
	y := b.Dy() / 2
	val := gray16.At(x, y).(color.Gray16).Y
	if int32(val) != gray16.IntVal(x, y) {
		t.Errorf("Error: IntVal failed for gray16 at pixel %d, %d", x, y)
	}
	if float64(val) != gray16.Val(x, y) {
		t.Errorf("Error: Val failed for gray16 at pixel %d, %d", x, y)
	}

	// Write
	name := filepath.Join(TestDir, "test16.png")
	err = cc.Write(&name)
	if err != nil {
		t.Fatal("Error writing gray16: " + err.Error())
	}
	// read it back in
	comp, err := Read(name)
	if err != nil {
		t.Fatal("Error reading written gray16")
	}
	if !reflect.DeepEqual(cc, comp) {
		t.Error("Error: written gray16 and read differ; written saved as " + name)
	} else {
		err = os.Remove(name)
	}

	// Thumbnail
	thm := cc.Thumbnail()
	name = filepath.Join(TestDir, "cosxcosy_thumb.png")
	gold, err := Read(name)
	if err != nil {
		t.Fatal("Error reading golden gray16 thumb")
	}
	if !reflect.DeepEqual(thm, gold) {
		t.Error("Error: golden gray16 thumbnail and generated differ")
	}
}
