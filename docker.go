package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

func dockerSignal(signal string, containerName string, dockerEndpoint string) (err error) {

	// Default we assume the containerRef is a name
	containerRefs := []string{
		containerName,
	}

	if strings.HasPrefix(containerName, "label:") {
		debug("found a docker label filter")
		label := strings.Split(containerName, ":")[1]
		dockerIds, err := dockerLookupByLabel(dockerEndpoint, label)
		debug("docker ids found by label %s", dockerIds)
		if err != nil {
			return err
		}
		containerRefs = dockerIds

	}

	// https://gist.github.com/ericchiang/c988d90edcb7eebd54de

	u, err := url.Parse(dockerEndpoint)
	if err != nil {
		debug("error parsing docker Endpoint %s", dockerEndpoint)
		return err
	}

	debug("Scheme %s , Path %s", u.Scheme, u.Path)
	dial, err := net.Dial(u.Scheme, u.Path)

	if err != nil {
		debug("error opening dial ")
		return err
	}
	defer dial.Close()

	form := url.Values{}
	form.Add("signal", signal)

	for _, containerRef := range containerRefs {
		debug("sending %s to %s", signal, containerRef)
		req, err := http.NewRequest("POST", "/containers/"+containerRef+"/kill", strings.NewReader(form.Encode()))
		// We need to add the type & length otherwise nothing will happen
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

		if err != nil {
			debug("error preparing kill request")
			return err
		}

		conn := httputil.NewClientConn(dial, nil)
		resp, err := conn.Do(req)
		fmt.Println(resp)
		defer resp.Body.Close()
		debug("killing it")

		if err != nil {
			debug("error sending kill signal to %s", dockerEndpoint)
			return err
		}
	}

	return nil

}

// --data-urlencode 'filters={"label":["com.amazonaws.ecs.container-name=ecs-nginx"]}'
func dockerLookupByLabel(dockerEndpoint string, containerLabel string) (containerIds []string, err error) {

	containerIds = []string{}

	u, err := url.Parse(dockerEndpoint)
	if err != nil {
		debug("error parsing docker Endpoint %s", dockerEndpoint)
		return containerIds, err
	}

	debug("Scheme %s , Path %s", u.Scheme, u.Path)
	dial, err := net.Dial(u.Scheme, u.Path)

	if err != nil {
		debug("error opening dial ")
		return containerIds, err
	}
	defer dial.Close()

	req, err := http.NewRequest("GET", "/containers/json", nil)
	if err != nil {
		debug("error creating docker container http request %s", err.Error())
		return containerIds, err
	}

	q := req.URL.Query()

	filter := map[string][]string{}
	filter["label"] = []string{containerLabel}
	filterJson, err := json.Marshal(filter)
	if err != nil {
		debug("error creating docker Filter %s", err.Error())
		return containerIds, err
	}

	q.Add("filters", string(filterJson))
	req.URL.RawQuery = q.Encode()

	conn := httputil.NewClientConn(dial, nil)
	resp, err := conn.Do(req)
	if err != nil {
		debug("error requesting docker info %s", err.Error())
		return containerIds, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	debug("%s", string(body))

	type DockerContainer struct{ Id string }
	type DockerContainers []DockerContainer

	dockerContainers := DockerContainers{}

	err = json.Unmarshal(body, &dockerContainers)
	if err != nil {
		debug("failed to parse docker Container json %s", err.Error())
		return containerIds, err
	}

	if len(dockerContainers) == 0 {
		debug("no containers found matching the label filter")
		return containerIds, nil
	}

	for _, container := range dockerContainers {
		containerIds = append(containerIds, container.Id)
	}

	//resp_body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return containerIds, nil
}
