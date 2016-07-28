// Copyright Raul Vera 2016

// Tests for package scomplex

package scomplex

import (
	. "image"
	"testing"
	)

import (
	. "github.com/Causticity/sipp/simage"
)

// A small image

const var smallPix []uint8 {
  1,  2,  3,  4,
  5,  6,  7,  8,
  9, 10, 11, 12,
 13, 14, 15, 16
}
 
const var gray image.Gray {
  smallPix, 
  4, 
  Rectangle{Point{0, 0}, Point{4, 4}
}

