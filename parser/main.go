package main

import (
	"database/sql"
	"log"
	"os"
	"path"
)

func main() {
	dataPath := "../data"

	db, err := sql.Open("host=localhost database=logs username=postgres", "postgers")
	if err != nil {
		log.Panic(err)
	}

	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}

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

	err = writeLogs(db, logs)
	if err != nil {
		log.Panic(err)
	}
}
