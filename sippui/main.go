// First hacky attempt at getting a UI up and running. This currently doesn't do
// anything useful at all.
package main

import (
//	"flag"
    "fmt"
    "image"
//    "image/draw"
    "os"

	"gopkg.in/qml.v1"

	"github.com/Causticity/sipp/simage"
	//"github.com/Causticity/sipp/stree"
    //"github.com/Causticity/sipp/sgrad"
    //"github.com/Causticity/sipp/shist"
)

//var srcName *string
//var k *int
var src simage.Sippimage


func main() {
	/*
	srcName = flag.String("in", "", "input image file; must be grayscale png")
	k = flag.Int("K", 0, "Number of bins to scale the max radius to. "+
						 "The histogram will be 2K+1 bins on a side.\n"+
						 "        This is used only for 16-bit images.\n"+
						 "        If K is omitted, it is computed from "+
						 "the maximum excursion of the gradient.\n"+
						 "        8-bit images always use a 511x511 histogram, "+
						 "as that covers the entire possible space.")
	flag.Parse()
	fmt.Println("input file:<", *srcName, ">")

	var err error
	src, err = simage.Read(srcName)
	if err != nil {
		fmt.Println("Error reading image:", err)
		os.Exit (1)
	}
	fmt.Println("source image read")

	if src.Bpp() == 8 {
		*k = 255
		fmt.Println("Image is 8-bit. K forced to 255.")
	}
	*/
	if err := qml.Run(run); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	engine := qml.NewEngine()
	engine.AddImageProvider("thumb", loadImage)
	engine.AddImageProvider("src", func(id string, width, height int) image.Image {
		return src
	})

	component, err := engine.LoadFile("sippui.qml")
	if err != nil {
		return err
	}

	win := component.CreateWindow(nil)
	win.Show()
	win.Wait()

	return nil
}

func loadImage(srcName string, width, height int) image.Image {
	fmt.Println("input file selected:<", srcName, ">")

	var err error
	src, err = simage.Read(&srcName)
	if err != nil {
		fmt.Println("Error reading image:", err)
		os.Exit (1)
	}
	//fmt.Println("source image read; returning thumbnail")
	return src.Thumbnail()
}
