package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/pkg/errors"
)

func main() {
	dataPath := "../data"

	files, err := os.ReadDir(dataPath)
	if err != nil {
		log.Panic(err)
	}

	type Log struct {
		Request   string `json:"request"`
		Timestamp string `json:"timestamp"`
	}
	var logs []Log

	for _, entry := range files {
		f, err := os.Open(path.Join(dataPath, entry.Name()))
		if err != nil {
			log.Panic(err)
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

				log.Fatalf("read file line error: %v", err)
				return
			}

			var logline Log
			err = json.Unmarshal([]byte(line), &logline)
			if err != nil {
				err = errors.Wrapf(err, "fail to parse line #%d '%s' on file '%s'", ln, logline, entry.Name())
				log.Panic(err)
			}

			logs = append(logs, logline)
		}
	}

	for i, l := range logs {
		b := bytes.NewBufferString(l.Request)
		req, err := http.ReadRequest(bufio.NewReader(b))
		if err != nil {
			log.Panic(err)
		}

		if req.URL.String() == "/" {
			continue
		}

		fmt.Printf("%d --------------------\n", i)

		fmt.Println("req.URL", req.URL)
		// fmt.Println("req.Header")
		// for k, v := range req.Header {
		// 	fmt.Println("\t", k, v)
		// }
	}

}
