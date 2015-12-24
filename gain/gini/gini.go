// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package gini contain function to calculating Gini gain.

Gini gain, which is an impurity-based criterion that measures the divergences
between the probability distribution of the target attributes' values.
*/
package gini

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/shuLhan/dsv/util"
	"github.com/shuLhan/go-mining/set"
)

/*
Gini contain slice of sorted index, slice of partition values, slice of Gini
index, Gini value for all samples.
*/
type Gini struct {
	// Skip if its true, the gain value would not be searched on this
	// instance.
	Skip bool
	// IsContinue define whether the Gini index came from continuous
	// attribute or not.
	IsContinu bool
	// Value of Gini index for all value in attribute.
	Value float64
	// MaxPartGain contain the index of partition which have the maximum
	// gain.
	MaxPartGain int
	// MaxGainValue contain maximum gain of index.
	MaxGainValue float64
	// MinIndexPart contain the index of partition which have the minimum
	// Gini index.
	MinIndexPart int
	// MinIndexGini contain minimum Gini index value.
	MinIndexValue float64
	// SortedIndex of attribute, sorted by values of attribute. This will
	// be used to reference the unsorted target attribute.
	SortedIndex []int
	// ContinuPart contain list of partition value for continuous attribute.
	ContinuPart []float64
	// DiscretePart contain the possible combination of discrete values.
	DiscretePart set.SetString
	// Index contain list of Gini Index for each partition.
	Index []float64
	// Gain contain information gain for each partition.
	Gain []float64
}

/*
ComputeDiscrete Given an attribute A with discreate value 'discval', and the
target attribute T which contain N classes in C, compute the information gain
of A.

The result is saved as gain value in MaxGainValue for each partition.
*/
func (gini *Gini) ComputeDiscrete(A *[]string, discval *[]string, T *[]string,
				C *[]string) {
	gini.IsContinu = false

	// create partition for possible combination of discrete values.
	gini.createDiscretePartition((*discval))

	glog.V(2).Infoln("part : ", gini.DiscretePart)

	gini.Index = make([]float64, len(gini.DiscretePart))
	gini.Gain = make([]float64, len(gini.DiscretePart))
	gini.MinIndexValue = 1.0

	// compute gini index for all samples
	gini.Value = gini.compute(T, C)

	gini.computeDiscreteGain(A, T, C)
}

/*
computeDiscreteGain will compute Gini index and Gain for each partition.
*/
func (gini *Gini) computeDiscreteGain(A *[]string, T *[]string, C *[]string) {
	// number of samples
	nsample := float64(len(*A))

	if glog.V(2) {
		glog.Infoln("sample:", T)
		glog.Infof("Gini(a=%s) = %f\n", (*A), gini.Value)
	}

	// compute gini index for each discrete values
	for i,subPart := range gini.DiscretePart {
		// check if sub partition has at least an element
		if len(subPart) <= 0 {
			continue
		}

		sumGI := 0.0
		for _,part := range subPart {
			ndisc := 0.0
			var subT []string

			for _,el := range part {
				for t,a := range (*A) {
					if a != el {
						continue
					}

					// count how many sample with this discrete value
					ndisc++
					// split the target by discrete value
					subT = append(subT, (*T)[t])
				}
			}

			// compute gini index for subtarget
			giniIndex := gini.compute(&subT, C)

			// compute probabilites of discrete value through all samples
			p := ndisc / nsample

			probIndex := p * giniIndex

			// sum all probabilities times gini index.
			sumGI += probIndex

			if glog.V(2) {
				glog.Infoln("subsample:", subT)
				glog.Infof("Gini(a=%s) = %f/%f * %f = %f\n",
						part, ndisc, nsample,
						giniIndex, probIndex)
			}
		}

		gini.Index[i] = sumGI
		gini.Gain[i] = gini.Value - sumGI

		if glog.V(2) {
			glog.Infoln("sample:", subPart)
			glog.Infof("Gain(a=%s) = %f - %f = %f\n",
					subPart, gini.Value, sumGI,
					gini.Gain[i])
		}

		if gini.MinIndexValue > gini.Index[i] && gini.Index[i] != 0 {
			gini.MinIndexValue = gini.Index[i]
			gini.MinIndexPart = i
		}

		if gini.MaxGainValue < gini.Gain[i] {
			gini.MaxGainValue = gini.Gain[i]
			gini.MaxPartGain = i
		}
	}
}

/*
createDiscretePartition will create possible combination for discrete value
in DiscretePart.
*/
func (gini *Gini) createDiscretePartition(discval []string) {
	// no discrete values ?
	if len(discval) <= 0{
		return
	}

	// use set partition function to group the discrete values into two
	// subset.
	gini.DiscretePart = set.PartitioningSetString(discval, 2)
}

/*
ComputeContinu Given an attribute A and the target attribute T which contain
N classes in C, compute the information gain of A.

The result of Gini partitions value, Gini Index, and Gini Gain is saved in
ContinuPart, Index, and Gain.
*/
func (gini *Gini) ComputeContinu(A *[]float64, T *[]string, C *[]string) {
	gini.IsContinu = true

	// make a copy of attribute and target.
	A2 := make([]float64, len(*A))
	copy(A2, *A)

	T2 := make([]string, len(*T))
	copy(T2, *T)

	gini.SortedIndex = util.IndirectSortFloat64(A2)

	if glog.V(1) {
		glog.Infoln(">>> attr :", A2)
	}

	// sort the target attribute using sorted index.
	util.SortStringSliceByIndex(&T2, &gini.SortedIndex)

	// create partition
	gini.createContinuPartition (&A2)

	// create holder for gini index and gini gain
	gini.Index = make([]float64, len(gini.ContinuPart))
	gini.Gain = make([]float64, len(gini.ContinuPart))
	gini.MinIndexValue = 1.0

	// compute gini index for all samples
	gini.Value = gini.compute (&T2, C)

	gini.computeContinuGain(&A2, &T2, C)
}

/*
createContinuPart for dividing class and computing Gini index.

This is assuming that the data `A` has been sorted in ascending order.
*/
func (gini *Gini) createContinuPartition(A *[]float64) {
	l := len(*A)
	gini.ContinuPart = make([]float64, 0)

	// loop from first index until last index - 1
	for i := 0; i < l - 1; i++ {
		sum := (*A)[i] + (*A)[i + 1]
		med := sum / 2.0

		// If median is zero, its mean both left and right value is
		// zero. We are not allowing this, because it will result the
		// minimum Gini Index or maximum Gain value.
		if med == 0 {
			continue
		}

		// Reject if median is contained in attribute's value.
		// We use equality because if both A[i] and A[i+1] value is
		// equal, the median is equal to both of them.
		exist := false
		for j := 0; j <= i; j++ {
			if (*A)[j] == med {
				exist = true
				break
			}
		}
		if ! exist {
			gini.ContinuPart = append(gini.ContinuPart, med)
		}
	}
}

/*
compute value for attribute T.

Return Gini value in the form of,

	1 - sum (probability of each classes in T)
*/
func (gini *Gini) compute(T *[]string, C *[]string) float64 {
	n := float64 (len(*T))

	if n == 0 {
		return 0
	}

	classCount := make(map[string]int, len(*C))

	for _,v := range *T {
		classCount[v]++
	}

	var p float64
	var sump2 float64

	for i := range classCount {
		p = float64(classCount[i]) / n
		sump2 += (p * p)

		if glog.V(2) {
			glog.Infof(" compute (%s): (%f/%f)^2 = %f\n", *T,
					float64(classCount[i]), n, p * p)
		}

	}

	return 1 - sump2
}

/*
computeContinuGain for each partition.

The Gini gain formula we used here is,

	Gain(part,S) = Gini(S) - ((count(left)/S * Gini(left))
				+ (count(right)/S * Gini(right)))

where,
	- left is sub-sample from S that is less than part value.
	- right is sub-sample from S that is greater than part value.
*/
func (gini *Gini) computeContinuGain(A *[]float64, T *[]string, C *[]string) {
	var a, nleft, nright, partidx int
	var pleft, pright float64
	var gleft, gright float64
	var tleft, tright []string

	nsample := len(*A)

	if glog.V(2) {
		glog.Infoln("sorted data:", A)
		glog.Infoln("Gini.Value:", gini.Value)
	}

	for p := range gini.ContinuPart {

		// find the split of samples between partition based on
		// partition value
		partidx = nsample
		for a = range *A {
			if (*A)[a] > gini.ContinuPart[p] {
				partidx = a
				break;
			}
		}

		nleft = partidx
		nright = nsample - partidx
		pleft = float64(nleft) / float64(nsample)
		pright = float64(nright) / float64(nsample)

		if partidx > 0 {
			tleft = (*T)[0:partidx]
			tright = (*T)[partidx:]

			gleft = gini.compute(&tleft, C)
			gright = gini.compute(&tright, C)
		} else {
			tleft = nil
			tright = (*T)[0:]

			gleft = 0
			gright = gini.compute(&tright, C)
		}

		// count class in partition
		gini.Index[p] = ((pleft * gleft) + (pright * gright))
		gini.Gain[p] = gini.Value - gini.Index[p]

		if glog.V(2) {
			glog.Infoln(tleft)
			glog.Infoln(tright)

			glog.Infof("GiniGain(%v) = %f - (%f * %f) + (%f * %f) = %f\n",
				gini.ContinuPart[p], gini.Value, pleft, gleft,
				pright, gright, gini.Gain[p])
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

/*
GetMaxPartGainValue return the partition that have the maximum Gini gain.
*/
func (gini *Gini) GetMaxPartGainValue() interface{} {
	if gini.IsContinu {
		return gini.ContinuPart[gini.MaxPartGain]
	}

	return gini.DiscretePart[gini.MaxPartGain]
}

/*
GetMaxGainValue return the value of partition which contain the maximum Gini
gain.
*/
func (gini *Gini) GetMaxGainValue() float64 {
	return gini.MaxGainValue
}

/*
GetMinIndexPartValue return the partition that have the minimum Gini index.
*/
func (gini *Gini) GetMinIndexPartValue() interface{} {
	if gini.IsContinu {
		return gini.ContinuPart[gini.MinIndexPart]
	}

	return gini.DiscretePart[gini.MinIndexPart]
}

/*
GetMinIndexValue return the minimum Gini index value.
*/
func (gini *Gini) GetMinIndexValue() float64 {
	return gini.MinIndexValue
}

/*
FindMaxGain find the attribute and value that have the maximum gain.
The returned value is index of attribute.
*/
func FindMaxGain(gains *[]Gini) (MaxGainIdx int) {
	var gainValue = 0.0
	var maxGainValue = 0.0

	for i := range *gains {
		if (*gains)[i].Skip {
			continue
		}
		gainValue = (*gains)[i].GetMaxGainValue()
		if gainValue > maxGainValue {
			maxGainValue = gainValue
			MaxGainIdx = i
		}
	}

	return
}

/*
FindMinIndex return the index of attribute that have the minimum Gini index.
*/
func FindMinGiniIndex(ginis *[]Gini) (MinIndexIdx int) {
	var indexV = 0.0
	var minIndexV = 1.0

	for i := range *ginis {
		indexV = (*ginis)[i].GetMinIndexValue()
		if indexV > minIndexV {
			minIndexV = indexV
			MinIndexIdx = i
		}
	}

	return
}

/*
String yes, it will print it JSON like format.
*/
func (gini Gini) String() (s string) {
	s = fmt.Sprint("{\n",
	"  Skip          :", gini.Skip, "\n",
	"  IsContinu     :", gini.IsContinu, "\n",
	"  Index         :", gini.Index, "\n",
	"  Value         :", gini.Value, "\n",
	"  Gain          :", gini.Gain, "\n",
	"  MaxPartGain   :", gini.MaxPartGain, "\n",
	"  MaxGainValue  :", gini.MaxGainValue, "\n",
	"  MinIndexPart  :", gini.MinIndexPart, "\n",
	"  MinIndexValue :", gini.MinIndexValue, "\n",
	"  SortedIndex   :", gini.SortedIndex, "\n",
	"  ContinuPart   :", gini.ContinuPart, "\n",
	"  DiscretePart  :", gini.DiscretePart, "\n",
	"}")
	return
}
