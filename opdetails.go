// The Isomorphic Go Project
// Copyright (c) Wirecog, LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license, which can be found in the LICENSE file.

package isokit

import "github.com/gopherjs/gopherjs/js"

type OperatingDetails struct {
	Environment int
	Runtime     int
}

const (
	ServerEnvironment = iota
	WebBrowserEnvironment
)

const (
	GoRuntime = iota
	JSRuntime
)

var (
	operatingEnvironment int
	operatingRuntime     int
)

func isJSRuntime() bool {
	return js.Global != nil
}

func isGoRuntime() bool {
	return !isJSRuntime()
}

func isWebBrowserEnvironment() bool {
	return isJSRuntime() && js.Global.Get("document") != js.Undefined
}

func isServerEnvironment() bool {

	if isGoRuntime() == true {
		return true
	} else if isJSRuntime() == true {
		return !isWebBrowserEnvironment()
	} else {
		return true
	}

}

func OperatingEnvironment() int {
	return operatingEnvironment
}

func OperatingRuntime() int {
	return operatingRuntime
}

func initializeOperatingDetails() {

	if isJSRuntime() == true {
		operatingRuntime = WebBrowserEnvironment
	}

	if isGoRuntime() == true {
		operatingRuntime = GoRuntime
	}

	if isServerEnvironment() == true {
		operatingEnvironment = ServerEnvironment
	}

	if isWebBrowserEnvironment() == true {
		operatingEnvironment = WebBrowserEnvironment
	}

}

func init() {

	initializeOperatingDetails()

}
