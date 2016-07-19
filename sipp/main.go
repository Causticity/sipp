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

	var in = flag.String("in", "", "Input image file; must be grayscale png")
	var out = flag.String("out", "", "Output image file prefix")
	var thb = flag.Bool("t", false, "Boolean; if true, write a thumbnail image")
	var grd = flag.Bool("g", false, "Boolean; if true, write the gradient"+
									" real and imaginary images")
	var hst = flag.Bool("h", false, "Boolean; if true, write a histogram image")
	var hsp = flag.Bool("hs", false, "Boolean; if true, write a histogram"+
									 " image with the center spike suppressed")
	var hse = flag.Bool("he", false, "Boolean; if true, write a histogram"+
									 "-entropy image")
	var gre = flag.Bool("ge", false, "Boolean; if true, write a gradient"+
									 "-entropy image")
	var e = flag.Bool("e", false, "Boolean; if true, write a conventional"+
									 " entropy image")
	var f = flag.Bool("f", false, "Boolean; if true, write the fft"+
									" real and imaginary images")
	var fls = flag.Bool("fls", false, "Boolean; if true, write the fft"+
									" log spectrum image")
	var a = flag.Bool("a", false, "Boolean; if true, write all the images")
	var k = flag.Int("K", 0, "Number of bins to scale the max radius to. "+
							  "The histogram will be 2K+1 bins on a side.\n"+
							  "        This is used only for 16-bit images.\n"+
							  "        If K is omitted, it is computed from "+
							  "the maximum excursion of the gradient.\n"+
							  "        8-bit images always use a 511x511 histogram, "+
							  "as that covers the entire possible space.")
	flag.Parse()
	if *a {
		*thb = true
		*grd = true
		*hst = true
		*hsp = true
		*hse = true
		*gre = true
		*e = true
		*f = true
		*fls = true
	}
	
	fmt.Println("input file:<", *in, ">")
	fmt.Println("output file prefix:<", *out, ">")

	src, err := simage.Read(in)
	if err != nil {
		fmt.Println("Error reading image:", err)
		os.Exit (1)
	}
	fmt.Println("source image read")
	
	if *thb {
		thumb := src.Thumbnail()
		fmt.Println("Thumbnail generated")
		tname := *out + "_thumb.png"
		err = thumb.Write(&tname)
		if err != nil {
			fmt.Println("Error writing thumbnail image:", err)
			os.Exit (1)
		}
	}
	
	if src.Bpp() == 8 {
		*k = 255
		fmt.Println("Image is 8-bit. K forced to 255.")
	}
	
	grad := sgrad.Fdgrad(src)
	fmt.Println("gradient image computed")

	if *grd {
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
	}

	hist := shist.Hist(grad, *k)
	
	if *hst {
		rhist := hist.Render()
		histName := *out + "_hist.png"	
		err = rhist.Write(&histName)
		if err != nil {
			fmt.Println("Error writing histogram image:", err)
			os.Exit (1)
		}
	}
	
	if *hsp {
		histSup := hist.RenderSuppressed()
		histSupName := *out + "_hist_sup.png"	
		err = histSup.Write(&histSupName)
		if err != nil {
			fmt.Println("Error writing suppressed histogram image:", err)
			os.Exit (1)
		}
	}

	gradEnt := hist.GradEntropy()
	fmt.Println("Entropy of the gradient image:", gradEnt)
	
	if *hse {
		histEntImg := hist.HistEntropyImage()
		histEntName := *out + "_hist_ent.png"
		err = histEntImg.Write(&histEntName)
		if err != nil {
			fmt.Println("Error writing the histogram entropy image", err)
			os.Exit (1)
		}
	}
	
	if *gre {
		gradEntImg := hist.GradEntropyImage()
		gradEntName := *out + "_grad_ent.png"
		err = gradEntImg.Write(&gradEntName)
		if err != nil {
			fmt.Println("Error writing the gradient entropy image", err)
			os.Exit (1)
		}
	}
	
	ent, entImg := shist.Entropy(src)
	fmt.Println("Conventional entropy of the source image:", ent)
	
	if *e {
		entName := *out + "_conv_ent.png"
		err = entImg.Write(&entName)
		if err != nil {
			fmt.Println("Error writing the conventional entropy image", err)
			os.Exit (1)
		}
	}
	
	fft := sfft.FFT(src)
	fmt.Println("fft computed");
	
	if *f {
		re, im := fft.Render()
		reName := *out + "_fft_real.png"
		imName := *out + "_fft_imag.png"
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
	}
	
	if *fls {
		ls := sfft.LogSpectrum(fft)
		fmt.Println("Log spectrum computed")
		lsName := *out + "_fft_spectrum.png"
		err = ls.Write(&lsName)
		if err != nil {
			fmt.Println("Error writing fft spectrum image:", err)
			os.Exit (1)
		}
	}
	
	elapsed := time.Since(start)
	fmt.Println("Elapsed time:" + elapsed.String())
	
}