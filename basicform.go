// The Isomorphic Go Project
// Copyright (c) Wirecog, LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license, which can be found in the LICENSE file.

package isokit

import "strings"

type BasicForm struct {
	formParams *FormParams

	autofillFields []string
	fields         map[string]string
	errors         map[string]string
}

func (c *BasicForm) AutofillFields() []string {
	return c.autofillFields
}

func (c *BasicForm) Fields() map[string]string {
	return c.fields
}

func (c *BasicForm) Errors() map[string]string {
	return c.errors
}

func (c *BasicForm) FormParams() *FormParams {
	return c.formParams

}

func (c *BasicForm) SetAutofillFields(autofillFields []string) {
	c.autofillFields = autofillFields
}

func (c *BasicForm) SetFields(fields map[string]string) {
	c.fields = fields
}

func (c *BasicForm) SetErrors(errors map[string]string) {
	c.errors = errors
}

func (c *BasicForm) SetFormParams(formParams *FormParams) {
	c.formParams = formParams
}

func (c *BasicForm) SetError(key string, message string) {
	c.errors[key] = message
}

func (c *BasicForm) ClearErrors() {
	c.errors = make(map[string]string)
}

func (c *BasicForm) PopulateFields() {
	for _, fieldName := range c.autofillFields {
		c.fields[fieldName] = FormValue(c.FormParams(), fieldName)
	}
}

func (c *BasicForm) DisplayErrors() {
	if OperatingEnvironment() == WebBrowserEnvironment && c.formParams.FormElement != nil {
		errorSpans := c.formParams.FormElement.QuerySelectorAll(".formError")
		for _, v := range errorSpans {
			v.SetInnerHTML(c.errors[strings.Replace(v.GetAttribute("id"), "Error", "", -1)])
		}
	}
}
