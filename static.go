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
	"regexp"
)

var StaticAssetsPath string
var CogStaticAssetsSearchPaths []string

func RegisterStaticAssetsSearchPath(path string) {
	//fmt.Println("cog search path: ", path)
	CogStaticAssetsSearchPaths = append(CogStaticAssetsSearchPaths, path)
}

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

func bundleJavaScript(jsfiles []string) {

	var result []byte = make([]byte, 0)

	if StaticAssetsPath == "" {
		return
	}

	if _, err := os.Stat(StaticAssetsPath + "/js"); os.IsNotExist(err) {
		os.Mkdir(StaticAssetsPath+"/js", 0711)
	}

	destinationFile := StaticAssetsPath + "/js/cogimports.js"

	for i := 0; i < len(jsfiles); i++ {
		b, err := ioutil.ReadFile(jsfiles[i])
		if err != nil {
			log.Println(err)
		}
		result = append(result, b...)
	}

	err := ioutil.WriteFile(destinationFile, result, 0644)
	if err != nil {
		log.Println(err)
	}

}

func bundleCSS(cssfiles []string) {

	var result []byte = make([]byte, 0)

	if StaticAssetsPath == "" {
		return
	}

	if _, err := os.Stat(StaticAssetsPath + "/css"); os.IsNotExist(err) {
		os.Mkdir(StaticAssetsPath+"/css", 0711)
	}

	destinationFile := StaticAssetsPath + "/css/cogimports.css"

	for i := 0; i < len(cssfiles); i++ {
		b, err := ioutil.ReadFile(cssfiles[i])
		if err != nil {
			log.Println(err)
		}
		result = append(result, b...)
	}

	err := ioutil.WriteFile(destinationFile, result, 0644)
	if err != nil {
		log.Println(err)
	}

}

func BundleStaticAssets() {

	if ShouldBundleStaticAssets == false {
		return
	}

	jsfiles := findStaticAssets(".js", CogStaticAssetsSearchPaths)
	bundleJavaScript(jsfiles)
	cssfiles := findStaticAssets(".css", CogStaticAssetsSearchPaths)
	bundleCSS(cssfiles)
}

func init() {
	CogStaticAssetsSearchPaths = make([]string, 0)
}
