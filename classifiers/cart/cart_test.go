// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cart_test

import (
	"fmt"
	"github.com/shuLhan/go-mining/classifiers/cart"
	"github.com/shuLhan/go-mining/dataset"
	"io"
	"reflect"
	"runtime/debug"
	"testing"
)

const (
	NRows = 150
)

func assert(t *testing.T, exp, got interface{}, equal bool) {
	if reflect.DeepEqual(exp, got) != equal {
		debug.PrintStack()
		t.Fatalf("\n"+
			">>> Expecting '%v'\n"+
			"          got '%v'\n", exp, got)
	}
}

func TestCART(t *testing.T) {
	ds, e := dataset.NewReader("../../testdata/iris/iris.dsv")

	if nil != e {
		t.Fatal(e)
	}

	e = ds.Read()

	if nil != e && e != io.EOF {
		t.Fatal(e)
	}

	assert(t, NRows, ds.GetNRow(), true)

	// Build CART tree.
	CART := cart.New(cart.SplitMethodGini, 0)

	e = CART.Build(ds)

	if e != nil {
		t.Fatal(e)
	}

	fmt.Println("CART Tree:\n", CART)

	// Reread the data
	ds.Reset()
	ds.Open()

	e = ds.Read()
	if nil != e && e != io.EOF {
		t.Fatal(e)
	}

	// Create test set
	testset, e := dataset.NewReader("../../testdata/iris/iris.dsv")

	if nil != e {
		t.Fatal(e)
	}

	e = testset.Read()

	if nil != e && e != io.EOF {
		t.Fatal(e)
	}

	testset.GetTarget().ClearValues()

	// Classifiy test set
	CART.ClassifySet(testset)

	assert(t, ds.GetTarget(), testset.GetTarget(), true)
}
