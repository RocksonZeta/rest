package rest

import (
	"regexp"
	"strconv"
)

var PATH_REG *regexp.Regexp = regexp.MustCompile(":([^/]+)")

func PathToRegString(path string) string {
	if 0 == len(path) {
		return path
	}
	var reg = PATH_REG.ReplaceAllString(path, "(?P<$1>[^/]+)")
	return "^" + reg + "$"
}
func PathToReg(path string) *regexp.Regexp {
	return regexp.MustCompile(PathToRegString(path))
}

func Matches(reg *regexp.Regexp, path string) map[string]string {
	if !reg.MatchString(path) {
		return nil
	}
	var keys = reg.SubexpNames()
	var values = reg.FindStringSubmatch(path)
	result := make(map[string]string, len(keys)-1)
	for i := 1; i < len(keys); i++ {
		if "" == keys[i] {
			result[strconv.Itoa(i)] = values[i]
		} else {
			result[keys[i]] = values[i]
		}

	}
	return result
}
