// First hacky attempt at getting a UI up and running. This currently doesn't do
// anything useful at all.
package main

import (
	"flag"
    "fmt"
    "image"
    "image/draw"
    "os"

	"github.com/andlabs/ui"

	"github.com/Causticity/sipp/simage"
    //"github.com/Causticity/sipp/sgrad"
    //"github.com/Causticity/sipp/shist"
)

var window ui.Window
var srcName *string
var k *int
var src simage.Sippimage

type areaHandler struct {
	img		*image.RGBA
}

func (a *areaHandler) Paint(rect image.Rectangle) *image.RGBA {
	return a.img.SubImage(rect).(*image.RGBA)
}

func (a *areaHandler) Mouse(me ui.MouseEvent) {
	if me.Up != 0 {
		fullsrc := createImageWindow(src)
		fullsrc.OnClosing(func() bool {
			return true
		})
		fullsrc.Show()
	}
}

func (a *areaHandler) Key(ke ui.KeyEvent) bool { return false }

func createImageWindow(img simage.Sippimage) ui.Window {
	b := img.Bounds()
    handler := areaHandler{image.NewRGBA(b)}
    draw.Draw(handler.img, b, img, b.Min, draw.Src)
	
	area := ui.NewArea(b.Dx(), b.Dy(), &handler)
	window = ui.NewWindow(*srcName, b.Dx(), b.Dy(), area)
	return window
}

func initUI() {
	thumb := src.Thumbnail()
	twin := createImageWindow(thumb)
    twin.OnClosing(func() bool {
        ui.Stop()
        return true
    })
	twin.Show()
}

func main() {
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

    go ui.Do(initUI)
    err = ui.Go()
    if err != nil {
        panic(err)
    }
}