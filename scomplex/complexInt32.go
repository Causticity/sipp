// Copyright Raul Vera 2020

package scomplex

// A ComplexInt32 is a complex number with int32 real and imaginary parts.
type ComplexInt32 struct {
	Re int32
	Im int32
}

// Adds a ComplexInt32 to this one, returning a new one.
func (a ComplexInt32) Add(b ComplexInt32) ComplexInt32 {
	return ComplexInt32{a.Re+b.Re, a.Im+b.Im}
}

// Subtracts a ComplexInt32 from this one, returning a new one.
func (a ComplexInt32) Sub(b ComplexInt32) ComplexInt32 {
	return ComplexInt32{a.Re-b.Re, a.Im-b.Im}
}

// Multiplies a ComplexInt32 to this one, returning a new one.
func (a ComplexInt32) Mult(b ComplexInt32) ComplexInt32 {
	return ComplexInt32{a.Re*b.Re-a.Im*b.Im, a.Re*b.Im+a.Im*b.Re}
}