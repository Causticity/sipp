// Copyright Raul Vera 2015-2016

package main

import (
	"flag"
    "fmt"
    "os"
    "time"
)

import (
	"github.com/Causticity/sipp/simage"
    "github.com/Causticity/sipp/sgrad"
    "github.com/Causticity/sipp/shist"
    "github.com/Causticity/sipp/sfft"
)

func main() {

	start := time.Now()

	var in = flag.String("in", "", "input image file; must be grayscale png")
	var out = flag.String("out", "", "output image file prefix")
	var k = flag.Int("K", 0, "Number of bins to scale the max radius to. "+
							  "The histogram will be 2K+1 bins on a side.\n"+
							  "        This is used only for 16-bit images.\n"+
							  "        If K is omitted, it is computed from "+
							  "the maximum excursion of the gradient.\n"+
							  "        8-bit images always use a 511x511 histogram, "+
							  "as that covers the entire possible space.")
	flag.Parse()
	fmt.Println("input file:<", *in, ">")
	fmt.Println("output file prefix:<", *out, ">")

	src, err := simage.Read(in)
	if err != nil {
		fmt.Println("Error reading image:", err)
		os.Exit (1)
	}
	fmt.Println("source image read")
	
	thumb := src.Thumbnail()
	fmt.Println("Thumbnail generated")
	tname := *out + "_thumb.png"
	err = thumb.Write(&tname)
	if err != nil {
		fmt.Println("Error writing thumbnail image:", err)
		os.Exit (1)
	}
	
	if src.Bpp() == 8 {
		*k = 255
		fmt.Println("Image is 8-bit. K forced to 255.")
	}
	
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

	ent, entImg := hist.Entropy()
	fmt.Println("Entropy of the gradient image:", ent)
	entName := *out + "_hist_ent.png"
	err = entImg.Write(&entName)
	if err != nil {
		fmt.Println("Error writing the histogram entropy image", err)
		os.Exit (1)
	}
		
	fft := sfft.FFT(src)
	fmt.Println("fft computed; rendering:");
	re, im = fft.Render()
	reName = *out + "_fft_real.png"
	imName = *out + "_fft_imag.png"
	err = re.Write(&reName)
	if err != nil {
		fmt.Println("Error writing real fft image:", err)
		os.Exit (1)
	}
	err = im.Write(&imName)
	if err != nil {
		fmt.Println("Error writing imag fft image:", err)
		os.Exit (1)
	}
	
	ls := sfft.LogSpectrum(fft)
	fmt.Println("Log spectrum computed")
	lsName := *out + "_fft_spectrum.png"
	err = ls.Write(&lsName)
	if err != nil {
		fmt.Println("Error writing fft spectrum image:", err)
		os.Exit (1)
	}

	elapsed := time.Since(start)
	fmt.Println("Elapsed time:" + elapsed.String())
	
}