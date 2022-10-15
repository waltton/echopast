package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const qInsert = `
	INSERT INTO logs (timestamp, hash, url, method, host, user_agent, accept_encoding, accept, cookie, ip, protocol, headers, body)
	VALUES %s
	ON CONFLICT DO NOTHING
`

func writeLogs(db *sql.DB, logs []Log) (n int, err error) {
	args := []interface{}{}
	for _, l := range logs {
		headers, err := json.Marshal(l.Request.Headers)
		if err != nil {
			return n, err
		}

		var ip *string
		if l.Request.IP != nil {
			ipv := l.Request.IP.String()
			ip = &ipv
		}

		args = append(args,
			l.Timestamp,
			l.Hash,
			l.Request.URL,
			l.Request.Method,
			l.Request.Host,
			l.Request.UserAgent,
			l.Request.AcceptEncoding,
			l.Request.Accept,
			l.Request.Cookie,
			ip,
			l.Request.Protocol,
			headers,
			l.Request.Body,
		)

		if len(args) > 50000 {
			q := fmt.Sprintf(qInsert, buildParams(13, len(args)))

			r, err := db.Exec(q, args...)
			if err != nil {
				return n, err
			}

			ra, err := r.RowsAffected()
			if err != nil {
				return n, err
			}

			n += int(ra)

			args = []interface{}{}
		}
	}

	q := fmt.Sprintf(qInsert, buildParams(13, len(args)))

	r, err := db.Exec(q, args...)
	if err != nil {
		return n, err
	}

	ra, err := r.RowsAffected()
	if err != nil {
		return n, err
	}

	n += int(ra)

	return n, nil
}

func buildParams(cols, args int) (params string) {
	var sb strings.Builder

	for i := 1; i <= args/cols; i++ {
		if i > 1 {
			sb.WriteString(",")
		}
		sb.WriteString("(")
		for j := 1; j <= cols; j++ {
			if j > 1 {
				sb.WriteString(",")
			}
			sb.WriteString("$")
			sb.WriteString(strconv.Itoa(j + (i-1)*cols))
		}
		sb.WriteString(")")
	}

	return sb.String()
}
