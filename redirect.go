// The Isomorphic Go Project
// Copyright (c) Wirecog, LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license, which can be found in the LICENSE file.

// +build !clientonly

package isokit

import (
	"errors"
	"net/http"

	"github.com/gopherjs/gopherjs/js"
)

// ServerRedirect performs a redirect when operating on the server-side.
func ServerRedirect(w http.ResponseWriter, r *http.Request, destinationURL string) {
	http.Redirect(w, r, destinationURL, 302)
}

// ClientRedirect performs a redirect when operating on the client-side.
func ClientRedirect(destinationURL string) {
	js := js.Global
	js.Get("location").Set("href", destinationURL)
}

type RedirectParams struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	URL            string
}

func Redirect(params *RedirectParams) error {

	if params.URL == "" {
		return errors.New("The URL must be specified!")
	}

	if OperatingEnvironment() == ServerEnvironment && (params.ResponseWriter == nil || params.Request == nil) {
		return errors.New("Either the response writer and/or the request is nil!")
	}

	switch OperatingEnvironment() {
	case WebBrowserEnvironment:
		ClientRedirect(params.URL)

	case ServerEnvironment:
		ServerRedirect(params.ResponseWriter, params.Request, params.URL)
	}

	return nil
}
