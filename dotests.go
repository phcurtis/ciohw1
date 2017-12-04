// Copyright 2017 phcurtis ciohw1 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// This file contains the main "do tests" loop for processing a testsMetricSet.

import (
	"log"
	"sort"

	"github.com/phcurtis/fn"
)

// TBD: yet to decide if make this func public and related doc maybe some refactor.
func doTests(pfix string, tms *testsMetricSet, verblvl int, fltrace bool) (testsRun, errCnt int) {
	defer fn.LogCondTrace(fltrace)() //may want to activate also on certain -verblvl

	failed := "FAILED"
	if verblvl > 0 {
		log.Print(pfix + "Beg:" + tms.desc)
	}
	indent := "  "
	for _, v := range tms.test {
		testresult := "PASSED"
		name := pfix + indent + tms.name + ":Test:" + v.name
		if verblvl > 1 {
			log.Print(indent + name + ":" + "processing ...\n")
		}
		var err Meterror
		switch v.ftype {
		case fMObjColItegChk:
			f := v.fn1.(func(url1 string, reset bool, loud int) Meterror)
			arg1 := v.arg1.(string)
			arg2 := v.arg2.(bool)
			arg3 := v.arg3.(int)
			if err = f(arg1, arg2, arg3); MetErrCodeDiff(err, v.res) {
				log.Printf(indent+"%s:%v", name, err)
				testresult = failed
			}
		case fGetMObjCol:
			f := v.fn1.(func(url1 string, dump bool, pfix string, trace bool) (*MObjCol, Meterror))
			arg1 := v.arg1.(string)
			arg2 := v.arg2.(bool)
			arg3 := v.arg3.(string)
			arg4 := v.arg4.(bool)
			if _, err = f(arg1, arg2, arg3, arg4); MetErrCodeDiff(err, v.res) {
				log.Printf(indent+"%s:%v", name, err)
				testresult = failed
			}
		case fPostEusrMetrics, fGetEusrMetrics:
			f := v.fn1.(func(url1, metrics string) (*EusrMetrics, Meterror))
			arg1 := v.arg1.(string)
			arg2 := v.arg2.(string)
			email := v.arg3.(string)
			vmetrics := v.arg4.([]string)
			sort.Strings(vmetrics)
			if verblvl > 1 {
				log.Printf(indent+"%s:Eusr:%q\n", name, email)
				log.Printf(indent+"%s:Metrics:%v\n", name, vmetrics)
			}

			var em *EusrMetrics
			if em, err = f(arg1, arg2); MetErrCodeDiff(err, v.res) {
				if err != nil {
					log.Printf(indent+"%s:%v", name, err)
				} else {
					log.Printf(indent+"%s:goterr:%v wantErrCode:%d(%v)\n",
						name, err, MetError{Code1: v.res}.Code(), MetError{Code1: v.res}.CodeText())
				}
				testresult = failed
			} else {
				if err != nil {
					if verblvl > 1 {
						if err.Code() == v.res {
							log.Printf(indent+"%s:expected+got:ErrCode:%d(%v)\n",
								name, err.Code(), err.CodeText())
						}
					}
				} else {
					if verblvl > 1 {
						log.Printf(indent+"%s:got:Eusr:%q\n", name, em.Email)
						log.Printf(indent+"%s:got:Metrics:%v\n", name, em.Metrics)
					}
					err = CompMetrics(em, email, vmetrics)
					if err != nil {
						log.Printf(indent+"%s:%v", name, err)
						testresult = failed
					}
				}
			}
		case fVerifyEusr:
			f := v.fn1.(func(url1, email string, metrics []string) Meterror)
			arg1 := v.arg1.(string)
			email := v.arg2.(string)
			vmetrics := v.arg3.([]string)
			sort.Strings(vmetrics)
			if verblvl > 1 {
				log.Printf(indent+"%s:Eusr:%q\n", name, email)
				log.Printf(indent+"%s:Metrics:%v\n", name, vmetrics)
			}
			if err = f(arg1, email, vmetrics); MetErrCodeDiff(err, v.res) {
				// TODO: MeusrMetricsWrong may want adjust output relating to this in future
				// TODO: MeusrDiffThanURL  may want adjust output relating to this in future
				log.Printf(indent+"%s:%v", name, err)
				testresult = failed
			}
		case fPutEusrMetric:
			f := v.fn1.(func(url1, email, metric string) (*EusrMetrics, Meterror))
			arg1 := v.arg1.(string)
			email := v.arg2.(string)
			metric := v.arg3.(string)
			vmetrics := v.arg4.([]string)
			if verblvl > 1 {
				log.Printf(indent+"%s:put:email:%q\n", name, email)
				log.Printf(indent+"%s:put:Metric:%v\n", name, metric)
			}

			var em *EusrMetrics
			if em, err = f(arg1, email, metric); MetErrCodeDiff(err, v.res) {
				if err != nil {
					log.Printf(indent+"%s:%v", name, err)
				} else {
					log.Printf(indent+"%s:goterr:%v wantErrCode:%d(%v)\n",
						name, err, MetError{Code1: v.res}.Code(), MetError{Code1: v.res}.CodeText())
				}
				testresult = failed
			} else {
				if err != nil {
					if verblvl > 1 {
						if err.Code() == v.res {
							log.Printf(indent+"%s:expected+got:ErrCode:%d(%v)\n",
								name, err.Code(), err.CodeText())
						}
					}
				} else {
					if verblvl > 1 {
						log.Printf(indent+"%s:got:email:%q\n", name, em.Email)
						log.Printf(indent+"%s:got:Metrics:%v\n", name, em.Metrics)
					}
					err = CompMetrics(em, email, vmetrics)
					if err != nil {
						log.Printf(indent+"%s:%v", name, err)
						testresult = failed
					}
				}
			}
		case fDelEusrMetricVerify:
			f := v.fn1.(func(url1, email, metric string, metrics []string, verblvl int) Meterror)
			arg1 := v.arg1.(string)
			email := v.arg2.(string)
			metric := v.arg3.(string)
			vmetrics := v.arg4.([]string)
			verblvl := v.arg5.(int)
			if verblvl > 1 {
				log.Printf(indent+"%s:email:%q\n", name, email)
				log.Printf(indent+"%s:del:Metric:%v\n", name, metric)
				log.Printf(indent+"%s:want:Metrics:%v\n", name, vmetrics)
			}
			if err = f(arg1, email, metric, vmetrics, verblvl); MetErrCodeDiff(err, v.res) {
				log.Printf(indent+"%s:%v", name, err)
				testresult = failed
			}
		default:
			log.Panicf("%v NOT IMPLEMENTED ftype:%d", name, v.ftype)
		}

		testsRun++
		if testresult == failed {
			errCnt++
		}
		if verblvl > 0 || testresult == failed {
			log.Printf(indent+"%s:%s\n", name, testresult)
		}
	}
	if verblvl > 0 {
		res := "PASSED"
		if errCnt > 0 {
			res = "FAILED"
		}
		log.Printf(pfix+"End:%s %s", tms.desc, res)
	}
	return testsRun, errCnt
}
