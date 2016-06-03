// Package stree provides a SippNode class for organising SIPP images,
// properties, and operations, as well as functions and variables for managing
// the list of trees. TODO: Separate the list from the node code?

package stree

import (
    "fmt"
    "image"

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
		newGuy.QmlNode.Set("title", url)
		newGuy.QmlNode.Call("setThumbSource", url)
		newGuy.QmlNode.On("focusChanged", findWindowWithFocus)
		newGuy.QmlNode.On("thumbClicked", newGuy.thumbClicked)
		newGuy.QmlNode.On("gradientClicked", newGuy.gradientClicked)
		newGuy.QmlNode.Show()
	}
}

var trees []*SippNode = make([]*SippNode, 10)

var currentTreeRootIndex int

func NewTree (url string) {
	firstFree := findFirstFreeTreeSlot()
		//fmt.Println("New Tree: first free index: ", firstFree)
	trees[firstFree] = NewSippTree(url)
	if trees[firstFree] == nil {
		// NewSippTree will have logged an error, so just return
		return
	}
	currentTreeRootIndex = firstFree
	trees[currentTreeRootIndex].BuildUI(url)
}

func findFirstFreeTreeSlot() int {
	// If the current one isn't filled yet, we're done
	if trees[currentTreeRootIndex] == nil {
		return currentTreeRootIndex
	}
	curLen := len(trees)
	// Are there any other slots available? Return the first one
	for i := currentTreeRootIndex; i < curLen; i++ {
		if trees[i] == nil {
			return i
		}
	}
	// Expand the slice and return the index of the first new slot
	temp := make([]*SippNode, cap(trees)+10)
	copy(temp, trees)
	trees = temp
	return curLen // Index of first available
}

func (victim *SippNode) Close() {
	victim.QmlNode.Destroy()
	if victim.QmlImage != nil {
		victim.QmlImage.Destroy()
	}
}

func CloseTree() {
	//fmt.Println("closing tree at index: ", currentTreeRootIndex)
	trees[currentTreeRootIndex].Close()
	// Use copy and nil out everything after the copy
	moved := copy(trees[currentTreeRootIndex:], trees[currentTreeRootIndex+1:])
	trees[currentTreeRootIndex+moved] = nil
	findWindowWithFocus()
}

// TODO: This doesn't deal with image windows. When an image window has focus,
// it chooses the first tree. Need an explicit nill, I think.
func findWindowWithFocus() (bool, int) {
	//fmt.Println("Finding window with focus")
	// Now find which one got focus in the UI and set current
	for i := 0; i < len(trees) && trees[i] != nil; i++ {
		if trees[i].QmlNode.Bool("active") == true {
			currentTreeRootIndex = i
			return true, 0
		}
	}
	currentTreeRootIndex = 0 // TODO: maybe separate first from none, explicitly?
	return true, 0
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

func SrcProvider(srcName string, width, height int) image.Image {
	return trees[currentTreeRootIndex].Src[0]
}

func ThumbProvider(srcName string, width, height int) image.Image {
	return trees[currentTreeRootIndex].Src[0].Thumbnail()
}
