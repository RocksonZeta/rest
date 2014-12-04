package rest

import (
	"regexp"
	"strings"
)

type Handler struct {
	Path    string
	PathReg *regexp.Regexp
	Method  string
	Handle  func(req *Request, res *Response, next func(e error))
}

func (this *Handler) Matches(method, path string) (params map[string]string, ok bool) {
	if nil == this.PathReg {
		this.PathReg = PathToReg(this.Path)
	}
	if 0 == len(this.Method) && nil == this.PathReg {
		ok = true
		return
	}
	if 0 != len(this.Method) && nil == this.PathReg {
		ok = (this.Method == strings.ToUpper(method))
		return
	}
	if nil != this.PathReg {
		params = Matches(this.PathReg, path)
		if nil != params {
			ok = true
		}
		return
	}
	ok = true
	return
}
