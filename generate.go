package main

import "bytes"
import "os"
import "bufio"
import "fmt"
import "net/url"
import "io/ioutil"
import "strings"
import "text/template"

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

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

func templateGet(filename string) (string, error) {
	u, err := url.Parse(filename)
	if err != nil {
		return "", err
	}

	if u.Scheme == "s3" {
		templateText, err := templateGetS3(u.Host, u.Path)
		if err != nil {
			return "", err
		}

		return templateText, nil
	} else {
		templateText, err := templateGetFile(filename)
		if err != nil {
			return "", err
		}
		return templateText, nil
	}
}

func templateGetFile(filename string) (string, error) {

	templateBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(templateBytes), nil
}

func templateGetS3(bucketName string, objectKey string) (string, error) {

	ses := session.New()
	//TODO get this from a global session
	svc := s3.New(ses, &aws.Config{Region: aws.String("eu-west-1")})

	params := &s3.GetObjectInput{
		Bucket: aws.String(bucketName), // Required
		Key:    aws.String(objectKey),  // Required
		/*
			IfMatch:                    aws.String("IfMatch"),
			IfModifiedSince:            aws.Time(time.Now()),
			IfNoneMatch:                aws.String("IfNoneMatch"),
			IfUnmodifiedSince:          aws.Time(time.Now()),
			Range:                      aws.String("Range"),
			RequestPayer:               aws.String("RequestPayer"),
			ResponseCacheControl:       aws.String("ResponseCacheControl"),
			ResponseContentDisposition: aws.String("ResponseContentDisposition"),
			ResponseContentEncoding:    aws.String("ResponseContentEncoding"),
			ResponseContentLanguage:    aws.String("ResponseContentLanguage"),
			ResponseContentType:        aws.String("ResponseContentType"),
			ResponseExpires:            aws.Time(time.Now()),
			SSECustomerAlgorithm:       aws.String("SSECustomerAlgorithm"),
			SSECustomerKey:             aws.String("SSECustomerKey"),
			SSECustomerKeyMD5:          aws.String("SSECustomerKeyMD5"),
			VersionId:                  aws.String("ObjectVersionId"),
		*/
	}
	resp, err := svc.GetObject(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return "", err
	}

	debug("fetching template from s3")
	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error reading content from s3", err)
		return "", err
	}

	return string(data), nil
}

func templateExecute(ecsWatchInfo EcsWatchInfo, options EcsWatchTrackOptions) (*bytes.Buffer, error) {
	templateInputFile := options.TemplateInputFile

	templateText, err := templateGet(templateInputFile)

	tmpl, err := templateNew("ecswatch-template").Parse(templateText)
	if err != nil {
		return nil, err
		//log.Fatalf("Unable to parse template: %s", err)
	}

	buf := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(buf, "ecswatch-template", &ecsWatchInfo)
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
