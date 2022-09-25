package whois

import (
	"io"
	"log"
	"net"
	"strings"
	"time"
)

func Lookup(param string) (result string, err error) {
	addr := "whois.iana.org"
	result, err = rawQuery(addr, param)
	if err != nil {
		return "", err
	}

	// fmt.Println("result", result)

	whois, err := parseWhois(result)
	if err != nil {
		return "", err
	}

	// data, err := json.Marshal(whois)
	// fmt.Println("data", string(data))
	// fmt.Println("err", err)

	if whois.Refer() == "" {
		return
	}

	result, err = rawQuery(whois.Refer(), param)
	if err != nil {
		return "", err
	}

	// fmt.Println("result", result)

	whois, err = parseWhois(result)
	if err != nil {
		return "", err
	}

	_ = whois

	// data, err := json.Marshal(whois)
	// fmt.Println("data", string(data))
	// fmt.Println("err", err)

	return
}

func rawQuery(addr, param string) (result string, err error) {
	conn, err := net.Dial("tcp", addr+":43")
	if err != nil {
		return "", err
	}

	if strings.Contains(addr, "arin") {
		param = "+ " + param
	}

	_, err = conn.Write([]byte(param + "\r\n"))
	if err != nil {
		return
	}

	for {
		conn.SetReadDeadline(time.Now().Add(time.Second * 30))

		data := make([]byte, 4096)
		n, err := conn.Read(data)
		if err != nil {
			if err == io.EOF {
				break
			} else if err.(net.Error).Timeout() {
				log.Print(err)
				break
			} else {
				log.Print(err)
			}
		}

		result += string(data[:n])
	}

	return
}
