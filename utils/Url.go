package utils

import (
	"net/url"
	"strings"
)

func ParseQueryString(queryString string) map[string][]string {
	var kvStrings = strings.Split(queryString, "&")

	result := make(map[string][]string, len(kvStrings))
	for _, kvString := range kvStrings {
		kvs := strings.Split(kvString, "=")
		if nil == result[kvs[0]] {
			if 2 <= len(kvs) {
				v, _ := url.QueryUnescape(kvs[1])
				result[kvs[0]] = []string{v}
			} else {
				result[kvs[0]] = []string{""}
			}
		} else {
			if 2 <= len(kvs) {
				v, _ := url.QueryUnescape(kvs[1])
				result[kvs[0]] = append(result[kvs[0]], v)
			} else {
				result[kvs[0]] = append(result[kvs[0]], "")
			}
		}
	}
	return result
}
