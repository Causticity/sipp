// Copyright Raul Vera 2020

// Tests for package sgrad.

package sgrad

import (
    //"fmt"
    "image"
    "math"
	"reflect"
	"testing"
)

import (
    . "github.com/Causticity/sipp/sipptesting/sipptestcore"
	. "github.com/Causticity/sipp/sipptesting"
)

var smallPicGrad = []complex128 {
    5 - 3i, 5 - 3i, 5 - 3i,
    5 - 3i, 5 - 3i, 5 - 3i,
    5 - 3i, 5 - 3i, 5 - 3i,
}

var smallPicGradMaxMod = math.Sqrt(34.0)

var identityKernel = SippGradKernel {
    {1 + 0i, 0 + 0i},
    {0 + 0i, 0 + 0i},
}

var imagIdentityKernel = SippGradKernel {
    {0 + 1i, 0 + 0i},
    {0 + 0i, 0 + 0i},
}
var kieransKernel = SippGradKernel {
    {-1 + 0i, 0 - 1i},
    {0 + 1i, 1 + 0i},
}

var smallPicGradKieransKernel = []complex128 {
    5 + 3i, 5 + 3i, 5 + 3i,
    5 + 3i, 5 + 3i, 5 + 3i,
    5 + 3i, 5 + 3i, 5 + 3i,
}

// TODO Make this table-driven.

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
        t.Errorf("Error: Incorrect max modulus. Expected: %v, got %v", 
            CosxCosyTinyGradMaxMod, grad.MaxMod)
    }
}

func TestFdgrad(t *testing.T) {
    grad := Fdgrad(Sgray)
    if grad.Rect.Dx() != Sgray.Rect.Dx()-1  || grad.Rect.Dy() != Sgray.Rect.Dy()-1 {
        t.Errorf("Error: Gradient image rect wrong size: expected [%d,%d], got [%d,%d]",
            Sgray.Rect.Dx()-1, Sgray.Rect.Dy()-1, grad.Rect.Dx(), grad.Rect.Dy())
    }
    if !reflect.DeepEqual(grad.Pix, smallPicGrad) {
        t.Error("Error: Gradient image incorrect. Expected:" +
            ComplexArrayToString(smallPicGrad, 3) + "Got:" +
            ComplexArrayToString(grad.Pix, 3))
    }
    if grad.MaxMod != smallPicGradMaxMod {
        t.Errorf("Error: Incorrect max modulus. Expected: %f, got %f", smallPicGradMaxMod, grad.MaxMod)
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
    idGrad := FdgradKernel(SgrayCosxCosyTiny, identityKernel)
	for y := 0; y < idGrad.Rect.Dy(); y++ {
		for x := 0; x < idGrad.Rect.Dx(); x++ { 
		    if real(idGrad.Pix[y*idGrad.Rect.Dy()+x]) != SgrayCosxCosyTiny.Val(x, y) {
		        t.Errorf("Error: Identity gradient incorrect at %v, %v. Expected %v, got %v\n",
		                 x, y, SgrayCosxCosyTiny.Val(x, y), real(idGrad.Pix[y*idGrad.Rect.Dy()+x]))
		    }
		}
    }
    imagIdGrad := FdgradKernel(SgrayCosxCosyTiny, imagIdentityKernel)
	for y := 0; y < imagIdGrad.Rect.Dy(); y++ {
		for x := 0; x < imagIdGrad.Rect.Dx(); x++ { 
		    if imag(imagIdGrad.Pix[y*imagIdGrad.Rect.Dy()+x]) != SgrayCosxCosyTiny.Val(x, y) {
		        t.Errorf("Error: ImagIdentity gradient incorrect at %v, %v. Expected %v, got %v\n",
		                 x, y, SgrayCosxCosyTiny.Val(x, y), imag(imagIdGrad.Pix[y*idGrad.Rect.Dy()+x]))
		    }
		}
    }
    kierGrad := FdgradKernel(Sgray, kieransKernel)
    if !reflect.DeepEqual(kierGrad.Pix, smallPicGradKieransKernel) {
        t.Error("Error: Gradient image incorrect. Expected:" +
            ComplexArrayToString(smallPicGradKieransKernel, 3) + "Got:" +
            ComplexArrayToString(kierGrad.Pix, 3))
    }
}