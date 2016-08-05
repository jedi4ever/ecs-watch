package main

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

func dockerSignal(signal string, containerName string, dockerEndpoint string) {
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
