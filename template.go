// The Isomorphic Go Project
// Copyright (c) Wirecog, LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license, which can be found in the LICENSE file.

package isokit

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"log"
	"strings"

	"honnef.co/go/js/dom"
)

const (
	TemplateRegular = iota
	TemplatePartial
	TemplateLayout
)

var (
	PrefixNamePartial            = "partials/"
	PrefixNameLayout             = "layouts/"
	TemplateFileExtension        = ".tmpl"
	TemplateFilesPath            = "./templates"
	UseStaticTemplateBundleFile  = false
	StaticTemplateBundleFilePath = ""
	ShouldBundleStaticAssets     = true
)

type Template struct {
	*template.Template
	templateType int8
}

const (
	PlacementAppendTo = iota
	PlacementReplaceInnerContents
	PlacementInsertBefore
)

type RenderParams struct {
	Data                          interface{}
	Writer                        io.Writer
	Element                       dom.Element
	Disposition                   int8
	Attributes                    map[string]string
	ShouldPopulateRenderedContent bool
	RenderedContent               string
	ShouldSkipFinalRenderStep     bool
	PageTitle                     string
}

func (t *Template) GetTemplateType() int8 {

	if t == nil {
		return -1
	} else {
		return t.templateType
	}
}

func (t *Template) NameWithPrefix() string {

	var prefixName string
	switch t.templateType {

	case TemplateRegular:
		prefixName = ""

	case TemplatePartial:
		prefixName = PrefixNamePartial

	case TemplateLayout:
		prefixName = PrefixNameLayout

	}

	if strings.HasPrefix(t.Name(), prefixName) {
		return t.Name()
	} else {
		return prefixName + t.Name()
	}

}

func (t *Template) Render(params *RenderParams) error {

	if OperatingEnvironment() == ServerEnvironment && (params.Writer == nil) {
		return errors.New("Either the response writer and/or the request is nil!")
	}

	if OperatingEnvironment() == WebBrowserEnvironment && params.Element == nil {
		return errors.New("The element to render relative to is nil!")
	}

	switch OperatingEnvironment() {
	case WebBrowserEnvironment:
		t.RenderTemplateOnClient(params)

	case ServerEnvironment:
		t.RenderTemplateOnServer(params)
	}

	return nil
}

func (t *Template) RenderTemplateOnClient(params *RenderParams) {

	var tpl bytes.Buffer

	if err := t.Execute(&tpl, params.Data); err != nil {
		log.Println("Error encountered when attempting to render template on client: ", err)
	}

	if params.ShouldPopulateRenderedContent == true {
		params.RenderedContent = string(tpl.Bytes())
	}

	if params.ShouldSkipFinalRenderStep == true {
		return
	}

	div := dom.GetWindow().Document().CreateElement("div").(*dom.HTMLDivElement)
	div.SetInnerHTML(string(tpl.Bytes()))

	if _, ok := params.Attributes["id"]; ok {
		div.SetID(params.Attributes["id"])
	}

	if _, ok := params.Attributes["class"]; ok {
		div.SetAttribute("class", params.Attributes["class"])
	}

	switch params.Disposition {
	case PlacementAppendTo:
		params.Element.AppendChild(div)
	case PlacementReplaceInnerContents:
		params.Element.SetInnerHTML(div.OuterHTML())
	case PlacementInsertBefore:
		params.Element.ParentNode().InsertBefore(div, params.Element)
	default:
		params.Element.AppendChild(div)
	}

	if params.PageTitle != "" && params.ShouldPopulateRenderedContent == false {
		dom.GetWindow().Document().Underlying().Set("title", params.PageTitle)
	}

}

func (t *Template) RenderTemplateOnServer(params *RenderParams) {

	w := params.Writer
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, params.Data); err != nil {
		log.Println("Error encountered when attempting to render template on server: ", err)
	}
	w.Write(tpl.Bytes())
}
