package rest

import (
	"log"
	pathUtils "path"
	"strings"
)

type Router struct {
	Handlers []Handler
}

func (this *Router) Mount(base string, mount *Router) {
	if nil == mount {
		return
	}
	for _, handler := range mount.Handlers {
		handler.Path = pathUtils.Join(base, handler.Path)
		this.Handlers = append(this.Handlers, handler)
	}
}

func (this *Router) Get(path string, handle func(req *Request, res *Response)) {
	this.Route("GET", path, handle)
}
func (this *Router) Post(path string, handle func(req *Request, res *Response)) {
	this.Route("POST", path, handle)
}
func (this *Router) Delete(path string, handle func(req *Request, res *Response)) {
	this.Route("DELETE", path, handle)
}
func (this *Router) Put(path string, handle func(req *Request, res *Response)) {
	this.Route("PUT", path, handle)
}
func (this *Router) Patch(path string, handle func(req *Request, res *Response)) {
	this.Route("PATCH", path, handle)
}
func (this *Router) Route(method string, path string, handle func(req *Request, res *Response)) {
	this.RouteNext(method, path, func(req *Request, res *Response, next func(e error)) {
		handle(req, res)
	})
}

func (this *Router) GetNext(path string, handle func(req *Request, res *Response, next func(e error))) {
	this.RouteNext("GET", path, handle)
}
func (this *Router) PostNext(path string, handle func(req *Request, res *Response, next func(e error))) {
	this.RouteNext("POST", path, handle)
}
func (this *Router) DeleteNext(path string, handle func(req *Request, res *Response, next func(e error))) {
	this.RouteNext("DELETE", path, handle)
}
func (this *Router) PutNext(path string, handle func(req *Request, res *Response, next func(e error))) {
	this.RouteNext("PUT", path, handle)
}
func (this *Router) PatchNext(path string, handle func(req *Request, res *Response, next func(e error))) {
	this.RouteNext("PATCH", path, handle)
}
func (this *Router) RouteNext(method string, path string, handle func(req *Request, res *Response, next func(e error))) {
	log.Printf("method:%s,path:%s\n", method, path)
	this.Handlers = append(this.Handlers, Handler{Method: strings.ToUpper(method), Path: path, Handle: handle})
}
