// Copyright Raul Vera 2020

// Tests for package sgrad.

package sgrad

import (
    "image"
	"reflect"
	"testing"
)

import (
    . "github.com/Causticity/sipp/sipptesting/sipptestcore"
	. "github.com/Causticity/sipp/sipptesting"
)

func TestFromComplex(t *testing.T) {
    grad := FromComplexArray(CosxCosyTinyGrad, 19)
    if !reflect.DeepEqual(grad.Pix, CosxCosyTinyGrad) {
        t.Error("Error: Gradient image array differs from the one constructed from");
    }
    rect := image.Rect(0,0,19,19)
    if !reflect.DeepEqual(grad.Rect, rect) {
        t.Errorf("Error: Gradient image rect incorrect, expected %v, got %v\n",
            rect, grad.Rect)
    }
    if grad.MaxMod != CosxCosyTinyGradMaxMod {
        t.Errorf("Error: Incorrect max modulus. Expected: %f, got %f", 
            CosxCosyTinyGradMaxMod, grad.MaxMod)
    }
}

func TestFdgrad(t *testing.T) {
    grad := Fdgrad(Sgray)
    if grad.Rect.Dx() != Sgray.Rect.Dx()-1  || grad.Rect.Dy() != Sgray.Rect.Dy()-1 {
        t.Errorf("Error: Gradient image rect wrong size: expected [%d,%d], got [%d,%d]",
            Sgray.Rect.Dx()-1, Sgray.Rect.Dy()-1, grad.Rect.Dx(), grad.Rect.Dy())
    }
    if !reflect.DeepEqual(grad.Pix, SmallPicGrad) {
        t.Error("Error: Gradient image incorrect. Expected:" +
            ComplexArrayToString(SmallPicGrad, 3) + "Got:" +
            ComplexArrayToString(grad.Pix, 3))
    }
    if grad.MaxMod != SmallPicGradMaxMod {
        t.Errorf("Error: Incorrect max modulus. Expected: %f, got %f", SmallPicGradMaxMod, grad.MaxMod)
    }

    testGrad := Fdgrad(SgrayCosxCosyTiny)
    if testGrad.Rect.Dx() != SgrayCosxCosyTiny.Rect.Dx()-1  || testGrad.Rect.Dy() != SgrayCosxCosyTiny.Rect.Dy()-1 {
        t.Errorf("Error: Gradient image rect wrong size: expected [%d,%d], got [%d,%d]",
            SgrayCosxCosyTiny.Rect.Dx()-1, SgrayCosxCosyTiny.Rect.Dy()-1, testGrad.Rect.Dx(), testGrad.Rect.Dy())
    }
    if !reflect.DeepEqual(testGrad.Pix, CosxCosyTinyGrad) {
        t.Error("Error: Gradient image incorrect. Expected:" +
            ComplexArrayToString(CosxCosyTinyGrad, 19) + "Got:" +
            ComplexArrayToString(testGrad.Pix, 19))
    }
    if testGrad.MaxMod != CosxCosyTinyGradMaxMod {
        t.Errorf("Error: Incorrect max modulus. Expected: %f, got %f", CosxCosyTinyGradMaxMod, testGrad.MaxMod)
    }
}