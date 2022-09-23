package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

const qInsert = `
	INSERT INTO logs (timestamp, url, host, user_agent, accept_encoding, accept, cookie, ip, protocol, headers)
	VALUES %s
`

func writeLogs(db *sql.DB, logs []Log) error {
	args := []interface{}{}
	for _, l := range logs {
		args = append(args,
			l.Timestamp,
			l.Request.URL,
			l.Request.Host,
			l.Request.UserAgent,
			l.Request.AcceptEncoding,
			l.Request.Accept,
			l.Request.Cookie,
			l.Request.IP,
			l.Request.Protocol,
			l.Request.Headers,
		)
	}

	q := fmt.Sprintf(qInsert, buildParams(10, len(logs)))

	_, err := db.Exec(q, args...)
	if err != nil {
		return err
	}

	return nil
}

func buildParams(cols, rows int) (params string) {
	var sb strings.Builder

	for i := 1; i <= rows; i++ {
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
