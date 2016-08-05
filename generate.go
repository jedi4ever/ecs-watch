package main

import "bytes"
import "path/filepath"
import "text/template"

func templateGenerate(ecsWatchInfo EcsWatchInfo, options EcsWatchTrackOptions) error {
	_, err := templateGenerateString(ecsWatchInfo, options)
	//result, err := templateGenerateString(ecsWatchInfo, options)

	if err != nil {
		debug("Generating template failed: %s", err.Error())
		return err
	}
	return nil
}

func templateGenerateString(ecsWatchInfo EcsWatchInfo, options EcsWatchTrackOptions) (string, error) {

	result, err := templateExecute(ecsWatchInfo, options)

	if err != nil {
		debug("Generating template failed: %s", err.Error())
		return "", err
	}

	return result.String(), nil
}

func templateExecute(ecsWatchInfo EcsWatchInfo, options EcsWatchTrackOptions) (*bytes.Buffer, error) {
	templateInputFile := options.TemplateInputFile

	tmpl, err := templateNew(filepath.Base(templateInputFile)).ParseFiles(templateInputFile)
	if err != nil {
		return nil, err
		//log.Fatalf("Unable to parse template: %s", err)
	}

	buf := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(buf, filepath.Base(templateInputFile), &ecsWatchInfo)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func templateNew(name string) *template.Template {
	tmpl := template.New(name).Funcs(template.FuncMap{})
	return tmpl
}
