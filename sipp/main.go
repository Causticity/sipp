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
	var out = flag.String("out", "", "output image file prefix")
	var k = flag.Int("K", 63, "Number of bins to scale the max radius to.\nThe histogram will be 2K+1 bins on a side")
	flag.Parse()
	fmt.Println("input file:<", *in, ">")
	fmt.Println("output file prefix:<", *out, ">")
	fmt.Println("histogram edge size:", (*k*2+1))

	var src, err = simage.Read(in)
	if err != nil {
		fmt.Println("Error reading image:", err)
		os.Exit (1)
	}
	fmt.Println("source image read")
	
	grad := sgrad.Fdgrad(src)
	fmt.Println("gradient image computed")

	re, im := grad.Render()
	reName := *out + "_grad_real.png"
	err = re.Write(&reName)
	if err != nil {
		fmt.Println("Error writing real gradient image:", err)
		os.Exit (1)
	}
	imName := *out + "_grad_imag.png"
	err = im.Write(&imName)
	if err != nil {
		fmt.Println("Error writing imag gradient image:", err)
		os.Exit (1)
	}

	hist := shist.Hist(grad, *k)
	
	rhist := hist.Render()
	histName := *out + "_hist.png"	
	err = rhist.Write(&histName)
	if err != nil {
		fmt.Println("Error writing histogram image:", err)
		os.Exit (1)
	}
	histSup := hist.RenderSuppressed()
	histSupName := *out + "_hist_sup.png"	
	err = histSup.Write(&histSupName)
	if err != nil {
		fmt.Println("Error writing suppressed histogram image:", err)
		os.Exit (1)
	}

	fmt.Println("Entropy of the gradient image:", hist.Entropy())
}