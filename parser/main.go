package main

import (
	"fmt"
	"log"
	"os"
	"path"
)

func main() {
	dataPath := "../data"

	files, err := os.ReadDir(dataPath)
	if err != nil {
		log.Panic(err)
	}

	var logs []Log

	for _, entry := range files {
		l, err := parseFile(path.Join(dataPath, entry.Name()))
		if err != nil {
			log.Panic(err)
		}

		logs = append(logs, l...)
	}

	for _, l := range logs {
		fmt.Println(l.Request)
	}
}
