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

func Shift(arr *[]string) string {
	if 0 == len(*arr) {
		return ""
	}
	r := (*arr)[0]
	*arr = (*arr)[1:]
	return r
}

func IsUrl(url string) {

}
