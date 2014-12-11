package rest

import (
	"log"
	"strings"
)

type Router struct {
	Handlers []*Handler
}

func (this *Router) Use(handle HandleFn) {
	this.UsePath("", handle)
}

func (this *Router) UsePath(path string, handle HandleFn) {
	if len(path) > 0 && !strings.HasPrefix(path, "^") {
		path = "^" + path
	}
	this.RouteNext("", path, handle)
}

func (this *Router) Get(path string, handle DoneFn) {
	this.Route("GET", path, handle)
}
func (this *Router) Post(path string, handle DoneFn) {
	this.Route("POST", path, handle)
}
func (this *Router) Delete(path string, handle DoneFn) {
	this.Route("DELETE", path, handle)
}
func (this *Router) Put(path string, handle DoneFn) {
	this.Route("PUT", path, handle)
}
func (this *Router) Patch(path string, handle DoneFn) {
	this.Route("PATCH", path, handle)
}
func (this *Router) Route(method string, path string, handle DoneFn) {
	this.RouteNext(method, path, func(req *Request, res *Response, next func(e error)) {
		handle(req, res)
	})
}

func (this *Router) GetNext(path string, handle HandleFn) {
	this.RouteNext("GET", path, handle)
}
func (this *Router) PostNext(path string, handle HandleFn) {
	this.RouteNext("POST", path, handle)
}
func (this *Router) DeleteNext(path string, handle HandleFn) {
	this.RouteNext("DELETE", path, handle)
}
func (this *Router) PutNext(path string, handle HandleFn) {
	this.RouteNext("PUT", path, handle)
}
func (this *Router) PatchNext(path string, handle HandleFn) {
	this.RouteNext("PATCH", path, handle)
}
func (this *Router) RouteNext(method string, path string, handle HandleFn) {
	log.Printf("method:%s,path:%s\n", method, path)
	this.Handlers = append(this.Handlers, &Handler{Method: strings.ToUpper(method), PathReg: PathToReg(path), Handle: handle})
}
