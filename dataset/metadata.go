// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dataset

import (
	"fmt"
	"github.com/shuLhan/dsv"
)

/*
Metadata extension for common dataset.
*/
type Metadata struct {
	// dsv.Metadata as our base
	dsv.Metadata
	// IsContinu indicated whether the data continu or not.
	IsContinu bool `json:"IsContinu"`
	// NominalValues contain list of known discrete values in data.
	NominalValues []string `json:"NominalValues"`
}

/*
GetNominalValue return the nominal value for discrete attribute.
If attribute is continuous, return nil.
*/
func (md *Metadata) GetNominalValue() []string {
	if md.IsContinu {
		return nil
	}
	return md.NominalValues
}

/*
String return the pretty print format of attribute.
*/
func (md Metadata) String() (s string) {
	s = fmt.Sprintf("{\n\tIsContinue: %v\n", md.IsContinu)

	if !md.IsContinu {
		s += fmt.Sprintf("\tNominalValues: %v\n", md.NominalValues)
	}

	s += "}"

	return
}
