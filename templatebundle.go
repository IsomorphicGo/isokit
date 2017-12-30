// The Isomorphic Go Project
// Copyright (c) Wirecog, LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license, which can be found in the LICENSE file.

package isokit

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type TemplateBundle struct {
	items map[string]string
}

func NewTemplateBundle() *TemplateBundle {

	return &TemplateBundle{
		items: map[string]string{},
		// Funcs:   template.FuncMap{},
	}

}

func (t *TemplateBundle) Items() map[string]string {
	return t.items
}

func (t *TemplateBundle) importTemplateFileContents(templatesPath string) error {

	templateDirectory := filepath.Clean(templatesPath)

	if err := filepath.Walk(templateDirectory, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, TemplateFileExtension) {
			name := strings.TrimSuffix(strings.TrimPrefix(path, templateDirectory+"/"), TemplateFileExtension)
			contents, err := ioutil.ReadFile(path)
			t.items[name] = string(contents)

			if err != nil {
				log.Println("error encountered while walking directory: ", err)
				return err
			}

		}
		return nil
	}); err != nil {
		return err
	}

	return nil

}

func (t *TemplateBundle) importTemplateFileContentsForCog(templatesPath string, prefixName string, templateFileExtension string) error {

	templateDirectory := filepath.Clean(templatesPath)
	RegisterStaticAssetsSearchPath(strings.Replace(templateDirectory, string(os.PathSeparator)+"templates", "", -1))
	//println("td: ", templateDirectory)
	if err := filepath.Walk(templateDirectory, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, templateFileExtension) {
			name := strings.TrimSuffix(strings.TrimPrefix(path, templateDirectory+string(os.PathSeparator)), TemplateFileExtension)
			name = prefixName + "/" + name
			contents, err := ioutil.ReadFile(path)
			t.items[name] = string(contents)

			if err != nil {
				log.Println("error encountered while walking directory: ", err)
				return err
			}

		}
		return nil
	}); err != nil {
		return err
	}

	return nil

}
