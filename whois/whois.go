package whois

import (
	"bufio"
	"io"
	"log"
	"net"
	"regexp"
	"strings"
	"time"
)

type Whois struct {
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

func Lookup(param string) (result string, err error) {
	addr := "whois.iana.org:43"
	result, err = rawQuery(addr, param)
	if err != nil {
		return "", err
	}

	whois, err := parseWhois(result)
	if err != nil {
		return "", err
	}

	_ = whois

	return
}

func rawQuery(addr, param string) (result string, err error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return "", err
	}

	_, err = conn.Write([]byte(param + "\r\n"))
	if err != nil {
		return
	}

	for {
		conn.SetReadDeadline(time.Now().Add(time.Second * 30))

		data := make([]byte, 4096)
		_, err := conn.Read(data)
		if err != nil {
			if err == io.EOF || err.(net.Error).Timeout() {
				log.Print(err)
				break
			} else {
				log.Print(err)
			}
		}

		result += string(data)
	}

	return
}

var reKV = regexp.MustCompile(`(\w*):\s*(.*)`)

func parseWhois(raw string) (w *Whois, err error) {
	w = new(Whois)

	var i int
	scanner := bufio.NewScanner(strings.NewReader(raw))
	for scanner.Scan() {
		i++

		if scanner.Text() == "" {
			continue
		}

		if strings.HasPrefix(scanner.Text(), "%") {
			w.Comments = append(w.Comments, scanner.Text())
			continue
		}

		m := reKV.FindStringSubmatch(scanner.Text())
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
			log.Printf("fail to parse line:'%s'", scanner.Text())
		}

	}

	err = scanner.Err()

	return
}
