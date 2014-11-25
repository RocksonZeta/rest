package utils

import (
	"strings"
)

func ParseQueryString(queryString string) map[string][]string {
	var kvStrings = strings.Split(queryString, "&")

	result = make(map[string][]string, len(kvStrings))
	for _, kvString := range kvStrings {
		kvs := strings.Split(kvString, "=")
		if nil == result[kvs[0]] {
			if 2 <= len(kvs) {
				result[kvs[0]] = []string{kv[1]}
			} else {
				result[kvs[0]] = []string{""}
			}
		} else {
			if 2 <= len(kvs) {
				append(result[kvs[0]], kv[1])
			} else {
				append(result[kvs[0]], "")
			}
		}
	}
}
