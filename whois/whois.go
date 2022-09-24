package whois

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func Lookup(param string) (result string, err error) {
	addr := "whois.iana.org:43"
	result, err = rawQuery(addr, param)
	if err != nil {
		return "", err
	}

	return
}

func rawQuery(addr, param string) (result string, err error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return "", err
	}

	n, err := conn.Write([]byte(param + "\r\n"))
	if err != nil {
		return
	}

	fmt.Println("n", n)

	for {
		conn.SetReadDeadline(time.Now().Add(time.Second * 30))

		data := make([]byte, 10)
		// data := make([]byte, 4096)
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
