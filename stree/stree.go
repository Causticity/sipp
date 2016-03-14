// Package stree provides a SippNode class for organising SIPP images,
// properties, and operations

package stree

import (
    "fmt"

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
	QmlRoot *qml.Window
	
	// The parameters to the SippOp that got us here.
	Params SippOpParams
	
	// The nodes that have been derived (and retained) from this node.
	// The slice itself is nil at a leaf.
	Children []*SippNode
	
	// The node that this one was derived from, nil at the root of the tree.
	Parent *SippNode
}

// A SippOp specifies a function that takes a source image and an arbitrary
// set of parameters, and returns a slice of result images.
type SippOp interface {
	Op(*SippImage, *SippOpParams) ([]*SippImage)
}

// A type used only to allow storing an arbitrary set of parameters as part of
// a SippNode.
type SippOpParams interface {}

var treeComponent qml.Object

// Package-wide initialisation, but it requires that the QML Engine be 
// available, so it can't be done as an actual package init, but must be 
// called explicitly once the Engine is created inside the QML run method.
func InitTreeComponents(engine *qml.Engine) error {
	var err error
	treeComponent, err = engine.LoadFile("SippTreeRoot.qml")
	return err
}

// NewSippTree initialises a new root SippNode by loading the given file. 
// Returns nil on error.
func NewSippTree(url string) *SippNode {
	newGuy := new(SippNode)
	newGuy.Src = make([]SippImage,1)
	var err error
	newGuy.Src[0], err = Read(&url)
	if err != nil {
		fmt.Println("Error reading image:", err)
		return nil
	}
	newGuy.Params = &url
	return newGuy
}

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
	if newGuy.QmlRoot == nil {
		newGuy.QmlRoot = treeComponent.CreateWindow(nil)
		newGuy.QmlRoot.Set("title", url)
		newGuy.QmlRoot.Call("setThumbSource", url)
		newGuy.QmlRoot.Show()
	}
}