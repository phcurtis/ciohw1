// Copyright 2017 phcurtis ciohw1 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// This file contains cli related hooks to support command line options.

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/phcurtis/fn"
)

// see corresponding init() for flagsStruct variables descriptions
type flagsStruct struct {
	baseurl               string //
	doteststrace          bool   //
	dotestsiterations     int    //
	dotestsiterationsshow bool   //
	logflags              int    //
	logpfix               string //
	fnlogflags            int    //
	fnlogpfix             string //
	mainfnlogcondtrace    bool   //
	mocdumpatend          bool   //
	mocichkverblvl        int    //
	mocichkreset          bool   //
	reqargsdump           bool   //
	//TODO:	resetb4testset bool  //probably better to add reset as a testset
	//TODO:resetaftestset bool   //probably better to add reset as a testset
	resbodydump       bool   //
	tsnameslist       bool   //
	tsnameslistrunseq string //
	showinvokeline    bool   //
	showversion       bool   //
	verblvl           int    //
}

var flags = flagsStruct{}

const defurl = "http://paul-cio-homework.getsandbox.com"

const deflogpfix = "LOG:"
const logCmnFlags = log.Ldate | log.Ltime | log.Lshortfile
const deflogflags = 0 //off

func init() {
	flag.StringVar(&flags.baseurl, "baseurl", defurl, "base url for metrics")
	flag.BoolVar(&flags.doteststrace, "doteststrace", false, "do tests trace using fn.LogCondTrace")
	flag.IntVar(&flags.dotestsiterations, "iterations", 1, "do tests iterations")
	flag.BoolVar(&flags.dotestsiterationsshow, "iterationsshow", false, "show do tests iterations useful with -verblvl=0")
	flag.IntVar(&flags.fnlogflags, "fnlogflags", fn.LogFlags(), "fn LogFlags see fn package for settings")
	flag.StringVar(&flags.fnlogpfix, "fnlogpfix", fn.LogPrefix(), "fn LogPrefix see fn package")

	logflagsHelp := fmt.Sprintf("see stdlib log for settings for 'date time shortfile' use:%d", logCmnFlags)
	flag.IntVar(&flags.logflags, "logflags", deflogflags, logflagsHelp)

	flag.StringVar(&flags.logpfix, "logpfix", deflogpfix, "prefix for logging")
	flag.BoolVar(&flags.mainfnlogcondtrace, "mainfnlogcondtrace", false, "trace main using fn.LogCondTrace")
	flag.BoolVar(&flags.mocdumpatend, "mocdumpatend", false, "dumps MObjCol at program end")
	flag.IntVar(&flags.mocichkverblvl, "mocichkverblvl", 0, "MObjCol integrity verbosity level")
	flag.BoolVar(&flags.mocichkreset, "mocichkreset", false, "reset MObjCol before integrity check")
	flag.BoolVar(&flags.reqargsdump, "reqargsdump", false, "dump req[uest] input arguments")
	flag.BoolVar(&flags.resbodydump, "resbodydump", false, "dump res[ponse] body")
	flag.BoolVar(&flags.tsnameslist, "tsnameslist", false, "list possible test set names and exit")
	flag.StringVar(&flags.tsnameslistrunseq, "tsnameslistrunseq", "", "run specified test set names CSV style if empty means (all); to see all use -tsnameslist")
	flag.BoolVar(&flags.showinvokeline, "invokeline", false, "show invocation line")
	flag.BoolVar(&flags.showversion, "version", false, "show version")
	flag.IntVar(&flags.verblvl, "verblvl", 0, "verbosity level")
}

func prelimsCLI() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s: (Version:%s)\n", os.Args[0], Version)
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "unrecognized %v\nUsage of ./%s (Version:%s):\n",
			flag.Args(), filepath.Base(os.Args[0]), Version)
		flag.PrintDefaults()
		os.Exit(2) // TBD: perhaps would like this to be some other non-zero value
	}

	log.SetFlags(flags.logflags)
	log.SetPrefix(flags.logpfix)
	fn.LogSetFlags(flags.fnlogflags)
	fn.LogSetPrefix(flags.fnlogpfix)

	if flags.showversion || flags.verblvl > 0 {
		log.Printf("./%s version=%s\n", filepath.Base(os.Args[0]), Version)
	}
	if flags.showinvokeline || flags.verblvl > 3 {
		log.Printf("%v\n", os.Args)
	}

	namesHelp := func() {
		_, namesHelp := testSetPosNames()
		log.Println("Possible values for -tsnameslistrunseq follow:")
		for _, v := range namesHelp {
			log.Println(v)
		}
	}

	if flags.tsnameslist {
		namesHelp()
		log.Println("exiting")
		os.Exit(3) // TBD: exit value
	}
	tsnameslistrunseq, _ = testSetPosNames()
	if len(flags.tsnameslistrunseq) > 0 {
		runseq := strings.Split(flags.tsnameslistrunseq, ",")
		for _, v := range runseq {
			if _, err := findTestSetForName(v); err != nil {
				log.Printf("err:%v given with -tsnameslistrunseq=%s\n", err, flags.tsnameslistrunseq)
				namesHelp()
				log.Println("exiting")
				os.Exit(2) // TBD: exit value
			}
		}
		tsnameslistrunseq = flags.tsnameslistrunseq
	}

	SetReqArgsDump(flags.reqargsdump)
	SetResBodyDump(flags.resbodydump)
}
