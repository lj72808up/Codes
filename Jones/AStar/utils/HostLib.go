package utils

import "strings"

func GetHostFromAddr(addr string) string {
	if !strings.Contains(addr, ":") {
		return ""
	} else {
		return strings.Split(addr, ":")[0]
	}
}