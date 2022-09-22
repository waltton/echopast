package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/microcosm-cc/bluemonday"
)

func main() {
	bucket := os.Getenv("GCLOUD_STORAGE_BUCKET")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sc, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer sc.Close()

	p := bluemonday.UGCPolicy()

	http.HandleFunc("/", indexHandler(sc, bucket, p))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func indexHandler(sc *storage.Client, bucket string, p *bluemonday.Policy) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ua := r.UserAgent()
		country := r.Header.Get("X-Appengine-Country")

		_, err := fmt.Fprintf(w, "hello my friend %s from %s :)", ua, country)
		if err != nil {
			log.Println("err:", err)
			return
		}

		data, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Println("err:", err)
			return
		}

		filename := time.Now().Format("2006-01-02")

		ctx := context.Background()
		var skipRead bool

		oattr, err := sc.Bucket(bucket).Object(filename).Attrs(ctx)
		if err != nil || oattr.Size > 52428800 {
			skipRead = true
		}

		sw := sc.Bucket(bucket).Object(filename).NewWriter(ctx)

		if !skipRead {
			sr, err := sc.Bucket(bucket).Object(filename).NewReader(ctx)
			if err != nil {
				log.Printf("Failed read from bucket '%s', file '%s': %v", bucket, filename, err)
			} else {
				var buf bytes.Buffer
				_, err := io.CopyN(sw, io.TeeReader(sr, &buf), 52428800) // 50MB
				if err != nil && err != io.EOF {
					log.Printf("Failed re-writing record: %v", err)
				}

				for {
					line, err := buf.ReadString('\n')
					if err != nil {
						if err == io.EOF {
							break
						}

						log.Fatalf("read file line error: %v", err)
						return
					}

					data := make(map[string]interface{})
					err = json.Unmarshal([]byte(line), &data)
					if err != nil {
						continue
					}

					if data["country"] == "AE" || data["country"] == nil {
						continue
					}

					_, err = fmt.Fprintf(w, p.Sanitize(fmt.Sprintf("\nsay hello to my old firend %v that visited at %v :)",
						data["ua"],
						data["timestamp"],
					)))

					if err != nil {
						log.Println("err:", err)
						return
					}
				}
			}
		}

		err = json.NewEncoder(sw).Encode(map[string]interface{}{
			"timestamp":   time.Now().Format(time.RFC3339),
			"remote_addr": r.RemoteAddr,
			"request":     string(data),
			"ua":          ua,
			"country":     country,
		})
		if err != nil {
			log.Printf("Failed writing record: %v", err)
			return
		}

		if err := sw.Close(); err != nil {
			log.Printf("Could not put file: %v", err)
			return
		}
	}
}
