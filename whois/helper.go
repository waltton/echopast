package whois

import "regexp"

var reKV = regexp.MustCompile(`([\w-]*):\s*(.*)`)
var reKVo = regexp.MustCompile(`([\w-]*):\s*(.*)?`)

func isZeroString(value string) bool {
	bs := []byte(value)
	var total int64
	for _, b := range bs {
		total += int64(b)
	}
	return total == 0
}
