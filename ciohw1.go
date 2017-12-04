// Copyright 2017 phcurtis ciohw1 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package main - provides a CLI interface for testing a metrics interface
// according to it specification as provided by client.
// To use this program to test all default test cases:
//
// RUN: ./ciohw1 on the command line after you have built the said
// executable, which with no other options will exercise the
// "default test cases".
// [to show that list see -tsnameslist -tsnamelistrunseq options]
// For a multitude of options invoke with option -h for more info.
// One immediate tip increasing X in -verblvl=X option will provide more
// "inside" info including key items and the why of failures.
// To see what tests are running use >0 for X with -verblvl
// i.e. ./ciohw1 -verblvl=1
// If all the requested tests pass then the program exit value will be zero,
// for other possible values ... grep for os.Exit and log.Fatal.
//
// for examples: see file: EXAMPLES
//
// INSTALL: via normal golang setup environment:
// 	go get github.com/phcurtis/ciohw1
// 	go get github.com/phcurtis/fn
// 	Then go build while in github.com/phcurtis/ciohw1
// 	Currently this has only been alpha tested in ubuntu.
package main

// This file contains the cli "main".

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/phcurtis/fn"
)

// Version of this program
const Version = "0.01b"

func main() {
	prelimsCLI()
	onexit := fn.LogCondTrace(flags.mainfnlogcondtrace)
	defer onexit()

	totErrCnt := 0
	totTestsRun := 0

	for i := 0; i < flags.dotestsiterations; i++ {
		pass := make(map[string]int, 1)
		pfix := ""
		if flags.dotestsiterations > 1 {
			pfix = fmt.Sprintf("\titer=%d:", i+1)
		}
		if flags.dotestsiterationsshow {
			log.Printf("Starting iteration:%d timeNow:%v\n", i+1, time.Now())
		}

		runseq := strings.Split(tsnameslistrunseq, ",")
		for _, v := range runseq {
			pass[v]++
			emailbase := fmt.Sprintf("%s-i%dp%d", v, i, pass[v])
			//fmt.Printf("emailbase:%s\n", emailbase)
			setTestsMetricSetFullList(v, flags.baseurl, emailbase,
				flags.verblvl, flags.mocichkverblvl, flags.mocichkreset)
			curTms, err := findTestSetForName(v)
			if err != nil {
				log.Panicf("%v given with -tsnameslistrunseq=%s\n", v, flags.tsnameslistrunseq)
			}
			testsRun, errCnt := doTests(pfix, curTms, flags.verblvl, flags.doteststrace)
			totTestsRun += testsRun
			totErrCnt += errCnt
		}
		if flags.dotestsiterations > 1 {
			// TODO: insert code  ... maybe dumping totals thus far likely need subtot set of counters
		}
	}
	if flags.mocdumpatend {
		_, err := GetMObjCol(flags.baseurl, true, fn.CurBase()+":", flags.verblvl > 2)
		if err != nil {
			log.Printf("error calling GetMObjCol:err:%v", err)
			totErrCnt++
		}
	}
	if totErrCnt != 0 {
		onexit() // have to do here because os.Exit does not allow defers to run.
		log.Fatalf("Exiting: %d of %d Test(s) FAILED", totErrCnt, totTestsRun)
	}

	msg := "PASSED"
	if flags.verblvl > 0 {
		msg = fmt.Sprintf("Exiting: all %d requested Test(s) PASSED", totTestsRun)
	}
	log.Println(msg)
}
