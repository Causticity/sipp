package main

import (
	"flag"
    "fmt"
    "os"
)

import (
	"github.com/Causticity/sipp/simage"
    "github.com/Causticity/sipp/sgrad"
    "github.com/Causticity/sipp/shist"
)

func main() {
	var in = flag.String("in", "", "input image file; must be grayscale png")
	var out = flag.String("out", "", "output image file")
	var k = flag.Int("K", 63, "Number of bins to scale the max radius to.\nThe histogram will be 2K+1 bins on a side")
	flag.Parse()
	fmt.Println("input file:<", *in, ">")
	fmt.Println("output file:<", *out, ">")
	fmt.Println("histogram edge size:", (*k*2+1))

	var src, err = simage.Read(in)
	if err != nil {
		fmt.Println("Error reading image:", err)
		os.Exit (1)
	}
	fmt.Println("source image read")
	
	grad := sgrad.Fdgrad(src)
	fmt.Println("gradient image computed")
	
	hist := shist.Hist(grad, *k)
	
	rhist := hist.Render()
		
	err = rhist.Write(out)
	if err != nil {
		fmt.Println("Error writing scatter-plot image:", err)
		os.Exit (1)
	}
}