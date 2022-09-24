package whois

import (
	"bufio"
	"io"
	"log"
	"strings"
)

type WhoisGEN []map[string][]string

func parseWhoisGEN(raw string) (w WhoisGEN, err error) {
	var i int
	var objc int
	var prevk string

	w = append(w, make(map[string][]string))
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
			if len(w[objc]) > 0 {
				objc++
				w = append(w, make(map[string][]string))
			}

			continue
		}

		if line[0] == '#' || line[0] == '%' {
			w[objc][string(line[0])] = append(w[objc][string(line[0])], strings.Trim(line, " #%"))
			continue
		}

		m := reKV.FindStringSubmatch(line)
		if len(m) == 3 {
			w[objc][m[1]] = append(w[objc][m[1]], m[2])
			prevk = m[1]
		} else {
			if strings.TrimSpace(w[objc][prevk][len(w[objc][prevk])-1]) != "" {
				w[objc][prevk][len(w[objc][prevk])-1] = w[objc][prevk][len(w[objc][prevk])-1] + "\n" + strings.TrimSpace(line)
			} else {
				log.Printf("last line for previous key is empty; objc: %d, prevk: %s\n", objc, prevk)
			}
		}

	}

	return
}

const (
	RegistryIANA    = "iana"
	RegistryARIN    = "arin"
	RegistryAPNIC   = "apnic"
	RegistryLACNIC  = "lacnic"
	RegistryAFRINIC = "afrinic"
	RegistryRIPE    = "ripe"
)

func (w WhoisGEN) Registry() string {
	if len(w) < 1 {
		return ""
	}

	val, ok := w[0]["%"]
	if ok {
		if len(val) > 0 {
			if val[0] == "IANA WHOIS server" {
				return RegistryIANA
			}
			if val[0] == "This is the AfriNIC Whois server." {
				return RegistryAFRINIC
			}
			if val[0] == "Joint Whois - whois.lacnic.net" {
				return RegistryLACNIC
			}
			if val[0] == "Whois data copyright terms    http://www.apnic.net/db/dbcopyright.html" {
				return RegistryAPNIC
			}
		}
	}

	val, ok = w[0]["#"]
	if ok {
		if len(val) > 1 {
			if val[1] == "ARIN WHOIS data and services are subject to the Terms of Use" {
				return RegistryARIN
			}
		}
	}

	val, ok = w[0]["source"]
	if ok {
		if len(val) > 0 {
			if val[0] == "RIPE" {
				return RegistryRIPE
			}
		}
	}

	return ""
}