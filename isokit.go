// The Isomorphic Go Project
// Copyright (c) Wirecog, LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license, which can be found in the LICENSE file.

// Package isokit provides common isomorphic functionality intended to be used
// in an Isomorphic Go web application.
package isokit

import (
	"fmt"
	"os"
)

var (
	WebAppRoot = ""
)

func init() {

	fmt.Println("The isokit package has moved:\t 'go get -u go.isomorphicgo.org/go/isokit'")

	os.Exit(1)
}
