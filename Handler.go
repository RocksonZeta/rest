package rest

import (
	"regexp"
	"strings"
)

type HandleFn func(req *Request, res *Response, next func())
type DoneFn func(req *Request, res *Response)

type Handler struct {
	PathReg *regexp.Regexp
	Method  string
	Handle  HandleFn
}

func (this *Handler) Matches(method, path string) (base string, params map[string]string) {
	if 0 != len(this.Method) {
		if this.Method != strings.ToUpper(method) {
			return
		}
		if nil != this.PathReg {
			base, params = NamedMatches(this.PathReg, path)
		} else {
			base = "/"
		}
	} else {
		if nil == this.PathReg {
			base = "/"
			return
		}
		base, params = NamedMatches(this.PathReg, path)
	}
	return

}
