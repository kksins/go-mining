// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package cart implement the Classification and Regression Tree by Breiman, et al.
CART is binary decision tree.

	Breiman, Leo, et al. Classification and regression trees. CRC press,
	1984.

The implementation is based on Data Mining book,

	Han, Jiawei, Micheline Kamber, and Jian Pei. Data mining: concepts and
	techniques: concepts and techniques. Elsevier, 2011.
*/
package cart

import (
	"github.com/shuLhan/go-mining/dataset"
	"github.com/shuLhan/go-mining/gain/gini"
	"github.com/shuLhan/go-mining/set"
	"github.com/shuLhan/go-mining/tree/binary"
)

const (
	// SplitMethodGini if defined in Input, the dataset will be splitted
	// using Gini gain for each possible value or partition.
	//
	// This option is used in Input.SplitMethod.
	SplitMethodGini = 0
)

/*
Input data for building CART.
*/
type Input struct {
	// SplitMethod define the criteria to used for splitting.
	SplitMethod int
	// tree in classification.
	Tree binary.Tree
}

/*
NewInput create new Input object.
*/
func NewInput(SplitMethod int) (*Input) {
	return &Input{
		SplitMethod: SplitMethod,
		Tree: binary.Tree{},
	}
}

/*
BuildTree will create a tree using CART algorithm.
*/
func (in *Input) BuildTree(D *dataset.Input) (e error) {
	in.Tree.Root = in.splitTree(D)

	return e
}

/*
splitTree calculate the gain in all dataset, and split into two node: left and
right.

Return node with the split information.
*/
func (in *Input) splitTree(D *dataset.Input) (node *binary.BTNode) {
	node = &binary.BTNode{}

	// if all dataset is in the same class, return node as leaf with class
	// is set to that class.
	if D.IsInSingleClass() {
		targetAttr := (*D).GetTargetAttrValues()
		node.Value = NodeValue{
				IsLeaf: true,
				Class: (*targetAttr)[0],
				Size: len(*targetAttr),
			}

		return node
	}

	// if dataset is empty return node labeled with majority classes in
	// dataset.
	if D.Size <= 0 {
		majorClass := D.GetMajorityClass()
		node.Value = NodeValue{
			IsLeaf: true,
			Class: majorClass,
		}
	}

	// calculate the Gini gain for each attribute.
	gains := in.computeGiniGain(D)

	// get attribute with maximum Gini gain.
	MaxGainIdx := gini.FindMaxGain(gains)
	MaxGain := (*gains)[MaxGainIdx]

	// using the sorted index in MaxGain, sort all field in dataset
	D.SortByIndex(MaxGainIdx, MaxGain.SortedIndex)

	// Now that we have attribute with max gain in AttrMaxGain, and their
	// gain dan partition value in Gains[AttrMaxGain] and
	// GetMaxPartValue(), we split the dataset based on type of max-gain
	// attribute.
	// If its continuous, split the attribute using numeric value.
	// If its discrete, split the attribute using subset (partition) of
	// nominal values.
	var splitV interface{}

	if MaxGain.IsContinu {
		splitV = MaxGain.GetMaxPartGainValue()
	} else {
		attrPartV := MaxGain.GetMaxPartGainValue()
		attrSubV := attrPartV.(set.SubsetString)
		splitV = attrSubV[0]
	}

	node.Value = NodeValue{
			IsLeaf: false,
			IsContinu: MaxGain.IsContinu,
			Size: D.Size,
			SplitAttrIdx: MaxGainIdx,
			SplitV: splitV,
		}

	splitD := D.SplitByAttrValue(MaxGainIdx, splitV)

	nodeLeft := in.splitTree(splitD)
	nodeRight := in.splitTree(D)

	node.SetLeft(nodeLeft)
	node.SetRight(nodeRight)

	return node
}

/*
computeGiniGain calculate the gini index for each value in each attribute.
*/
func (in *Input) computeGiniGain(D *dataset.Input) (*[]gini.Gini) {
	var gains []gini.Gini

	switch in.SplitMethod {
	case SplitMethodGini:
		// create gains value for all attribute minus target class.
		gains = make([]gini.Gini, len(D.Attrs))
	}

	targetAttr := D.GetTargetAttr()
	targetValues := targetAttr.GetDiscreteValues()
	classes := targetAttr.NominalValues

	for i := range (*D).Attrs {
		// skip class attribute.
		if i == D.ClassIdx {
			continue
		}

		target := make([]string, len(*targetValues))
		copy(target, (*targetValues))

		// compute gain.
		if (*D).Attrs[i].IsContinu {
			attr := (*D).Attrs[i].GetContinuValues()

			gains[i].ComputeContinu(attr, &target, &classes)
		} else {
			attr := (*D).Attrs[i].GetDiscreteValues()
			attrValues := (*D).Attrs[i].NominalValues

			gains[i].ComputeDiscrete(attr, &attrValues, &target,
							&classes)
		}
	}
	return &gains
}

/*
String yes, it will print it JSON like format.
*/
func (in *Input) String() (s string) {
	s = in.Tree.String()
	return s
}
