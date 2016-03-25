// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smote_test

import (
	"fmt"
	"github.com/shuLhan/dsv"
	"github.com/shuLhan/go-mining/knn"
	"github.com/shuLhan/go-mining/resampling/smote"
	"github.com/shuLhan/tabula"
	"testing"
)

var (
	fcfg        = "../../testdata/phoneme/phoneme.dsv"
	PercentOver = 100
	K           = 5
)

func TestSmote(t *testing.T) {
	smot := &smote.Runtime{
		Runtime: knn.Runtime{
			DistanceMethod: knn.TEuclidianDistance,
			ClassIndex:     5,
			K:              K,
		},
		PercentOver: PercentOver,
		Synthetic:   nil,
	}

	// Read samples.
	dataset := tabula.Claset{}

	_, e := dsv.SimpleRead(fcfg, &dataset)
	if nil != e {
		t.Fatal(e)
	}

	fmt.Println("[smote_test] Total samples:", dataset.Len())

	minorset := dataset.GetMinorityRows()

	fmt.Println("[smote_test] # minority samples:", minorset.Len())

	e = smot.Resampling(*minorset)
	if e != nil {
		t.Fatal(e)
	}

	fmt.Println("[smote_test] # synthetic:", smot.Synthetic.Len())

	e = smot.Write("phoneme_smote.csv")
	if e != nil {
		t.Fatal(e)
	}
}
