// Copyright Raul Vera 2020

// Tests for package sgrad.

package sgrad

import (
	"reflect"
	"testing"
)

import (
	. "github.com/Causticity/sipp/sipptesting"
)

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
}