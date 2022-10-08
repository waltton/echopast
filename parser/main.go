package main

import (
	"database/sql"
	"embed"
	"log"
	"os"
	"path"

	"net/http"
	_ "net/http/pprof"

	_ "github.com/lib/pq"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	dataPath := "../data"

	db, err := sql.Open("postgres", "host=localhost database=logs user=postgres sslmode=disable")
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
		if entry.IsDir() {
			continue
		}

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
