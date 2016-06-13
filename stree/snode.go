// Package stree provides a SippNode class for organising SIPP images,
// properties, and operations, as well as functions and variables for managing
// the list of trees. This file contains code for SippNodes.

package stree

import (
    "fmt"
    //"image"
    "path/filepath"
	"gopkg.in/qml.v1"

	. "github.com/Causticity/sipp/simage"
)

// We store all images as a tree, retaining the operations and parameters that
// produced each image from its parent in the tree. A SippNode is a node in the
// tree.
type SippNode struct {
	
	// The Image(s) at this node.
	Src []SippImage
	
	// The operation that got us here from the Parent node, i.e. Src = Op(...)
	SippOp
	
	// The UI object that corresponds to this node.
	QmlNode *qml.Window
	
	// The UI object(s) that correspond(s) to the window(s) displaying the full
	// image(s) for this node.
	QmlImage *qml.Window
	
	// The parameters to the SippOp that got us here.
	Params SippOpParams
	
	// The nodes that have been derived (and retained) from this node.
	// The slice itself is nil at a leaf.
	Children []*SippNode
	
	// The node that this one was derived from, nil at the root of the tree.
	Parent *SippNode
	
	// The name of this image, to be used for window titles and as a key
	// in a map of nodes. Must be unique so that a map can work as a lookup
	// mechanism.
	Name string
}

// A SippOp specifies a function that takes a slice of source images and an 
// arbitrary set of parameters, and returns a slice of result images.
type SippOp interface {
	Op([]*SippImage, *SippOpParams) ([]*SippImage)
}

// A type used only to allow storing an arbitrary set of parameters as part of
// a SippNode.
type SippOpParams interface {}

var treeComponent qml.Object
var srcImageComponent qml.Object

// Package-wide initialisation, but it requires that the QML Engine be 
// available, so it can't be done as an actual package init, but must be 
// called explicitly once the Engine is created inside the QML run method.
func InitTreeComponents(engine *qml.Engine) error {
	var err error
	treeComponent, err = engine.LoadFile("SippTreeRoot.qml")
	if err != nil {
		return err
	}
	srcImageComponent, err = engine.LoadFile("SrcImageViewer.qml")
	if err != nil {
		return err
	}
	
	return err
}

// NewSippRootNode initialises a new root SippNode by loading the given file. 
// Returns nil on error.
func NewSippRootNode(url string) *SippNode {
	newGuy := new(SippNode)
	newGuy.Src = make([]SippImage,1)
	var err error
	// strip off the "file://" prefix by referencing the string from index 7
	filename := url[7:]
	newGuy.Src[0], err = Read(&filename)
	if err != nil {
		fmt.Println("Error reading image:", err)
		return nil
	}
	newGuy.Params = url
	newGuy.Name = filepath.Base(url)
	return newGuy
}

var xBase, yBase int

// BuildUI sets up the QML elements for this tree. As some of the setup done
// here can result in callbacks that might depend on the return value of
// NewSippTree, this is broken out into a separate function to be called
// once the node is obtained from NewSippTree. Panics if InitTreeComponents
// has not been called. Can be called multiple times on the same object; if a
// window already exists, this does nothing.
func (newGuy *SippNode) BuildUI(url string) {
	if treeComponent == nil {
		panic("InitTreeComponents must be called before building tree UIs.")
	}
	if newGuy.QmlNode == nil {
		newGuy.QmlNode = treeComponent.CreateWindow(nil)
		if xBase == 0 {
			xBase = newGuy.QmlNode.Int("x")
			yBase = newGuy.QmlNode.Int("y")
		} else {
			xBase += 40
			yBase += 40
			newGuy.QmlNode.Set("x", xBase)
			newGuy.QmlNode.Set("y", yBase)
		}
		newGuy.QmlNode.Set("title", newGuy.Name)
		newGuy.QmlNode.Call("setThumbSource", url)
		newGuy.QmlNode.On("focusChanged", findWindowWithFocus)
		newGuy.QmlNode.On("thumbClicked", newGuy.thumbClicked)
		newGuy.QmlNode.On("gradientClicked", newGuy.gradientClicked)
		newGuy.QmlNode.Show()
	}
}

func (victim *SippNode) Close() {
	victim.QmlNode.Destroy()
	if victim.QmlImage != nil {
		victim.QmlImage.Destroy()
	}
}

func (victim *SippNode) thumbClicked() {
	if victim.QmlImage == nil {
		victim.QmlImage = srcImageComponent.CreateWindow(nil)
		victim.QmlImage.Call("open", "image://src/")
	}
}

// These have to be able to take a receiver, or this whole thing doesn't work.
func (victim *SippNode) gradientClicked() {
	fmt.Println("I'ma do you a gradient!")
}

