package whois

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"strings"
	"time"
)

func Lookup(param string) (result *Whois, err error) {
	addr := "whois.iana.org"

	rr, err := rawQuery(addr, param)
	if err != nil {
		return nil, err
	}

	whois, err := parseWhois(rr)
	if err != nil {
		return nil, err
	}

	if whois.Refer() == "" {
		return
	}

	rr, err = rawQuery(whois.Refer(), param)
	if err != nil {
		return nil, err
	}

	whois, err = parseWhois(rr)
	if err != nil {
		return nil, err
	}

	// data, err := json.Marshal(whois)
	// fmt.Println("data", string(data))
	// fmt.Println("err", err)

	return &whois, nil
}

func rawQuery(addr, param string) (result string, err error) {
	result, err = getFromCache(addr, param)
	if err == nil {
		log.Print("hit", addr, param)
	} else {
		log.Print("miss", addr, param, err)
	}

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

	err = writeToCache(addr, param, result)
	if err != nil {
		log.Print(err)
	}

	return
}

func getFromCache(addr, param string) (result string, err error) {
	cacheFile := path.Join("/Users/waltton/projects/echopast/cache", fmt.Sprintf("%s-%s", param, addr))
	f, err := os.Open(cacheFile)
	if err != nil {
		return "", err
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func writeToCache(addr, param, result string) (err error) {
	cacheFile := path.Join("/Users/waltton/projects/echopast/cache", fmt.Sprintf("%s-%s", param, addr))

	f, err := os.Create(cacheFile)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(f, result)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return
}
