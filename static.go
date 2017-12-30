// The Isomorphic Go Project
// Copyright (c) Wirecog, LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license, which can be found in the LICENSE file.

// +build !clientonly

package isokit

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
)

var StaticAssetsPath string
var ShouldMinifyStaticAssets bool

func findStaticAssets(ext string, paths []string) []string {

	var files []string

	for i := 0; i < len(paths); i++ {
		//fmt.Println("file search path: ", paths[i])
		filepath.Walk(paths[i], func(path string, f os.FileInfo, _ error) error {
			if !f.IsDir() {
				r, err := regexp.MatchString(ext, f.Name())
				if err == nil && r {
					files = append(files, path)
				}
			}
			return nil
		})

	}
	return files
}

func bundleJavaScript(jsfiles []string, shouldMinify bool) {

	outputFileName := "cogimports.js"
	if shouldMinify == true {
		outputFileName = "cogimports.min.js"
	}

	var result []byte = make([]byte, 0)

	if StaticAssetsPath == "" {
		return
	}

	if _, err := os.Stat(filepath.Join(StaticAssetsPath, "js")); os.IsNotExist(err) {
		os.Mkdir(filepath.Join(StaticAssetsPath, "js"), 0711)
	}

	destinationFile := filepath.Join(StaticAssetsPath, "js", outputFileName)

	for i := 0; i < len(jsfiles); i++ {
		b, err := ioutil.ReadFile(jsfiles[i])
		if err != nil {
			log.Println(err)
		}
		result = append(result, b...)
	}

	if shouldMinify == true {

		m := minify.New()
		m.AddFunc("text/javascript", js.Minify)
		b, err := m.Bytes("text/javascript", result)

		if err != nil {
			log.Println(err)
		}

		err = ioutil.WriteFile(destinationFile, b, 0644)

		if err != nil {
			log.Println(err)
		}

	} else {
		err := ioutil.WriteFile(destinationFile, result, 0644)

		if err != nil {
			log.Println(err)
		}

	}

}

func bundleCSS(cssfiles []string, shouldMinify bool) {

	outputFileName := "cogimports.css"
	if shouldMinify == true {
		outputFileName = "cogimports.min.css"
	}

	var result []byte = make([]byte, 0)

	if StaticAssetsPath == "" {
		return
	}

	if _, err := os.Stat(filepath.Join(StaticAssetsPath, "css")); os.IsNotExist(err) {
		os.Mkdir(filepath.Join(StaticAssetsPath, "css"), 0711)
	}

	destinationFile := filepath.Join(StaticAssetsPath, "css", outputFileName)

	for i := 0; i < len(cssfiles); i++ {
		b, err := ioutil.ReadFile(cssfiles[i])
		if err != nil {
			log.Println(err)
		}
		result = append(result, b...)
	}

	if shouldMinify == true {

		m := minify.New()
		m.AddFunc("text/css", css.Minify)
		b, err := m.Bytes("text/css", result)

		if err != nil {
			log.Println(err)
		}

		err = ioutil.WriteFile(destinationFile, b, 0644)

	} else {
		err := ioutil.WriteFile(destinationFile, result, 0644)
		if err != nil {
			log.Println(err)
		}

	}

}

func BundleStaticAssets() {

	if ShouldBundleStaticAssets == false {
		return
	}

	jsfiles := findStaticAssets(".js", CogStaticAssetsSearchPaths)
	bundleJavaScript(jsfiles, ShouldMinifyStaticAssets)
	cssfiles := findStaticAssets(".css", CogStaticAssetsSearchPaths)
	bundleCSS(cssfiles, ShouldMinifyStaticAssets)
}

func init() {
	CogStaticAssetsSearchPaths = make([]string, 0)
}
