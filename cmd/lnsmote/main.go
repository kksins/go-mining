// Copyright 2016 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/shuLhan/dsv"
	"github.com/shuLhan/go-mining/resampling/lnsmote"
	"github.com/shuLhan/tabula"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

var (
	// DEBUG level, can be set from environment variable.
	DEBUG = 0
	// percentOver contain percentage of over sampling.
	percentOver = 100
	// knn contain number of nearest neighbours considered when
	// oversampling.
	knn = 5
	// synFile flag for synthetic file output.
	synFile = ""
	// merge flag, if its true the original and synthetic will be merged
	// into `synFile`.
	merge = false
)

var usage = func() {
	cmd := os.Args[0]
	fmt.Fprintf(os.Stderr, "Usage of %s:\n"+
		"[-percentover number] "+
		"[-knn number] "+
		"[-syntheticfile string] "+
		"[-merge bool] "+
		"[config.dsv]\n", cmd)
	flag.PrintDefaults()
}

func init() {
	var e error

	v := os.Getenv("DEBUG")
	DEBUG, e = strconv.Atoi(v)
	if e != nil {
		DEBUG = 0
	}

	flagUsage := []string{
		"Percentage of oversampling (default 100)",
		"Number of nearest neighbours (default 5)",
		"File where synthetic samples will be written (default '')",
		"If true then original and synthetic will be merged when" +
			" written to file (default false)",
	}

	flag.IntVar(&percentOver, "percentover", -1, flagUsage[0])
	flag.IntVar(&knn, "knn", -1, flagUsage[1])
	flag.StringVar(&synFile, "syntheticfile", "", flagUsage[2])
	flag.BoolVar(&merge, "merge", false, flagUsage[3])
}

func trace(s string) (string, time.Time) {
	fmt.Println("[START]", s)
	return s, time.Now()
}

func un(s string, startTime time.Time) {
	endTime := time.Now()
	fmt.Println("[END]", s, "with elapsed time",
		endTime.Sub(startTime))
}

//
// createLnsmote will create and initialize SMOTE object from config file and
// from command parameter.
//
func createLnsmote(fcfg string) (lnsmoteRun *lnsmote.Runtime, e error) {
	lnsmoteRun = &lnsmote.Runtime{}

	config, e := ioutil.ReadFile(fcfg)
	if e != nil {
		return nil, e
	}

	e = json.Unmarshal(config, lnsmoteRun)
	if e != nil {
		return nil, e
	}

	// Use option value from command parameter.
	if percentOver > 0 {
		lnsmoteRun.PercentOver = percentOver
	}
	if knn > 0 {
		lnsmoteRun.K = knn
	}

	if DEBUG >= 1 {
		fmt.Println("[lnsmote]", lnsmoteRun)
	}

	return
}

//
// runLnsmote will select minority class from dataset and run oversampling.
//
func runLnsmote(lnsmoteRun *lnsmote.Runtime, dataset *tabula.Claset) (e error) {
	e = lnsmoteRun.Resampling(dataset)
	if e != nil {
		return
	}

	if DEBUG >= 1 {
		fmt.Println("[lnsmote] # synthetics:",
			lnsmoteRun.GetSynthetics().Len())
	}

	return
}

// runMerge will append original dataset to synthetic file.
func runMerge(lnsmoteRun *lnsmote.Runtime, dataset *tabula.Claset) (e error) {
	writer, e := dsv.NewWriter("")
	if e != nil {
		return
	}

	e = writer.ReopenOutput(lnsmoteRun.SyntheticFile)
	if e != nil {
		return
	}

	sep := dsv.DefSeparator
	n, e := writer.WriteRawDataset(dataset, &sep)
	if e != nil {
		return
	}

	if DEBUG >= 1 {
		fmt.Println("[lnsmote] # appended:", n)
	}

	return writer.Close()
}

func main() {
	defer un(trace("lnsmote"))

	flag.Parse()

	if len(flag.Args()) <= 0 {
		usage()
		os.Exit(1)
	}

	fcfg := flag.Arg(0)

	// Parsing config file and parameter.
	lnsmoteRun, e := createLnsmote(fcfg)
	if e != nil {
		panic(e)
	}

	// Get dataset.
	dataset := tabula.Claset{}
	_, e = dsv.SimpleRead(fcfg, &dataset)
	if e != nil {
		panic(e)
	}

	fmt.Println("[lnsmote] Dataset:", &dataset)

	row := dataset.GetRow(0)
	fmt.Println("[lnsmote] sample:", row)

	e = runLnsmote(lnsmoteRun, &dataset)
	if e != nil {
		panic(e)
	}

	if !merge {
		return
	}

	e = runMerge(lnsmoteRun, &dataset)
	if e != nil {
		panic(e)
	}
}
