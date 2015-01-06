package rest

import (
	"regexp"
	"strconv"
	"strings"
)

var PATH_REG *regexp.Regexp = regexp.MustCompile(":([^/]+)")

func NamedPath(path string) string {
	return PATH_REG.ReplaceAllString(path, "(?P<$1>[^/]+)")
}

//convert request path to regular string
func PathToRegString(path string) string {
	if 0 == len(path) {
		return path
	}
	if strings.HasPrefix(path, "^") {
		return "(?i)" + NamedPath(path)
	}
	return "(?i)^" + NamedPath(path) + "$"
}

//convert path to regular expression
func PathToReg(path string) *regexp.Regexp {
	if 0 == len(path) {
		return nil
	}
	return regexp.MustCompile(PathToRegString(path))
}

//get name group from matched result
func NamedMatches(reg *regexp.Regexp, path string) (base string, result map[string]string) {
	values := reg.FindStringSubmatch(path)
	if 0 == len(values) {
		return
	}
	base = values[0]
	keys := reg.SubexpNames()
	result = make(map[string]string, len(keys)-1)
	for i := 1; i < len(keys); i++ {
		if "" == keys[i] {
			result[strconv.Itoa(i)] = values[i]
		} else {
			result[keys[i]] = values[i]
		}

	}
	return
}
