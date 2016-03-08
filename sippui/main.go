// First hacky attempt at getting a UI up and running. This currently doesn't do
// anything useful at all.
package main

import (
	"flag"
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

var srcName *string
//var k *int
var src simage.Sippimage

func main() {
	/*
	// This will become a parameter to the gradient op.
	k = flag.Int("K", 0, "Number of bins to scale the max radius to. "+
						 "The histogram will be 2K+1 bins on a side.\n"+
						 "        This is used only for 16-bit images.\n"+
						 "        If K is omitted, it is computed from "+
						 "the maximum excursion of the gradient.\n"+
						 "        8-bit images always use a 511x511 histogram, "+
						 "as that covers the entire possible space.")

	// This test will move to the gradient op. Specifically, it won't be 
	// available in the UI for 8-bit images. But it will be displayed in the
	// info display.
	if src.Bpp() == 8 {
		*k = 255
		fmt.Println("Image is 8-bit. K forced to 255.")
	}
	*/
	srcName = flag.String("in", "", "input image file; must be grayscale png")
	flag.Parse()
	if err := qml.Run(run); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

var app *qml.Window
var treeRoot *qml.Window

func run() error {
	engine := qml.NewEngine()
	engine.AddImageProvider("thumb", loadImage)
	engine.AddImageProvider("src", func(id string, width, height int) image.Image {
		return src
	})

	appComponent, err := engine.LoadFile("sippui.qml")
	if err != nil {
		return err
	}
	
	treeComponent, err := engine.LoadFile("SippTreeRoot.qml")
	if err != nil {
		return err
	}
	
	app = appComponent.CreateWindow(nil)
	app.Show()
	
	treeRoot = treeComponent.CreateWindow(nil)
	treeRoot.On("gotFile", imageName)
	
	newTree := app.ObjectByName("newTree")
	newTree.On("triggered", func() {treeRoot.Call("getFile")})

	if len(*srcName) > 0 {
		imageName(*srcName)
	}
	
	app.Wait()

	return nil
}

func imageName (url string) {
	 treeRoot.Call("setThumbSource", url)
	 treeRoot.Set("title", url)
	 treeRoot.Show()
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
