// Copyright 2017 phcurtis ciohw1 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// This file contains api auxilary related functions an user
// of 'metrics' api may find helpful when verifying its compliance to its spec.

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"

	"github.com/phcurtis/fn"
)

// MetErrCodeDiff - returns true if Meterror.Code is different than "wantCode".
func MetErrCodeDiff(got Meterror, wantCode MetErrCode) bool {
	if got != nil {
		return got.Code() != wantCode
	}
	return wantCode != MnoError
}

// EusrAddrDups - returns any email-user address duplicates list
// if any found in "moc".
func (moc *MObjCol) EusrAddrDups(trace bool) (dups []string) {
	defer fn.LogCondTrace(trace)()
	emailmap := make(map[string]int, 1)
	for i := 0; i < len(moc.Grp); i++ {
		email := moc.Grp[i].Email
		emailmap[email]++
	}
	for k := range emailmap {
		if emailmap[k] > 1 {
			dups = append(dups, fmt.Sprintf("%q(%d)", k, emailmap[k]))
		}
	}
	return dups
}

// MetricsDups - returns a list of any email-user metrics duplicates
// found in "metrics".
func MetricsDups(metrics []string) (dups []string) {
	metricmap := make(map[string]int, 1)
	for i := 0; i < len(metrics); i++ {
		metricmap[metrics[i]]++
	}
	for i := range metricmap {
		if metricmap[i] > 1 {
			dups = append(dups, fmt.Sprintf("%q(%d)", i, metricmap[i]))
		}
	}
	return dups
}

// MetricsValid - returns a non-nil error if any invalid email-user
// metrics are found in "metrics".
func MetricsValid(metrics []string) error {
	for _, v := range metrics {
		if !metricPossibleMap[v] {
			return fmt.Errorf("metric-INVALID:%q", v)
		}
	}
	return nil
}

// EusrMetricsValidate - returns a list of any metrics problems with an
// email-user in "moc", such as duplicates or an invalid metric type.
func (moc *MObjCol) EusrMetricsValidate(trace bool) (problems []string) {
	defer fn.LogCondTrace(trace)()
	for i := 0; i < len(moc.Grp); i++ {
		if err := MetricsValid(moc.Grp[i].Metrics); err != nil {
			problems = append(problems, fmt.Sprintf("%v", err))
		}
		if dups := MetricsDups(moc.Grp[i].Metrics); len(dups) > 0 {
			problems = append(problems, fmt.Sprintf("dups[%s]:%v", moc.Grp[i].Email, dups))
		}
	}
	return problems
}

// EusrAddrExists - returns whether or not an email-user exists within "moc".
func (moc *MObjCol) EusrAddrExists(email string) bool {
	for i := 0; i < len(moc.Grp); i++ {
		if moc.Grp[i].Email == email {
			return true
		}
	}
	return false
}

// EusrMetricsMatch - returns nil if an email-user metrics match ones within "moc".
func (moc *MObjCol) EusrMetricsMatch(email string, metrics []string) Meterror {
	for i := 0; i < len(moc.Grp); i++ {
		if moc.Grp[i].Email == email {
			got := fmt.Sprintf("%v", moc.Grp[i].Metrics)
			want := fmt.Sprintf("%v", metrics)
			if got != want {
				return MetError{MeusrMetricsDiffInMOC, 0,
					fmt.Errorf("MObjCol-WRONG:Eusr:%s\n  \tgot:%v \n \twant:%v", email, got, want),
					"func:" + fn.CurBase()}
			}
			return nil
		}
	}
	return MetError{MeusrNotFoundInMOC, 0,
		fmt.Errorf("Eusr:%s NOT-FOUND in MObjCol", email),
		"func:" + fn.CurBase()}
}

// CompMetrics ... returns nil if the given EusrMetrics struct matches
// the email-user address and metrics.
func CompMetrics(em *EusrMetrics, email string, metrics []string) Meterror {
	sort.Strings(metrics)
	got := fmt.Sprintf("%v", em.Metrics)
	want := fmt.Sprintf("%v", metrics)
	if got != want {
		return MetError{MeusrMetricsWrong, 0,
			fmt.Errorf("got:%v want:%v", got, want), fn.CurBase()}
	}
	if em.Email != email {
		return MetError{MeusrDiffThanURL, 0,
			fmt.Errorf("got:%v want:%v", em.Email, email), fn.CurBase()}
	}
	return nil
}

// VerifyEusr ... returns nil if a email-user is found in
// the metrics obj collection (found at the specified "url") and its "metrics" match.
func VerifyEusr(url1, email string, metrics []string) Meterror {
	em, err := GetEusrMetrics(url1, email)
	if err != nil {
		return err
	}
	return CompMetrics(em, email, metrics)
}

// DelEusrMetricVerify ... returns nil if can successfully delete
// an email-user's specified metric and the resulting email-user metrics
// match "metrics".
func DelEusrMetricVerify(url1, email, metric string, metrics []string, verblvl int) Meterror {
	defer fn.LogCondTrace(verblvl > 2)()
	err := DeleteEusrMetric(url1, email, metric)
	if err != nil {
		return err
	}
	if verblvl > 2 {
		log.Printf("  returned from DeleteEusrMetric:err:%v\n", err)
	}
	em, err := GetEusrMetrics(url1, email)
	if err != nil {
		return err
	}
	return CompMetrics(em, email, metrics)
}

// VerifyEusrInMo ... returns nil if a email-user is found in
// the "moc" and its "metrics" match.
func (moc *MObjCol) VerifyEusrInMo(email string, metrics []string) (err Meterror) {
	if !moc.EusrAddrExists(email) {
		return MetError{MeusrNotFoundInMOC, 0,
			fmt.Errorf("Eusr:%q NOT-FOUND", email), fn.CurBase()}
	}
	if err = moc.EusrMetricsMatch(email, metrics); err != nil {
		return err
	}
	return nil
}

// DumpMObjCol - dumps via MarshalIndent the "moc" with 'pfix'.
func (moc *MObjCol) DumpMObjCol(pfix string, trace bool) Meterror {
	defer fn.LogCondTrace(trace)()
	out, jerr := json.MarshalIndent(moc, pfix, "\t")
	if jerr != nil {
		return MetError{MerrJSONmarshalIndent, 0, jerr, "func:" + fn.CurBase()}
	}
	log.Printf("MObjColItems:%d\n%s%s\n", len(moc.Grp), pfix, string(out))
	return nil
}

// MObjColItegChk - returns nil if metrics obj collection passes a series of
// validity or integrity checks. These includes the following:
//	- checking that there are no email-user address duplicates
//  - checking that email-user metrics are valid (no dups or unknown).
// reset - if true resets the metrics obj collection before checks begin;
// verblvl - enables some level of verbosity during the above process.
func MObjColItegChk(url1 string, reset bool, verblvl int) Meterror {
	defer fn.LogCondTrace(verblvl > 0)()
	if reset {
		err := ResetMetrics(url1, verblvl > 0)
		if err != nil {
			return err
		}
	}
	moc, err := GetMObjCol(url1, verblvl > 1, fn.CurBase()+":", verblvl > 0)
	if err != nil {
		return err
	}
	if dups := moc.EusrAddrDups(verblvl > 0); len(dups) > 0 {
		return MetError{MeusrDupsFound, 0,
			fmt.Errorf("duplicates-FOUND:%v", dups), fn.CurBase()}
	}

	if problems := moc.EusrMetricsValidate(verblvl > 0); len(problems) > 0 {
		return MetError{MeusrMetricsInvalid, 0,
			fmt.Errorf("problems-FOUND:%v", problems), fn.CurBase()}
	}
	return nil
}
