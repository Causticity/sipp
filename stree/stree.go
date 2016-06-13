// Package stree provides a SippNode class for organising SIPP images,
// properties, and operations, as well as functions and variables for managing
// the list of trees. This file contains code for managing the list of trees.

package stree

import (
    //"fmt"
    "image"
)

var trees []*SippNode = make([]*SippNode, 10)

// It would be good to eliminate this variable and the corresponding code.
// The window system should manage this, if we can only ensure functions can
// be called with appropriate receivers.
var currentTreeRootIndex int

func NewTree (url string) {
	//fmt.Println("NewTree called with name: <", url,">")
	firstFree := findFirstFreeTreeSlot()
		//fmt.Println("New Tree: first free index: ", firstFree)
	trees[firstFree] = NewSippRootNode(url)
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

// This should take a node receiver, and possibly move to snode. Or it should
// be broken up into a function that finds the root node for a given node, which
// lives in snode, and a function that destroys the tree given a root node. This
// will have to walk the tree and destroy both the UI and the SippNode for every
// node, then delete the tree, if there even is a representation of the tree
// once we are done converting to receivers.
func CloseTree() {
	//fmt.Println("closing tree at index: ", currentTreeRootIndex)
	trees[currentTreeRootIndex].Close()
	// Use copy and nil out everything after the copy
	moved := copy(trees[currentTreeRootIndex:], trees[currentTreeRootIndex+1:])
	trees[currentTreeRootIndex+moved] = nil
	findWindowWithFocus()
}

// TODO: This doesn't deal with image windows. When an image window has focus,
// it chooses the first tree. Need an explicit nill, I think. But this is the
// first function that should disappear in favour of receivers.
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

// The image-provider mechanism gets in the way of hooking the images up to
// the corresponding UI elements. Solution: Use a map of nodes with a
// guaranteed-unique version of the image name as the map key. Then the image
// source in the UI can include the key, so that when the name is passed to
// the providers below they can look up the correct node.
func SrcProvider(srcName string, width, height int) image.Image {
	return trees[currentTreeRootIndex].Src[0]
}

func ThumbProvider(srcName string, width, height int) image.Image {
	return trees[currentTreeRootIndex].Src[0].Thumbnail()
}
