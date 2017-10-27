// The Isomorphic Go Project
// Copyright (c) Wirecog, LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license, which can be found in the LICENSE file.

// +build !clientonly

package isokit

import (
	"net/http"

	"honnef.co/go/js/dom"
)

type FormParams struct {
	FormElement                *dom.HTMLFormElement
	ResponseWriter             http.ResponseWriter
	Request                    *http.Request
	UseFormFieldsForValidation bool
	FormFields                 map[string]string
}
