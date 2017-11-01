// The Isomorphic Go Project
// Copyright (c) Wirecog, LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license, which can be found in the LICENSE file.

// +build !clientonly

package isokit

import (
	"honnef.co/go/js/dom"
)

type Form interface {
	Validate() bool
	Fields() map[string]string
	Errors() map[string]string
	FormParams() *FormParams
	PrefillFields()
	SetFields(fields map[string]string)
	SetErrors(errors map[string]string)
	SetFormParams(formParams *FormParams)
	SetPrefillFields(prefillFields []string)
}

func FormValue(fp *FormParams, key string) string {

	var result string

	if OperatingEnvironment() == ServerEnvironment && fp.Request == nil {
		return ""
	}

	switch OperatingEnvironment() {

	case ServerEnvironment:

		if fp.UseFormFieldsForValidation == true {
			result = fp.FormFields[key]
		} else {
			result = fp.Request.FormValue(key)
		}

	case WebBrowserEnvironment:

		field := fp.FormElement.QuerySelector("[name=" + key + "]")

		switch field.(type) {
		case *dom.HTMLInputElement:
			result = field.(*dom.HTMLInputElement).Value
		case *dom.HTMLTextAreaElement:
			result = field.(*dom.HTMLTextAreaElement).Value
		case *dom.HTMLSelectElement:
			result = field.(*dom.HTMLSelectElement).Value

		}
	}

	return result
}
