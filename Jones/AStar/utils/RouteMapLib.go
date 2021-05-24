package utils

import "strconv"

func SparkRoute(host string, port int, route string) string {
	return "http://" + host + ":" + strconv.Itoa(port) + "/" + route
}
