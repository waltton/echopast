package whois

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
)

type WhoisIANA struct {
	Comments []string `json:"comments"`

	Refer        string `json:"refer"`
	InetNum      string `json:"inetnum"`
	InetNumBegin net.IP `json:"inetnum_begin"`
	InetNumEnd   net.IP `json:"inetnum_end"`
	Organisation string `json:"organisation"`
	Status       string `json:"status"`
	Whois        string `json:"whois"`
	Changed      string `json:"changed"`
	Source       string `json:"source"`
}

func parseWhoisIANA(raw string) (w *WhoisIANA, err error) {
	w = new(WhoisIANA)

	var i int
	rd := bufio.NewReader(strings.NewReader(raw))
	for {
		i++

		line, err := rd.ReadString('\n')
		if isZeroString(line) && err != nil {
			if err == io.EOF {
				break
			}

			log.Fatalf("read file line error: %v", err)
			break
		}

		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "%") {
			w.Comments = append(w.Comments, line)
			continue
		}

		m := reKV.FindStringSubmatch(line)
		if len(m) == 3 {
			switch m[1] {
			case "refer":
				w.Refer = m[2]
			case "inetnum":
				w.InetNum = m[2]
				ins := strings.Split(w.InetNum, " - ")
				w.InetNumBegin = net.ParseIP(ins[0])
				w.InetNumEnd = net.ParseIP(ins[1])
			case "organisation":
				w.Organisation = m[2]
			case "status":
				w.Status = m[2]
			case "whois":
				w.Whois = m[2]
			case "changed":
				w.Changed = m[2]
			case "source":
				w.Source = m[2]
			}
		} else {
			log.Printf("fail to parse line: #%d '%+v'", i, line)
		}

	}

	return
}
