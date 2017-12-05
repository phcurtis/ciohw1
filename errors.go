// Copyright 2017 phcurtis ciohw1 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// This file contains error codes used throughout this package.

import (
	"fmt"
	"log"
)

// MetErrCode ... metrics error code type.
type MetErrCode uint32

// MetError codes ... set metErrText map for details
const (
	//_ MetErrCode = iota
	MnoError MetErrCode = 0 // always set to zero

	// allow some gaps for future growth by category/types...
	// so it lends to errors codes being grouped and less subject to
	// today's existing codes values from changing.
	MerrBodyReadAll       MetErrCode = 100
	MerrClientDoErr       MetErrCode = 110
	MerrNewReqWithBody    MetErrCode = 121
	MerrNewReqWoBody      MetErrCode = 122
	MerrJSONmarshal       MetErrCode = 150
	MerrJSONmarshalIndent MetErrCode = 151
	MerrJSONunMarshal     MetErrCode = 152
	MerrScodeBadRequest   MetErrCode = 160
	MerrScodeNotFound     MetErrCode = 161
	MerrScodeWrong        MetErrCode = 162

	MeusrDiffThanURL      MetErrCode = 200
	MeusrDupsFound        MetErrCode = 201
	MeusrMetricsDiffInMOC MetErrCode = 202
	MeusrMetricsInvalid   MetErrCode = 203
	MeusrMetricsWrong     MetErrCode = 204
	MeusrNotFound         MetErrCode = 205
	MeusrNotFoundInMOC    MetErrCode = 206
	MerrTBD               MetErrCode = 999
)

var metErrText = map[MetErrCode]string{
	MnoError:              "no error",
	MerrBodyReadAll:       "error during body ReadAll",
	MerrClientDoErr:       "error http client do",
	MerrNewReqWithBody:    "error http.NewRequest with body",
	MerrNewReqWoBody:      "error http.NewRequest without body",
	MerrJSONmarshal:       "error json.Marshal",
	MerrJSONmarshalIndent: "error json.MarshalIndent",
	MerrJSONunMarshal:     "error json.Unmarshal",
	MerrScodeBadRequest:   "error http.statusCode='Bad Request'",
	MerrScodeNotFound:     "error http.statusCode='Not Found'",
	MerrScodeWrong:        "error http.statusCode wrong",

	MeusrDiffThanURL:      "Eusr address diff than URL",
	MeusrDupsFound:        "Eusr metrics duplicates found",
	MeusrMetricsDiffInMOC: "Eusr metrics diff from MObjCol",
	MeusrMetricsInvalid:   "Eusr metrics invalid",
	MeusrMetricsWrong:     "Eusr metrics wrong",
	MeusrNotFound:         "Eusr not found",
	MeusrNotFoundInMOC:    "Eusr not Found in MObjCol",
	MerrTBD:               "error To-Be-Determined",
}

// Meterror - metrics error interface
type Meterror interface {
	Error() string
	Code() MetErrCode
	ResCode() int
	CodeText() string // may want this to be outside interface TBD
}

// MetError - implements Meterror interface
type MetError struct {
	Code1    MetErrCode
	ResCode1 int
	Err      error
	Msg      string
}

func (m MetError) Error() string {
	// TODO:need to report bug on GoVet
	// if use this var format := "errCode:%d:(%s) ReqCode:%d err:%v, Msg:%s"
	// in statement go vet did not catch error as opposed to putting
	// format string directly in Sprintf

	text, ok := metErrText[m.Code1]
	if !ok {
		// TODO: appears log.Panicf HERE will not cause a normal panic,
		// but only prints the arguments passed to log.Panicf ...
		// should be investigated why it behaves as such;
		// log.Panicf("metErrText[%v] not found:", m.Code1)

		// TODO: need automated way to verify all MetErrCode have text representation (are in metErrText)
		log.Printf("WARNING-PLEASE-INFORM-DEVELOPER: metErrText[%v] not found:", m.Code1)
		text = "notFound-see previous log message"
	}
	return fmt.Sprintf("errCode:%d:(%s) ResCode:%d err:%v, Msg:%s",
		m.Code1, text, m.ResCode1, m.Err, m.Msg)
}

// Code - returns error code value
func (m MetError) Code() MetErrCode {
	return m.Code1
}

// CodeText - returns error code text value
func (m MetError) CodeText() string {
	text, ok := metErrText[m.Code1]
	if !ok {
		log.Printf("WARNING-PLEASE-INFORM-DEVELOPER: metErrText[%v] not found:", m.Code1)
	}
	return text
}

// ResCode - should contain http.ResCode if any otherwise should contain zero.
func (m MetError) ResCode() int {
	return m.ResCode1
}
