// Copyright 2017 phcurtis ciohw1 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// This file contains the "testlists" and related hooks.

import (
	"fmt"
	"strings"
	"time"

	"github.com/phcurtis/fn"
)

var tsnameslistrunseq string // if not set than default to run all in default sequence

type testArgs struct {
	name  string
	fn1   interface{}
	ftype int
	res   MetErrCode
	arg1  interface{}
	arg2  interface{}
	arg3  interface{}
	arg4  interface{}
	arg5  interface{}
}

type testsMetricSet struct {
	tsid string // test set id currently not used but may prove useful.
	name string // should limit to 11 or less Alphanumeric and NO spaces
	// ... and assumed and needs to be unique and non-blank, especially useful for
	// ... clarification when verblvl=0 and there are error(s).
	desc string // heading for tests ... displayed when verblvl is > 0.
	test []testArgs
}

// these constants are used with indicating what function
// is to be called; its ftype in testArgs struct.
const (
	_ = iota
	fDelEusrMetricVerify
	fGetEusrMetrics
	fGetMObjCol
	fMObjColItegChk
	fPostEusrMetrics
	fPutEusrMetric
	fVerifyEusr
)

func genUniqTestSetIDInit(basename string) func() string {
	cnt := 0
	bname := basename
	return func() string {
		cnt++
		return fmt.Sprintf("%s%02d", bname, cnt)
	}
}

var gutsid = genUniqTestSetIDInit("tsid")

var gm = func(name string, num int) string {
	return fmt.Sprintf("e-%s-%d@metrics.com", name, num)
}

func createTsMapEmail(name string) string {
	if len(name) < 1 || len(name) > 11 {
		panic(fmt.Sprintf("name too long %q\n", name))
	}
	name = strings.ToLower(name)
	// TODO: later check that is alphanumeric for now when define test make sure they are
	if _, ok := tsMapEmail[name]; !ok {
		tsMapEmail[name] = ""
	}
	return name
}

var mocintegchkTS = func(url1 string, reset bool, verblvl int) (string, testsMetricSet) {
	name := createTsMapEmail("mocintegchk")
	return name, testsMetricSet{
		gutsid(),
		name,
		"Verifying: Integrity (a sanity check) of the Metrics Object Collection",
		[]testArgs{
			testArgs{"T01-MObjColItegChk............", MObjColItegChk, fMObjColItegChk, MnoError,
				url1, reset, verblvl, "", ""},
		},
	}
}

var apibasicsTS = func(url1, email string, verblvl int) (string, testsMetricSet) {
	name := createTsMapEmail("apibasics")
	return name, testsMetricSet{
		gutsid(),
		name,
		"Verifying: can successfully exercise all Metrics API basics",
		[]testArgs{
			testArgs{"T01-GetMObjCol................", GetMObjCol, fGetMObjCol, MnoError,
				url1, verblvl > 3, "", verblvl > 3, ""},

			testArgs{"T02a-PostEusrMetrics-noMetrics.", PostEusrMetrics, fPostEusrMetrics, MnoError,
				url1, GenMetricsJSONStr(gm(email, 1), []string{}), gm(email, 1), []string{}, ""},

			testArgs{"T02b-GetEusrMetrics-verify.....", VerifyEusr, fVerifyEusr, MnoError,
				url1, gm(email, 1), []string{}, "", ""},

			testArgs{"T02c-PutEusrMetric-delivered...", PutEusrMetric, fPutEusrMetric, MnoError,
				url1, gm(email, 1), Gen1MetricJSONStr(gm(email, 1), Mdelivered), []string{Mdelivered}, ""},

			testArgs{"T02d-DeleteEusrMetric-verify...", DelEusrMetricVerify, fDelEusrMetricVerify, MnoError,
				url1, gm(email, 1), Gen1MetricJSONStr(gm(email, 1), Mdelivered), []string{}, verblvl},
		},
	}
}

var dupusersTS = func(url1, email string) (string, testsMetricSet) {
	name := createTsMapEmail("dupusers")
	return name, testsMetricSet{
		gutsid(),
		name,
		"Verifying: duplicate users (by email address) not allowed",
		[]testArgs{
			testArgs{"T01a-PostEusrMetrics-noMetrics.........", PostEusrMetrics, fPostEusrMetrics, MnoError,
				url1, GenMetricsJSONStr(gm(email, 1), []string{}), gm(email, 1), []string{}, ""},

			testArgs{"T01b-PostEusrMetrics-noMetrics-dupuser.", PostEusrMetrics, fPostEusrMetrics, MerrScodeBadRequest,
				url1, GenMetricsJSONStr(gm(email, 1), []string{}), gm(email, 1), []string{}, ""},
		},
	}
}

var dupmetricsTS = func(url1, email string) (string, testsMetricSet) {
	name := createTsMapEmail("dupmetrics")
	return name, testsMetricSet{
		gutsid(),
		name,
		"Verifying: duplicate metrics not allowed",
		[]testArgs{
			testArgs{"T01a-createUser-noMetrics......", PostEusrMetrics, fPostEusrMetrics, MnoError,
				url1, GenMetricsJSONStr(gm(email, 1), []string{}), gm(email, 1), []string{}, ""},

			testArgs{"T01b-Put-delivered.............", PutEusrMetric, fPutEusrMetric, MnoError,
				url1, gm(email, 1), Gen1MetricJSONStr(gm(email, 1), Mdelivered), []string{Mdelivered}, ""},

			testArgs{"T01c-Put-delivered again(dup)..", PutEusrMetric, fPutEusrMetric, MerrScodeBadRequest,
				url1, gm(email, 1), Gen1MetricJSONStr(gm(email, 1), Mdelivered), []string{Mdelivered}, ""},

			testArgs{"T02-createUser-spammed-twice...", PostEusrMetrics, fPostEusrMetrics, MerrScodeBadRequest,
				url1, GenMetricsJSONStr(gm(email, 2), []string{Mspammed, Mspammed}), gm(email, 2), []string{Mspammed}, ""},

			testArgs{"T03a-createUser-withAllMetrics..", PostEusrMetrics, fPostEusrMetrics, MnoError,
				url1, GenMetricsJSONStr(gm(email, 3), MetricsPossible[:]), gm(email, 3), MetricsPossible[:], ""},

			testArgs{"T03b-Put-bounced again(dup).....", PutEusrMetric, fPutEusrMetric, MerrScodeBadRequest,
				url1, gm(email, 3), Gen1MetricJSONStr(gm(email, 3), Mbounced), MetricsPossible[:], ""},

			testArgs{"T03c-Put-clicked again(dup).....", PutEusrMetric, fPutEusrMetric, MerrScodeBadRequest,
				url1, gm(email, 3), Gen1MetricJSONStr(gm(email, 3), Mclicked), MetricsPossible[:], ""},

			testArgs{"T03d-Put-delivered again(dup)...", PutEusrMetric, fPutEusrMetric, MerrScodeBadRequest,
				url1, gm(email, 3), Gen1MetricJSONStr(gm(email, 3), Mdelivered), MetricsPossible[:], ""},

			testArgs{"T03e-Put-opened again(dup)...", PutEusrMetric, fPutEusrMetric, MerrScodeBadRequest,
				url1, gm(email, 3), Gen1MetricJSONStr(gm(email, 3), Mopened), MetricsPossible[:], ""},

			testArgs{"T03f-Put-sent again(dup).....", PutEusrMetric, fPutEusrMetric, MerrScodeBadRequest,
				url1, gm(email, 3), Gen1MetricJSONStr(gm(email, 3), Msent), MetricsPossible[:], ""},

			testArgs{"T03g-Put-spammed again(dup)...", PutEusrMetric, fPutEusrMetric, MerrScodeBadRequest,
				url1, gm(email, 3), Gen1MetricJSONStr(gm(email, 3), Mspammed), MetricsPossible[:], ""},

			testArgs{"T03e-Put-suppressed again(dup)", PutEusrMetric, fPutEusrMetric, MerrScodeBadRequest,
				url1, gm(email, 3), Gen1MetricJSONStr(gm(email, 3), Msuppressed), MetricsPossible[:], ""},
		},
	}
}

var invalidmetTS = func(url1, email string) (string, testsMetricSet) {
	name := createTsMapEmail("invalidmet")
	return name, testsMetricSet{
		gutsid(),
		name,
		"Verifying: invalid metrics not allowed",
		[]testArgs{
			testArgs{"T01a-createUser-noMetrics.........", PostEusrMetrics, fPostEusrMetrics, MnoError,
				url1, GenMetricsJSONStr(gm(email, 1), []string{}), gm(email, 1), []string{}, ""},

			testArgs{"T01b-Put-invalidmet...............", PutEusrMetric, fPutEusrMetric, MerrScodeBadRequest,
				url1, gm(email, 1), Gen1MetricJSONStr(gm(email, 1), "invalidmet"), []string{}, ""},

			testArgs{"T01c-Put-invalidmet...............", PutEusrMetric, fPutEusrMetric, MerrScodeBadRequest,
				url1, gm(email, 1), Gen1MetricJSONStr(gm(email, 1), "Spammed"), []string{"Spammed"}, ""},

			testArgs{"T01d-GetEusrMetrics-shouldHaveNone", VerifyEusr, fVerifyEusr, MnoError,
				url1, gm(email, 1), []string{}, "", ""},
		},
	}
}

var delmetricsTS = func(url1, email string, verblvl int) (string, testsMetricSet) {
	name := createTsMapEmail("delmetrics")
	return name, testsMetricSet{
		gutsid(),
		name,
		"Verifying: delete tracked metrics do as such",
		[]testArgs{
			testArgs{"T01a-createUser-bounced-spammed.", PostEusrMetrics, fPostEusrMetrics, MnoError,
				url1, GenMetricsJSONStr(gm(email, 1), []string{Mbounced, Mspammed}), gm(email, 1), []string{Mbounced, Mspammed}, ""},

			testArgs{"T01b-DeleteMetric-bounced.......", DelEusrMetricVerify, fDelEusrMetricVerify, MnoError,
				url1, gm(email, 1), Gen1MetricJSONStr(gm(email, 1), Mbounced), []string{Mspammed}, verblvl},

			testArgs{"T01c-DeleteMetric-spammed.......", DelEusrMetricVerify, fDelEusrMetricVerify, MnoError,
				url1, gm(email, 1), Gen1MetricJSONStr(gm(email, 1), Mspammed), []string{}, verblvl},
		},
	}
}

var metallusrsTS = func(url1, email string) (string, testsMetricSet) {
	name := createTsMapEmail("metallusrs")
	return name, testsMetricSet{
		gutsid(),
		name,
		"Verifying: that adding 8 users with diff metrics and as such so exist in Metrics Object Col ...",
		[]testArgs{
			testArgs{"T01-MObjColItegChk.with reset..", MObjColItegChk, fMObjColItegChk, MnoError,
				url1, true, 0, "", ""},

			testArgs{"T02a-createUser-noMetrics......", PostEusrMetrics, fPostEusrMetrics, MnoError,
				url1, GenMetricsJSONStr(gm(email, 1), []string{}), gm(email, 1), []string{}, ""},

			testArgs{"T02b-createUser-bounced...", PostEusrMetrics, fPostEusrMetrics, MnoError,
				url1, GenMetricsJSONStr(gm(email, 2), []string{Mbounced}), gm(email, 2), []string{Mbounced}, ""},

			testArgs{"T02c-createUser-clicked...", PostEusrMetrics, fPostEusrMetrics, MnoError,
				url1, GenMetricsJSONStr(gm(email, 3), []string{Mclicked}), gm(email, 3), []string{Mclicked}, ""},

			testArgs{"T02d-createUser-delivered.", PostEusrMetrics, fPostEusrMetrics, MnoError,
				url1, GenMetricsJSONStr(gm(email, 4), []string{Mdelivered}), gm(email, 4), []string{Mdelivered}, ""},

			testArgs{"T02e-createUser-opened....", PostEusrMetrics, fPostEusrMetrics, MnoError,
				url1, GenMetricsJSONStr(gm(email, 5), []string{Mopened}), gm(email, 5), []string{Mopened}, ""},

			testArgs{"T02f-createUser-sent......", PostEusrMetrics, fPostEusrMetrics, MnoError,
				url1, GenMetricsJSONStr(gm(email, 6), []string{Msent}), gm(email, 6), []string{Msent}, ""},

			testArgs{"T02g-createUser-spammed...", PostEusrMetrics, fPostEusrMetrics, MnoError,
				url1, GenMetricsJSONStr(gm(email, 7), []string{Mspammed}), gm(email, 7), []string{Mspammed}, ""},

			testArgs{"T02h-createUser-suppressed", PostEusrMetrics, fPostEusrMetrics, MnoError,
				url1, GenMetricsJSONStr(gm(email, 8), []string{Msuppressed}), gm(email, 8), []string{Msuppressed}, ""},
		},
	}
}
var testsMetricSetFullList []testsMetricSet
var tsfuncl = make(map[string]interface{}, 1)

var tsMapEmail = make(map[string]string, 1)
var tsEmailGroup = 0

func getyymmddhhmmss() string {
	t1 := time.Now()
	tbase := fmt.Sprintf("%02d%02d%02d%02d%02d%02d",
		t1.Year()-2000, t1.Month(), t1.Day(), t1.Hour(), t1.Minute(), t1.Second())
	return tbase
}

func rtnMset(name string, tms testsMetricSet) testsMetricSet {
	tsfuncl[name] = tms // left here not currently implemented may be valuable for refactoring
	return tms
}

func init() {
	tsname := "" // must be "" for initing
	setTestsMetricSetFullList(tsname, "", "", 0, 0, false)
}

// TODO: should be refactored and only reinit specific test since email typically is
// designed to indicate by iteration of do tests and which pass of specific test.
// This needs to be called with tsname = "" during func init()
func setTestsMetricSetFullList(tsname, url1, email string, verblvl, mociVerblvl int, mociReset bool) {
	if tsname == "" {
		email = ""
	} else {
		if _, ok := tsMapEmail[tsname]; !ok {
			panic(fmt.Sprintf("unknown tsname=%q func:%s", tsname, fn.CurBase()))
		}
		tsEmailGroup++
		tsMapEmail[tsname] = fmt.Sprintf("%s-g%d-%s", getyymmddhhmmss(), tsEmailGroup, email)
	}

	testsMetricSetFullList = []testsMetricSet{
		rtnMset(mocintegchkTS(url1, mociReset, mociVerblvl)),
		rtnMset(apibasicsTS(url1, tsMapEmail["apibasics"], verblvl)),
		rtnMset(dupusersTS(url1, tsMapEmail["dupusers"])),
		rtnMset(dupmetricsTS(url1, tsMapEmail["dupmetrics"])),
		rtnMset(invalidmetTS(url1, tsMapEmail["invalidmet"])),
		rtnMset(delmetricsTS(url1, tsMapEmail["delmetrics"], verblvl)),
		rtnMset(metallusrsTS(url1, tsMapEmail["metallusrs"])),
	}
	return
}

func findTestSetForName(name string) (*testsMetricSet, error) {
	for i := range testsMetricSetFullList {
		if name == testsMetricSetFullList[i].name {
			return &testsMetricSetFullList[i], nil
		}
	}
	return nil, fmt.Errorf("%s:%q NOT-FOUND", fn.CurBase(), name)
}

func testSetPosNames() (names string, namesHelp []string) {
	sep := ""
	for i := range testsMetricSetFullList {
		v := testsMetricSetFullList[i]
		namesHelp = append(namesHelp, fmt.Sprintf("testsetname:%s \tdesc:%q", v.name, v.desc))
		names += sep + v.name
		sep = ","
	}
	namesHelp = append(namesHelp, fmt.Sprintf("example: -tsnameslistrunseq:%s", names))
	namesHelp = append(namesHelp, "one can also list names multiple times; "+
		"it is recommended when using this option to also use -invokeline option")

	return names, namesHelp
}
