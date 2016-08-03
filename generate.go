package main

//import "os"
//import "strconv"
import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/service/ecs"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"
)

func generate(svc *ecs.ECS, clusterName string, templateFile string, outputFile string) error {

	tickChan := time.NewTicker(time.Second * 2).C
	doneChan := make(chan bool)

	result, err := generateOne(svc, clusterName, templateFile, outputFile)

	if err != nil {
		debug("[%s] Generating template %s failed : %s", clusterName, templateFile, err.Error())
		return err
	}
	prevResult := result

	for {
		select {
		case <-tickChan:
			result, err := generateOne(svc, clusterName, templateFile, outputFile)
			if err != nil {
				debug("[%s] Generating template %s failed : %s", clusterName, templateFile, err.Error())
				return err
			}

			if result != prevResult {
				prevResult = result
				debug("******** CHANGED DETECTED ****")
				signal("SIGHUP", "ecswatch_nginx_1", "/var/run/docker.sock")
			}

		case <-doneChan:
			fmt.Println("Done")
			return nil
		}
	}

}

func signal(signal string, containerName string, socketFile string) {
	// https://gist.github.com/ericchiang/c988d90edcb7eebd54de
	form := url.Values{}
	form.Add("signal", "SIGHUP")

	dial, err := net.Dial("unix", socketFile)
	if err != nil {
		debug("error opening dial ")
		return
	}
	defer dial.Close()

	req, err := http.NewRequest("POST", "/containers/"+containerName+"/kill", strings.NewReader(form.Encode()))
	// We need to add the type & length otherwise nothing will happen
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	if err != nil {
		debug("error opening socker file /var/run/docker.sock")
		return
	}

	conn := httputil.NewClientConn(dial, nil)
	resp, err := conn.Do(req)
	fmt.Println(resp)
	defer resp.Body.Close()
	debug("killing it")

	if err != nil {
		debug("error opening socker file /var/run/docker.sock")
		return
	}

}

func generateOne(svc *ecs.ECS, clusterName string, templateFile string, outputFile string) (string, error) {
	var watchInfo, err = getEcsWatchInfo(svc, clusterName)

	if err != nil {
		debug("[%s] Retrieving report ECS Cluster failed: %s", clusterName, err.Error())
		return "", err
	}

	result, err := executeTemplate(templateFile, *watchInfo)

	if err != nil {
		debug("[%s] Generating template failed: %s", clusterName, err.Error())
		return "", err
	}

	return result.String(), nil
}

func executeTemplate(templatePath string, ecsWatchInfo EcsWatchInfo) (*bytes.Buffer, error) {
	tmpl, err := newTemplate(filepath.Base(templatePath)).ParseFiles(templatePath)
	if err != nil {
		return nil, err
		//log.Fatalf("Unable to parse template: %s", err)
	}

	buf := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(buf, filepath.Base(templatePath), &ecsWatchInfo)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func newTemplate(name string) *template.Template {
	tmpl := template.New(name).Funcs(template.FuncMap{})
	return tmpl
}
