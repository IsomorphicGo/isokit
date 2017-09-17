// The Isomorphic Go Project
// Copyright (c) Wirecog, LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license, which can be found in the LICENSE file.

package isokit

import (
	"html/template"
	"io/ioutil"
	"strings"
)

type TemplateSet struct {
	members           map[string]*Template
	Funcs             template.FuncMap
	bundle            *TemplateBundle
	TemplateFilesPath string
}

func NewTemplateSet() *TemplateSet {
	return &TemplateSet{
		members: map[string]*Template{},
		Funcs:   template.FuncMap{},
	}
}

func (t *TemplateSet) Members() map[string]*Template {
	return t.members
}

func (t *TemplateSet) Bundle() *TemplateBundle {
	return t.bundle
}

func (t *TemplateSet) AddTemplateFile(name, filename string, templateType int8) error {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	tpl, err := template.New(name).Funcs(t.Funcs).Parse(string(contents))
	template := Template{
		Template:     tpl,
		templateType: templateType,
	}

	t.members[tpl.Name()] = &template
	return nil

}

func (t *TemplateSet) MakeAllAssociations() error {

	for _, template := range t.members {

		for _, member := range t.members {

			if member.Lookup(template.NameWithPrefix()) == nil {

				if _, err := member.AddParseTree(template.NameWithPrefix(), template.Tree); err != nil {
					println(err)
					return err
				}
			}

		}

	}
	return nil
}

func (t *TemplateSet) ImportTemplatesFromMap(templateMap map[string]string) error {

	for name, templateContents := range templateMap {

		var templateType int8
		if strings.Contains(name, PrefixNamePartial) {
			templateType = TemplatePartial
		} else if strings.Contains(name, PrefixNameLayout) {
			templateType = TemplateLayout
		} else {
			templateType = TemplateRegular
		}

		tpl, err := template.New(name).Funcs(t.Funcs).Parse(templateContents)

		if err != nil {
			println("Encountered error when attempting to parse template: ", err)

			return err
		}

		tmpl := Template{
			Template:     tpl,
			templateType: templateType,
		}
		t.members[name] = &tmpl

	}
	t.MakeAllAssociations()
	return nil
}

func (t *TemplateSet) Render(templateName string, params *RenderParams) {

	t.Members()[templateName].Render(params)

}

func (t *TemplateSet) GatherTemplates() {

	bundle := NewTemplateBundle()

	templatesPath := t.TemplateFilesPath
	if templatesPath == "" {
		templatesPath = TemplateFilesPath
	}
	bundle.importTemplateFileContents(templatesPath)
	t.ImportTemplatesFromMap(bundle.Items())
	t.bundle = bundle

}

func (t *TemplateSet) GatherCogTemplates(cogTemplatePath string, prefixName string, templateFileExtension string) {

	bundle := NewTemplateBundle()

	templatesPath := cogTemplatePath
	bundle.importTemplateFileContentsForCog(templatesPath, prefixName, templateFileExtension)
	t.ImportTemplatesFromMap(bundle.Items())

	for k, v := range bundle.Items() {
		t.bundle.items[k] = v
	}
	//	t.bundle = bundle

}
