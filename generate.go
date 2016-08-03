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

func generate(svc *ecs.ECS, clusterName string, templateFile string, options map[string]string) error {

	tickChan := time.NewTicker(time.Second * 2).C
	doneChan := make(chan bool)

	result, err := generateOne(svc, clusterName, templateFile, options)

	if err != nil {
		debug("[%s] Generating template %s failed : %s", clusterName, templateFile, err.Error())
		return err
	}
	prevResult := result

	for {
		select {
		case <-tickChan:
			result, err := generateOne(svc, clusterName, templateFile, options)
			if err != nil {
				debug("[%s] Generating template %s failed : %s", clusterName, templateFile, err.Error())
				return err
			}
			signal(options["docker-signal"], options["docker-container"], options["docker-endpoint"])

			if result != prevResult {
				prevResult = result
				debug("******** CHANGED DETECTED ****")
				signal(options["docker-signal"], options["docker-container"], options["docker-endpoint"])
			}

		case <-doneChan:
			fmt.Println("Done")
			return nil
		}
	}

}

func signal(signal string, containerName string, dockerEndpoint string) {
	// https://gist.github.com/ericchiang/c988d90edcb7eebd54de

	u, err := url.Parse(dockerEndpoint)
	if err != nil {
		debug("error parsing docker Endpoint %s", dockerEndpoint)
		return
	}

	debug("Scheme %s , Path %s", u.Scheme, u.Path)
	dial, err := net.Dial(u.Scheme, u.Path)

	if err != nil {
		debug("error opening dial ")
		return
	}
	defer dial.Close()

	form := url.Values{}
	form.Add("signal", signal)

	debug("sending %s to %s", signal, containerName)
	req, err := http.NewRequest("POST", "/containers/"+containerName+"/kill", strings.NewReader(form.Encode()))
	// We need to add the type & length otherwise nothing will happen
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	if err != nil {
		debug("error preparing kill request")
		return
	}

	conn := httputil.NewClientConn(dial, nil)
	resp, err := conn.Do(req)
	fmt.Println(resp)
	defer resp.Body.Close()
	debug("killing it")

	if err != nil {
		debug("error sending kill signal to %s", dockerEndpoint)
		return
	}

}

func generateOne(svc *ecs.ECS, clusterName string, templateFile string, options map[string]string) (string, error) {
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
