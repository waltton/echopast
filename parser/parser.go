package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type Log struct {
	RawRequest string `json:"request"`
	Timestamp  string `json:"timestamp"`
	Request    Request
}

type Request struct {
	URL            string
	Host           string
	UserAgent      string
	AcceptEncoding string
	Accept         string
	Cookie         string
	IP             net.IP
	Protocol       string
	Headers        map[string][]string
}

func parseFile(name string) (logs []Log, err error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	buf := bufio.NewReader(f)

	ln := 0
	for {
		ln++

		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		var logline Log
		err = json.Unmarshal([]byte(line), &logline)
		if err != nil {
			err = errors.Wrapf(err, "fail to parse line #%d '%s' on file '%s'", ln, logline, name)
			log.Panic(err)
		}

		logline.Request = parseRequest(logline)

		logs = append(logs, logline)
	}

	return logs, nil
}

var reForwarded = regexp.MustCompile(`for="(.*)";proto=(.*),`)

func parseRequest(l Log) Request {
	b := bytes.NewBufferString(l.RawRequest)
	req, err := http.ReadRequest(bufio.NewReader(b))
	if err != nil {
		log.Panic(err)
	}

	var ip net.IP
	var protocol string
	m := reForwarded.FindStringSubmatch(req.Header.Get("Forwarded"))
	if len(m) >= 3 {
		ip = net.ParseIP(m[1])
		protocol = m[2]
	}

	headers := make(map[string][]string)
	for k, v := range req.Header {
		if strings.HasPrefix(k, "X-Appengine") || k == "Via" || k == "User-Agent" || k == "Traceparent" ||
			k == "X-Cloud-Trace-Context" || k == "X-Google-Serverless-Node-Envoy-Config-Gae" || k == "X-Forwarded-Proto" ||
			k == "X-Forwarded-For" || k == "Forwarded" || k == "Accept-Encoding" || k == "Accept" || k == "Cookie" {
			continue
		}

		headers[k] = v
	}

	return Request{
		URL:            req.URL.String(),
		Host:           req.Host,
		UserAgent:      req.Header.Get("User-Agent"),
		AcceptEncoding: req.Header.Get("Accept-Encoding"),
		Accept:         req.Header.Get("Accept"),
		Cookie:         req.Header.Get("Cookie"),
		IP:             ip,
		Protocol:       protocol,
		Headers:        headers,
	}
}
