// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cart

import (
	"fmt"
	"reflect"
)

/*
NodeValue of tree in CART.
*/
type NodeValue struct {
	// IsLeaf define whether node is a leaf or not.
	IsLeaf bool
	// Class of leaf node.
	Class string
	// IsContinu define whether the node split is continuous or discrete.
	IsContinu bool
	// Size define number of value that this node hold.
	Size int
	// SplitAttrIdx define the attribute which cause the split.
	SplitAttrIdx int
	// SplitV define the split value.
	SplitV interface{}
}

/*
String will return the value of node for printable.
*/
func (nodev *NodeValue) String() (s string) {
	if nodev.IsLeaf {
		s = fmt.Sprintf("Class: %s", nodev.Class)
	} else {
		s = fmt.Sprintf("(SplitValue: %v)",
			reflect.ValueOf(nodev.SplitV))
	}

	return s
}
