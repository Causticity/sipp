package shist

import (
	. "github.com/Causticity/sipp/simage"
	. "github.com/Causticity/sipp/sgrad"
)

type Sipphist struct {
	Grad *Gradimage
	Max float64
}

func Hist(grad *Gradimage, k int) (hist *Sipphist) {
	return
}

func (hist *Sipphist) Render() (rnd *Sippimage) {
	return
}