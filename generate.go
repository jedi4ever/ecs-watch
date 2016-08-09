package main

import "bytes"
import "os"
import "bufio"
import "fmt"
import "strings"
import "path/filepath"
import "text/template"

func templateGenerate(ecsWatchInfo EcsWatchInfo, options EcsWatchTrackOptions) error {
	result, err := templateGenerateString(ecsWatchInfo, options)
	//result, err := templateGenerateString(ecsWatchInfo, options)

	debug("Generating template")

	if options.TemplateOutputFile != "" {
		f, err := os.Create(options.TemplateOutputFile)
		if err != nil {
			debug("Error creating outpurt file output failed: %s", err.Error())
			return err
		}
		w := bufio.NewWriter(f)
		n, err := w.WriteString(result)
		debug("wrote %d bytes", n)
		if err != nil {
			debug("Error writing output file output failed: %s", err.Error())
			return err
		}
		w.Flush()

	} else {
		fmt.Println(result)
	}
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
	tmpl := template.New(name).Funcs(template.FuncMap{
		"groupByVirtualHost": groupByVirtualHost,
		"replace":            strings.Replace,
	})
	return tmpl
}

func groupByVirtualHost(ecsWatchInfo EcsWatchInfo) map[string]EcsWatchInfo {

	infoByHosts := make(map[string]EcsWatchInfo)

	for _, infoItem := range ecsWatchInfo {
		if infoItem.HostPort != 0 {
			virtualHost, found := infoItem.Environment["VIRTUAL_HOST"]
			if found {
				infoByHosts[virtualHost] = append(infoByHosts[virtualHost], infoItem)
			}
		}
	}

	return infoByHosts
}
