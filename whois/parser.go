package whois

import (
	"bufio"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type Whois struct {
	Data []map[string][]string
}

var reKV = regexp.MustCompile(`([\w-]*):\s*(.*)`)

func parseWhois(raw string) (w Whois, err error) {
	var i int
	var objc int
	var prevk string

	w.Data = append(w.Data, make(map[string][]string))
	rd := bufio.NewReader(strings.NewReader(raw))
	for {
		i++

		line, err := rd.ReadString('\n')
		if line == "" && err != nil {
			if err == io.EOF {
				break
			}

			log.Fatalf("read file line error: %v", err)
			break
		}

		line = strings.TrimSpace(line)

		if line == "" {
			if len(w.Data[objc]) > 0 {
				objc++
				w.Data = append(w.Data, make(map[string][]string))
			}

			continue
		}

		if line[0] == '#' || line[0] == '%' {
			w.Data[objc][string(line[0])] = append(w.Data[objc][string(line[0])], strings.Trim(line, " #%"))
			continue
		}

		m := reKV.FindStringSubmatch(line)
		if len(m) == 3 {
			w.Data[objc][m[1]] = append(w.Data[objc][m[1]], m[2])
			prevk = m[1]
		} else {
			if len(w.Data[objc][prevk]) == 0 {
				log.Printf("no w.Datahere to put '%s'; objc: %d, prevk: %s\n", line, objc, prevk)
			} else if strings.TrimSpace(w.Data[objc][prevk][len(w.Data[objc][prevk])-1]) == "" {
				log.Printf("last line for previous key is empty; objc: %d, prevk: %s\n", objc, prevk)
			} else {
				w.Data[objc][prevk][len(w.Data[objc][prevk])-1] = w.Data[objc][prevk][len(w.Data[objc][prevk])-1] + "\n" + strings.TrimSpace(line)
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
	if len(w.Data) < 1 {
		return ""
	}

	val, ok := w.Data[0]["%"]
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
			case "This is the RIPE Database query service.":
				return RegistryRIPE
			}
		}
	}

	val, ok = w.Data[0]["#"]
	if ok {
		if len(val) > 1 {
			if val[1] == "ARIN WHOIS data and services are subject to the Terms of Use" {
				return RegistryARIN
			}
		}
	}

	val, ok = w.Data[0]["source"]
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
		if len(w.Data) > 2 {
			val, ok := w.Data[1]["refer"]
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
		if len(w.Data) > 4 {
			_, ok := w.Data[3]["inetnum"]
			if ok {
				val, ok := w.Data[3]["country"]
				if ok {
					if len(val) == 1 {
						return val[0]
					}
				}
			}
		}
	case RegistryARIN:
		begin, end := w.arinGetBlock()
		for _, obj := range w.Data[begin:end] {
			val, ok := obj["Country"]
			if ok {
				if len(val) == 1 {
					return val[0]
				}
			}
		}
	case RegistryLACNIC:
		if len(w.Data) > 4 {
			_, ok := w.Data[3]["inetnum"]
			if ok {
				val, ok := w.Data[3]["country"]
				if ok {
					if len(val) == 1 {
						return val[0]
					}
				}
			}
		}
	case RegistryAFRINIC:
		if len(w.Data) > 5 {
			_, ok := w.Data[4]["inetnum"]
			if ok {
				val, ok := w.Data[4]["country"]
				if ok {
					if len(val) == 1 {
						return val[0]
					}
				}
			}
		}
	case RegistryRIPE:
		begin, end := w.ripeGetBlock()

		for _, obj := range w.Data[begin:end] {
			_, ok := obj["inetnum"]
			if ok {
				val, ok := obj["country"]
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

func (w Whois) arinGetBlock() (begin, end int) {
	var starts, ends []int
	bnc, enc := -1, -1
	for i := range w.Data {
		val, ok := w.Data[i]["#"]
		if ok {
			if len(val) == 1 {
				if val[0] == "start" {
					starts = append(starts, i)
				} else if val[0] == "end" {
					ends = append(ends, i)
				}
			}
		} else {
			if bnc == -1 {
				bnc = i
			}
			enc = i
		}
	}

	if len(starts) == 0 {
		return bnc, enc
	}

	var cidrsSize []int
	for _, sidx := range starts {
		if len(w.Data) > sidx+1 {
			val, ok := w.Data[sidx+1]["CIDR"]
			if ok {
				if len(val) == 1 {
					sval := strings.Split(val[0], "/")
					if len(sval) == 2 {
						size, err := strconv.Atoi(sval[1])
						if err == nil {
							cidrsSize = append(cidrsSize, size)
						}
					}
				}
			}
		}
	}

	var maxValIdx int
	for i, val := range cidrsSize {
		if val > cidrsSize[maxValIdx] {
			maxValIdx = i
		}
	}

	return starts[maxValIdx] + 1, ends[maxValIdx] - 1
}

func (w Whois) ripeGetBlock() (begin, end int) {
	begin, end = -1, -1
	for i := range w.Data {
		_, ok := w.Data[i]["%"]
		if !ok {
			if begin == -1 {
				begin = i
			}
			end = i
		}
	}
	return
}
