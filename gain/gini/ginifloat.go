// Copyright 2016 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gini

import (
	"fmt"
	"github.com/shuLhan/numerus"
)

//
// ComputeContinuFloat Given an attribute A and the target attribute T which contain
// N classes in C, compute the information gain of A.
//
// The result of Gini partitions value, Gini Index, and Gini Gain is saved in
// ContinuPart, Index, and Gain.
//
// Algorithm,
// (0) Sort the attribute.
// (1) Sort the target attribute using sorted index.
// (2) Create continu partition.
// (3) Create temporary space for gini index and gini gain.
// (4) Compute gini index for all target.
//
func (gini *Gini) ComputeContinuFloat(A, T, C *[]float64) {
	gini.IsContinu = true

	gini.SortedIndex = numerus.Floats64IndirectSort(*A, true)

	if DEBUG >= 1 {
		fmt.Println("[gini] attr sorted :", A)
	}

	// (1)
	numerus.Floats64SortByIndex(T, gini.SortedIndex)

	// (2)
	gini.createContinuPartition(A)

	// (3)
	gini.Index = make([]float64, len(gini.ContinuPart))
	gini.Gain = make([]float64, len(gini.ContinuPart))
	gini.MinIndexValue = 1.0

	// (4)
	gini.Value = gini.computeFloat(T, C)

	gini.computeContinuGainFloat(A, T, C)
}

//
// computeFloat will compute Gini value for attribute T.
//
// Gini value is computed using formula,
//
//	1 - sum (probability of each classes in T)
//
func (gini *Gini) computeFloat(T, C *[]float64) float64 {
	n := float64(len(*T))
	if n == 0 {
		return 0
	}

	classCount := numerus.Floats64Counts(*T, *C)

	var sump2 float64

	for x, v := range classCount {
		p := float64(v) / n
		sump2 += (p * p)

		if DEBUG >= 3 {
			fmt.Printf("[gini] compute (%s): (%f/%f)^2 = %f\n",
				(*C)[x], v, n, p*p)
		}

	}

	return 1 - sump2
}

//
// computeContinuGainFloat will compute gain for each partition.
//
// The Gini gain formula we used here is,
//
//	Gain(part,S) = Gini(S) - ((count(left)/S * Gini(left))
//				+ (count(right)/S * Gini(right)))
//
// where,
//	- left is sub-sample from S that is less than part value.
//	- right is sub-sample from S that is greater than part value.
//
// Algorithm,
// (0) For each partition value,
// (0.1) Find the split of samples between partition based on partition value.
// (0.2) Count class in partition.
//
func (gini *Gini) computeContinuGainFloat(A, T, C *[]float64) {
	var gainLeft, gainRight float64
	var tleft, tright []float64

	nsample := len(*A)

	if DEBUG >= 2 {
		fmt.Println("[gini] sorted data:", A)
		fmt.Println("[gini] Gini.Value:", gini.Value)
	}

	// (0)
	for p, contVal := range gini.ContinuPart {

		// (0.1)
		partidx := nsample
		for x, attrVal := range *A {
			if attrVal > contVal {
				partidx = x
				break
			}
		}

		nleft := float64(partidx)
		nright := float64(nsample - partidx)
		probLeft := nleft / float64(nsample)
		probRight := nright / float64(nsample)

		if partidx > 0 {
			tleft = (*T)[0:partidx]
			tright = (*T)[partidx:]

			gainLeft = gini.computeFloat(&tleft, C)
			gainRight = gini.computeFloat(&tright, C)
		} else {
			tleft = nil
			tright = (*T)[0:]

			gainLeft = 0
			gainRight = gini.computeFloat(&tright, C)
		}

		// (0.2)
		gini.Index[p] = ((probLeft * gainLeft) +
			(probRight * gainRight))
		gini.Gain[p] = gini.Value - gini.Index[p]

		if DEBUG >= 3 {
			fmt.Println("[gini] tleft:", tleft)
			fmt.Println("[gini] tright:", tright)

			fmt.Printf("[gini] GiniGain(%v) = %f - (%f * %f) + (%f * %f) = %f\n",
				contVal, gini.Value, probLeft, gainLeft,
				probRight, gainRight, gini.Gain[p])
		}

		if gini.MinIndexValue > gini.Index[p] && gini.Index[p] != 0 {
			gini.MinIndexValue = gini.Index[p]
			gini.MinIndexPart = p
		}

		if gini.MaxGainValue < gini.Gain[p] {
			gini.MaxGainValue = gini.Gain[p]
			gini.MaxPartGain = p
		}
	}
}
