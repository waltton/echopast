package whois

import (
	"bufio"
	"io"
	"log"
	"regexp"
	"strings"
)

type Whois []map[string][]string

var reKV = regexp.MustCompile(`([\w-]*):\s*(.*)`)

func isZeroString(value string) bool {
	bs := []byte(value)
	var total int64
	for _, b := range bs {
		total += int64(b)
	}
	return total == 0
}

func parseWhois(raw string) (w Whois, err error) {
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
		line = strings.Trim(line, "\u0000")

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
			if len(w[objc][prevk]) == 0 {
				log.Printf("no where to put '%s'; objc: %d, prevk: %s\n", line, objc, prevk)
			} else if strings.TrimSpace(w[objc][prevk][len(w[objc][prevk])-1]) == "" {
				log.Printf("last line for previous key is empty; objc: %d, prevk: %s\n", objc, prevk)
			} else {
				w[objc][prevk][len(w[objc][prevk])-1] = w[objc][prevk][len(w[objc][prevk])-1] + "\n" + strings.TrimSpace(line)
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

func (w Whois) Registry() string {
	if len(w) < 1 {
		return ""
	}

	val, ok := w[0]["%"]
	if ok {
		if len(val) > 0 {
			switch val[0] {
			case "IANA WHOIS server":
				return RegistryIANA
			case "This is the AfriNIC Whois server.":
				return RegistryAFRINIC
			case "Joint Whois - whois.lacnic.net":
				return RegistryLACNIC
			case "[whois.apnic.net]":
				fallthrough
			case "Whois data copyright terms    http://www.apnic.net/db/dbcopyright.html":
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

func (w Whois) Refer() string {
	switch w.Registry() {
	case RegistryIANA:
		if len(w) > 2 {
			val, ok := w[1]["refer"]
			if ok {
				if len(val) == 1 {
					return val[0]
				}
			}
		}
	}
	return ""
}

func (w Whois) Country() string {
	switch w.Registry() {
	case RegistryAPNIC:
		if len(w) > 4 {
			_, ok := w[3]["inetnum"]
			if ok {
				val, ok := w[3]["country"]
				if ok {
					if len(val) == 1 {
						return val[0]
					}
				}
			}
		}
	case RegistryARIN:
		if len(w) > 4 {
			val, ok := w[3]["Country"]
			if ok {
				if len(val) == 1 {
					return val[0]
				}
			}
		}
	case RegistryLACNIC:
		if len(w) > 4 {
			_, ok := w[3]["inetnum"]
			if ok {
				val, ok := w[3]["country"]
				if ok {
					if len(val) == 1 {
						return val[0]
					}
				}
			}
		}
	case RegistryAFRINIC:
		if len(w) > 5 {
			_, ok := w[4]["inetnum"]
			if ok {
				val, ok := w[4]["country"]
				if ok {
					if len(val) == 1 {
						return val[0]
					}
				}
			}
		}
	case RegistryRIPE:
		if len(w) > 1 {
			_, ok := w[0]["inetnum"]
			if ok {
				val, ok := w[0]["country"]
				if ok {
					if len(val) == 1 {
						return val[0]
					}
				}
			}
		}
	}
	return ""
}
