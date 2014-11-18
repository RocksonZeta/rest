package rest

import (
	"log"
	"regexp"
	"strings"
)

type Handler struct {
	path   *regexp.Regexp
	method string
	handle func(req Request, res Response, next func(e error))
}

func (this *Handler) Matches(method, path string) bool {
	log.Printf("rest#Matches - this.method:%s,this.path:%s to match %s %s \n", this.method, this.path, method, path)
	if 0 == len(this.method) && nil == this.path {
		return true
	}
	if 0 != len(this.method) && nil == this.path {
		return this.method == strings.ToUpper(method)
	}
	if nil != this.path {
		return this.path.MatchString(path)
	}

	return true
}
