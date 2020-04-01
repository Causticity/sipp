// Copyright Raul Vera 2016

// Tests for package scomplex

package scomplex

import (
	"reflect"
	"testing"
	)

import (
	. "github.com/Causticity/sipp/sipptesting"
)

func TestComplex (t *testing.T) {
    comp := ToShiftedComplex(Sgray)
    if !reflect.DeepEqual(ShiftedPic, comp.Pix) {
        t.Error("Shifted complex doesn't match Gray!")
    }
    
    comp = ToShiftedComplex(Sgray16)
    if !reflect.DeepEqual(ShiftedPic, comp.Pix) {
        t.Error("Shifted complex doesn't match Gray16!")
    }

    re, im := comp.Render()
    
    if !reflect.DeepEqual(re.Pix(), ScaledShiftedPic) {
        t.Error("real unexpected")
    }
    
    if !reflect.DeepEqual(im, SgrayZero) {
        t.Error("imaginary not zero")
    }
}

