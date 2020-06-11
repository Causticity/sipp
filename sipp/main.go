// Copyright Raul Vera 2015-2016

package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

import (
	"github.com/Causticity/sipp/sentropy"
	"github.com/Causticity/sipp/sfft"
	"github.com/Causticity/sipp/sgrad"
	"github.com/Causticity/sipp/shist"
	"github.com/Causticity/sipp/simage"
)

func main() {

	start := time.Now()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println("Source code for this program may be found at (https://github.com/Causticity/sipp)")
	}

	var in = flag.String("in", "", "Input image file; must be grayscale png")
	var out = flag.String("out", "", "Output image file prefix")
	var thb = flag.Bool("t", false, "Boolean; if true, write a thumbnail image")
	var grd = flag.Bool("g", false, "Boolean; if true, write the gradient"+
		" real and imaginary images")
	var hst = flag.Bool("h", false, "Boolean; if true, write a histogram image")
	var hsp = flag.Bool("hs", false, "Boolean; if true, write a histogram"+
		" image with the center spike suppressed")
	var hde = flag.Bool("hde", false, "Boolean; if true, write a histogram"+
		"delentropy image")
	var de = flag.Bool("de", false, "Boolean; if true, write a delentropy"+
		" image")
	var e = flag.Bool("e", false, "Boolean; if true, write a conventional"+
		" entropy image")
	var f = flag.Bool("f", false, "Boolean; if true, write the fft"+
		" real and imaginary images")
	var fls = flag.Bool("fls", false, "Boolean; if true, write the fft"+
		" log spectrum image")
	var a = flag.Bool("a", false, "Boolean; if true, write all the images")
	var radius = flag.Uint("r", 0, "Number of bins to scale the max radius to. "+
		"The histogram will be 2r+1 bins on a side.\n"+
		"        This is used only for 16-bit images.\n"+
		"        If r is omitted, it is computed from "+
		"the maximum excursion of the gradient.\n")
	var rf8 = flag.Bool("rf8", false, "Boolean; if true, force 8-bit images "+
		" to use a 511x511 histogram, "+
		"as that covers the entire possible space.")
	var v = flag.Bool("v", false, "Boolean; if true, verbosely report "+
		"everything done")
	var csv = flag.Bool("csv", false, "Boolean: if true, write the name of the"+
		"image, a comma, and the delentropy,"+
		"on a single line.")

	flag.Parse()
	if *a {
		*thb = true
		*grd = true
		*hst = true
		*hsp = true
		*hde = true
		*de = true
		*e = true
		*f = true
		*fls = true
	}

	if *v {
		fmt.Println("input file:<", *in, ">")
		fmt.Println("output file prefix:<", *out, ">")
	}

	src, err := simage.Read(*in)
	if err != nil {
		fmt.Println("Error reading image:", err)
		os.Exit(1)
	}
	if *v {
		fmt.Println("source image read")
	}

	if *thb {
		thumb := src.Thumbnail()
		if *v {
			fmt.Println("Thumbnail generated")
		}
		tname := *out + "_thumb.png"
		err = thumb.Write(&tname)
		if err != nil {
			fmt.Println("Error writing thumbnail image:", err)
			os.Exit(1)
		}
	}

	if src.Bpp() == 8 && *rf8 {
		*radius = 255
		if *v {
			fmt.Println("Rdius forced to 255 for 8-bit image.")
		}
	}

	grad := sgrad.Fdgrad(src)
	if *v {
		fmt.Println("gradient image computed")
	}

	if *grd {
		re, im := grad.Render()
		reName := *out + "_grad_real.png"
		err = re.Write(&reName)
		if err != nil {
			fmt.Println("Error writing real gradient image:", err)
			os.Exit(1)
		}
		imName := *out + "_grad_imag.png"
		err = im.Write(&imName)
		if err != nil {
			fmt.Println("Error writing imag gradient image:", err)
			os.Exit(1)
		}
	}

	hist := shist.Hist(grad, uint16(*radius))

	if *hst {
		rhist := hist.Render()
		histName := *out + "_hist.png"
		err = rhist.Write(&histName)
		if err != nil {
			fmt.Println("Error writing histogram image:", err)
			os.Exit(1)
		}
	}

	if *hsp {
		histSup := hist.RenderSuppressed()
		histSupName := *out + "_hist_sup.png"
		err = histSup.Write(&histSupName)
		if err != nil {
			fmt.Println("Error writing suppressed histogram image:", err)
			os.Exit(1)
		}
	}

	sippDel := sentropy.Delentropy(hist)
	delentropy := sippDel.Delentropy / 2.0

	if *csv {
		fmt.Printf("%s,%.2f\n", *in, delentropy)
	} else {
		fmt.Println("Delentropy:", delentropy)
	}
	if *hde {
		histEntImg := sippDel.HistDelentropyImage()
		histEntName := *out + "_hist_delent.png"
		err = histEntImg.Write(&histEntName)
		if err != nil {
			fmt.Println("Error writing the histogram delentropy image", err)
			os.Exit(1)
		}
	}

	if *de {
		delEntImg := sippDel.DelEntropyImage()
		delEntName := *out + "_delent.png"
		err = delEntImg.Write(&delEntName)
		if err != nil {
			fmt.Println("Error writing the delentropy image", err)
			os.Exit(1)
		}
	}

	ent, entImg := sentropy.Entropy(src)
	if *v {
		fmt.Println("Conventional entropy of the source image:", ent)
	}

	if *e {
		entName := *out + "_conv_ent.png"
		err = entImg.Write(&entName)
		if err != nil {
			fmt.Println("Error writing the conventional entropy image", err)
			os.Exit(1)
		}
	}

	fft := sfft.FFT(src)
	if *v {
		fmt.Println("fft computed")
	}

	if *f {
		re, im := fft.Render()
		reName := *out + "_fft_real.png"
		imName := *out + "_fft_imag.png"
		err = re.Write(&reName)
		if err != nil {
			fmt.Println("Error writing real fft image:", err)
			os.Exit(1)
		}
		err = im.Write(&imName)
		if err != nil {
			fmt.Println("Error writing imag fft image:", err)
			os.Exit(1)
		}
	}

	if *fls {
		ls := sfft.LogSpectrum(fft)
		fmt.Println("Log spectrum computed")
		lsName := *out + "_fft_spectrum.png"
		err = ls.Write(&lsName)
		if err != nil {
			fmt.Println("Error writing fft spectrum image:", err)
			os.Exit(1)
		}
	}

	elapsed := time.Since(start)
	if *v {
		fmt.Println("Elapsed time:" + elapsed.String())
	}
}
