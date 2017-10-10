package isokit

import (
	"bytes"
	"encoding/gob"
	"html/template"

	"honnef.co/go/js/xhr"
)

func FetchTemplateBundle(templateSetChannel chan *TemplateSet) {

	data, err := xhr.Send("POST", "/template-bundle", nil)
	if err != nil {
		println("Encountered error: ", err)
		println(err)
	}
	var templateBundleMap map[string]string
	var templateBundleMapBuffer bytes.Buffer

	dec := gob.NewDecoder(&templateBundleMapBuffer)
	templateBundleMapBuffer = *bytes.NewBuffer(data)
	err = dec.Decode(&templateBundleMap)

	if err != nil {
		println("Encountered error: ", err)
		panic(err)
	}

	templateSet := NewTemplateSet()
	err = templateSet.ImportTemplatesFromMap(templateBundleMap)

	if err != nil {
		println("Encountered import error: ", err)
		panic(err)
	}

	templateSetChannel <- templateSet
}

func FetchTemplateBundleWithSuppliedFunctionMap(templateSetChannel chan *TemplateSet, funcMap template.FuncMap) {

	data, err := xhr.Send("POST", "/template-bundle", nil)
	if err != nil {
		println("Encountered error: ", err)
		println(err)
	}
	var templateBundleMap map[string]string
	var templateBundleMapBuffer bytes.Buffer

	dec := gob.NewDecoder(&templateBundleMapBuffer)
	templateBundleMapBuffer = *bytes.NewBuffer(data)
	err = dec.Decode(&templateBundleMap)

	if err != nil {
		println("Encountered error: ", err)
		panic(err)
	}

	templateSet := NewTemplateSet()
	templateSet.Funcs = funcMap
	err = templateSet.ImportTemplatesFromMap(templateBundleMap)

	if err != nil {
		println("Encountered import error: ", err)
		panic(err)
	}

	templateSetChannel <- templateSet
}
