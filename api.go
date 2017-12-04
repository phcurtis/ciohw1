// Copyright 2017 phcurtis ciohw1 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// This file contains core api stuff and core "helper" supporting funcs.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/phcurtis/fn"
)

// EusrMetrics - email user metrics struct.
type EusrMetrics struct {
	// ID      float64  `json:"id"` //commented out since don't need nor have spec to know its type for sure
	Email   string   `json:"email"`
	Metrics []string `json:"metrics"`
}

// MObjCol - metrics obj collection struct.
type MObjCol struct {
	Grp []EusrMetrics `json:"metrics"`
}

// possible individual email metrics.
const (
	Mbounced    = "bounced"
	Mclicked    = "clicked"
	Mdelivered  = "delivered"
	Mopened     = "opened"
	Msent       = "sent"
	Mspammed    = "spammed"
	Msuppressed = "suppressed"
)

// MetricsPossible string array.
var MetricsPossible = [...]string{
	Mbounced,
	Mclicked,
	Mdelivered,
	Mopened,
	Msent,
	Mspammed,
	Msuppressed,
}

var metricPossibleMap map[string]bool

func init() {
	metricPossibleMap = make(map[string]bool, 1)
	for _, v := range MetricsPossible {
		metricPossibleMap[v] = true
	}
}

// GenMetricsJSONStr - generates a JSON string for given 'email' and zero
// or more 'metrics'. Validity of input arguments are user responsibility;
// will panic if json.Marshal is not happy.
func GenMetricsJSONStr(email string, metrics []string) string {
	out, jerr := json.Marshal(struct {
		Email   string   `json:"email"`
		Metrics []string `json:"metrics"`
	}{email, metrics})
	if jerr != nil {
		log.Panic(jerr)
	}
	//DBG fmt.Println("out:", string(out))
	return string(out)
}

// Gen1MetricJSONStr - generates a JSON string for given 'email' and
// a single 'metric'. Validity of input arguments are user responsibility;
// will panic if json.Marshal is not happy.
func Gen1MetricJSONStr(email, metric string) string {
	out, jerr := json.Marshal(struct {
		Email   string `json:"email"`
		Metrics string `json:"metric"`
	}{email, metric})
	if jerr != nil {
		log.Panic(jerr)
	}
	//DBG fmt.Println("out:", string(out))
	return string(out)
}

type mapiControl struct {
	reqArgsDump bool
	resBodyDump bool
}

var mapiCtrl mapiControl

// SetReqArgsDump - sets flag whether or not to dump Req[uest] input arguments.
func SetReqArgsDump(v bool) { mapiCtrl.reqArgsDump = v }

// ReqArgsDump - returns current setting for dumping Req[uest] input arguments.
func ReqArgsDump() bool { return mapiCtrl.reqArgsDump }

// SetResBodyDump - sets flag whether or not to dump Res[ponse] body.
func SetResBodyDump(v bool) { mapiCtrl.resBodyDump = v }

// ResBodyDump - returns current setting for dumping Res[ponse] body.
func ResBodyDump() bool { return mapiCtrl.resBodyDump }

func (mapi mapiControl) reqExec(verb, url1, email, jsonStr string, statusCode int) ([]byte, Meterror) {
	var err error
	var req *http.Request
	var res *http.Response
	var body []byte
	verb = strings.ToUpper(verb)
	if mapi.reqArgsDump {
		log.Printf("reqArgs:%s:verb:%q url:%q email:%q jsonStr:%q\n", fn.CurBase(), verb, url1, email, jsonStr)
	}

	if len(jsonStr) > 0 {
		if req, err = http.NewRequest(verb, url1+email, bytes.NewBuffer([]byte(jsonStr))); err != nil {
			return nil, MetError{MerrNewReqWithBody, 0, err, "func:" + fn.CurBase()}
		}
	} else {
		if req, err = http.NewRequest(verb, url1+email, nil); err != nil {
			return nil, MetError{MerrNewReqWoBody, 0, err, "func:" + fn.CurBase()}
		}
	}
	req.Header.Add("Content-Type", `application/json`)
	if res, err = http.DefaultClient.Do(req); err != nil {
		return nil, MetError{MerrClientDoErr, 0, err,
			fmt.Sprintf("verb:%q url1:%q email:%q func:%s", verb, url1, email, fn.CurBase())}
	}
	defer func() { _ = res.Body.Close() }()
	if body, err = ioutil.ReadAll(res.Body); err != nil {
		return nil, MetError{MerrBodyReadAll, 0, err, "func:" + fn.CurBase()}
	}
	if mapi.resBodyDump {
		log.Printf("%s:resBody:%v\n", fn.CurBase(), string(body))
	}
	if res.StatusCode != statusCode {
		meterr := MerrScodeWrong
		if res.StatusCode == http.StatusNotFound {
			meterr = MerrScodeNotFound
		} else if res.StatusCode == http.StatusBadRequest {
			meterr = MerrScodeBadRequest
		}
		return body, MetError{meterr, res.StatusCode, fmt.Errorf("StatusCode>got:%d want:%d funcPar:%s resBody:%v",
			res.StatusCode, statusCode, fn.LvlInfoCmn(fn.Lpar), string(body)), "ref:" + fn.LvlInfoCmn(fn.Lme)}
	}
	return body, nil
}

func unMarshalEusrMetrics(body []byte) (*EusrMetrics, Meterror) {
	var em EusrMetrics
	jerr := json.Unmarshal(body, &em)
	if jerr != nil {
		return nil, MetError{MerrJSONunMarshal, 0, jerr, "func:" + fn.CurBase()}
	}
	sort.Strings(em.Metrics)
	return &em, nil
}

//***********************************************************************************
// Exercising Core Metric API public entry points follow:
// which includes: "GET /metrics", "POST /metrics", "GET /metrics/{email}",
// "PUT /metrics/{email}", "DELETE /metrics/{email}", "GET /reset".
//***********************************************************************************

// GetMObjCol - fetches the metrics Object Collection (JSON) at the given 'url' and
// optionally dumps via MarshalIndent with provided 'pfix'.
func GetMObjCol(url1 string, dump bool, pfix string, trace bool) (*MObjCol, Meterror) {
	defer fn.LogCondTrace(trace)()
	body, err := mapiCtrl.reqExec("GET", url1+"/metrics", "", "", http.StatusOK)
	if err != nil {
		return nil, err
	}

	var moc MObjCol
	if jerr := json.Unmarshal(body, &moc); err != nil {
		return nil, MetError{MerrJSONunMarshal, 0, jerr, "func:" + fn.CurBase()}
	}
	for i := 0; i < len(moc.Grp); i++ {
		sort.Strings(moc.Grp[i].Metrics)
	}
	if dump {
		err = moc.DumpMObjCol(pfix, trace)
	}
	return &moc, err
}

// PostEusrMetrics - creates a new email-user (contained within 'metrics')
// with 'metrics' given at the specified 'url'.
func PostEusrMetrics(url1, metrics string) (*EusrMetrics, Meterror) {
	body, err := mapiCtrl.reqExec("POST", url1+"/metrics", "", metrics, http.StatusCreated)
	if err != nil {
		return nil, err
	}
	return unMarshalEusrMetrics(body)
}

// GetEusrMetrics - fetches the given email-user metrics at the given specified 'url'.
func GetEusrMetrics(url1, email string) (*EusrMetrics, Meterror) {
	body, err := mapiCtrl.reqExec("GET", url1+"/metrics/", email, "", http.StatusOK)
	if err != nil {
		if err.Code() == MerrScodeNotFound {
			merr := err.(MetError)
			merr.Code1 = MeusrNotFound
			err = merr
		}
		//MeusrNotFoundInMO    MetErrCode = 206
		return nil, err
	}
	return unMarshalEusrMetrics(body)
}

// PutEusrMetric - updates an email-user 'metric', at the given specified 'url'.
func PutEusrMetric(url1, email, metric string) (*EusrMetrics, Meterror) {
	///*DBG*/ fmt.Println(url1, email, metric)
	body, err := mapiCtrl.reqExec("PUT", url1+"/metrics/", email, metric, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return unMarshalEusrMetrics(body)
}

// DeleteEusrMetric - deletes an email-user's specified 'metric' at a given 'url'.
func DeleteEusrMetric(url1, email, metric string) Meterror {
	_, err := mapiCtrl.reqExec("DELETE", url1+"/metrics/", email, metric, http.StatusOK)
	return err
}

// ResetMetrics - resets the metrics obj collection at the specified 'url'.
func ResetMetrics(url1 string, trace bool) Meterror {
	defer fn.LogCondTrace(trace)()
	_, err := mapiCtrl.reqExec("GET", url1+"/reset", "", "", http.StatusOK)
	return err
}
