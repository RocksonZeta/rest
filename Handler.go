package rest

import (
	"log"
	"regexp"
	"strings"
)

type Handler struct {
	path    string
	pathReg *regexp.Regexp
	method  string
	handle  func(req *Request, res *Response, next func(e error))
}

func (this *Handler) Matches(method, path string) (params map[string]string, ok bool) {
	log.Printf("Handle#Matches - this.method:%s,this.path:%s to match %s %s \n", this.method, this.path, method, path)
	if nil == this.pathReg {
		this.pathReg = PathToReg(this.path)
	}
	if 0 == len(this.method) && nil == this.pathReg {
		ok = true
		return
	}
	if 0 != len(this.method) && nil == this.pathReg {
		ok = (this.method == strings.ToUpper(method))
		return
	}
	if nil != this.pathReg {
		params = Matches(this.pathReg, path)
		if nil != params {
			ok = true
		}
		return
	}
	ok = true
	return
}
