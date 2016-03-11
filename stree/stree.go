// Package stree provides a SippNode class for organising SIPP images,
// properties, and operations

package stree

import (
	. "github.com/Causticity/sipp/simage"
)

// We store all images as a tree, retaining the operations and parameters that
// produced each image from its parent in the tree. A SippNode is a node in the
// tree.
type SippNode struct {
	// The Image at this node.
	src *SippImage
	// The operation that got us here from the Parent node.
	Op *SippOp
	
	// The UI object that corresponds to this node.
	
	Params *SippOpParams
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