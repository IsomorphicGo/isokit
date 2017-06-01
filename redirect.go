// The Isomorphic Go Project
// Copyright (c) Wirecog, LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license, which can be found in the LICENSE file.

package isokit

import (
	"net/http"

	"github.com/go-humble/detect"
	"github.com/gopherjs/gopherjs/js"
)

// ServerRedirect performs a redirect when operating on the server-side.
func ServerRedirect(w http.ResponseWriter, r *http.Request, destinationURL string) {
	http.Redirect(w, r, destinationURL, 301)
}

// ClientRedirect performs a redirect when operating on the client-side.
func ClientRedirect(destinationURL string) {
	js := js.Global
	js.Get("location").Set("href", destinationURL)
}

// Redirect performs a redirect operation based on the environment that
// the program is operating under.
func Redirect(items ...interface{}) {
	var url string
	var w http.ResponseWriter
	var r *http.Request

	if detect.IsServer() && len(items) != 3 {
		return
	}

	if detect.IsBrowser() && len(items) != 1 {
		return
	}

	for _, item := range items {

		switch t := item.(type) {
		case http.ResponseWriter:
			w = t
		case *http.Request:
			r = t
		case string:
			url = t
		}

	}

	if detect.IsServer() && (w == nil || r == nil) {
		return
	}

	if url == "" {
		return
	}

	switch {
	case detect.IsBrowser():
		ClientRedirect(url)

	case detect.IsServer():
		ServerRedirect(w, r, url)
	}
}
