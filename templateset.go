// The Isomorphic Go Project
// Copyright (c) Wirecog, LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license, which can be found in the LICENSE file.

package isokit

import (
	"bytes"
	"encoding/gob"
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var CogStaticAssetsSearchPaths []string

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
			log.Println("Encountered error when attempting to parse template: ", err)

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

func (t *TemplateSet) PersistTemplateBundleToDisk() error {

	dirPath := filepath.Dir(StaticTemplateBundleFilePath)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {

		return errors.New("The specified directory for the StaticTemplateBundleFilePath, " + dirPath + ", does not exist!")

	} else {

		var templateContentItemsBuffer bytes.Buffer
		enc := gob.NewEncoder(&templateContentItemsBuffer)
		m := t.bundle.Items()
		err := enc.Encode(&m)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(StaticTemplateBundleFilePath, templateContentItemsBuffer.Bytes(), 0644)
		if err != nil {
			return err
		} else {
			return nil
		}

	}

}

func (t *TemplateSet) RestoreTemplateBundleFromDisk() error {

	if _, err := os.Stat(StaticTemplateBundleFilePath); os.IsNotExist(err) {
		return errors.New("The StaticTemplateBundleFilePath, " + StaticTemplateBundleFilePath + ", does not exist")
	} else {

		data, err := ioutil.ReadFile(StaticTemplateBundleFilePath)
		if err != nil {
			return err
		}

		var templateBundleMap map[string]string
		var templateBundleMapBuffer bytes.Buffer
		dec := gob.NewDecoder(&templateBundleMapBuffer)
		templateBundleMapBuffer = *bytes.NewBuffer(data)
		err = dec.Decode(&templateBundleMap)

		if err != nil {
			return err
		}

		t.ImportTemplatesFromMap(templateBundleMap)
		bundle := &TemplateBundle{items: templateBundleMap}
		t.bundle = bundle

		return nil
	}
}

func (t *TemplateSet) GatherTemplates() {

	if UseStaticTemplateBundleFile == true {
		err := t.RestoreTemplateBundleFromDisk()
		if err != nil {
			log.Println("Didn't find a template bundle from disk, will generate a new template bundle.")
		} else {
			return
		}
	}

	bundle := NewTemplateBundle()

	templatesPath := t.TemplateFilesPath
	if templatesPath == "" {
		templatesPath = TemplateFilesPath
	}
	bundle.importTemplateFileContents(templatesPath)
	t.ImportTemplatesFromMap(bundle.Items())
	t.bundle = bundle

	if StaticTemplateBundleFilePath != "" {
		err := t.PersistTemplateBundleToDisk()
		if err != nil {
			log.Println("Failed to persist the template bundle to disk, in GatherTemplates, with error: ", err)
		}
	}

}

func (t *TemplateSet) GatherCogTemplates(cogTemplatePath string, prefixName string, templateFileExtension string) {

	if ShouldBundleStaticAssets == false || UseStaticTemplateBundleFile == true {
		return
	}

	bundle := NewTemplateBundle()

	templatesPath := cogTemplatePath
	bundle.importTemplateFileContentsForCog(templatesPath, prefixName, templateFileExtension)
	t.ImportTemplatesFromMap(bundle.Items())

	for k, v := range bundle.Items() {
		t.bundle.items[k] = v
	}

	if StaticTemplateBundleFilePath != "" {
		err := t.PersistTemplateBundleToDisk()
		if err != nil {
			log.Println("Failed to persist the template bundle to disk, in GatherCogTemplates, with error: ", err)
		}
	}

}

func StaticTemplateBundleFileExists() bool {

	if _, err := os.Stat(StaticTemplateBundleFilePath); os.IsNotExist(err) {
		return false
	} else {
		return true
	}

}

func RegisterStaticAssetsSearchPath(path string) {
	//fmt.Println("cog search path: ", path)
	CogStaticAssetsSearchPaths = append(CogStaticAssetsSearchPaths, path)
}
