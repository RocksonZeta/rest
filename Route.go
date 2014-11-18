package rest

import (
	pathUtils "path"
	"strings"
)

type Route struct {
	handlers []Handler
}

func (this *Route) Mount(base string, mount *Route) {
	if nil == mount {
		return
	}
	for _, handler := range mount.handlers {
		handler.path = pathUtils.Join(base, handler.path)
		this.handlers = append(this.handlers, handler)
	}
}

func (this *Route) Get(path string, handle func(req Request, res Response)) {
	this.Method("GET", path, handle)
}
func (this *Route) Post(path string, handle func(req Request, res Response)) {
	this.Method("POST", path, handle)
}
func (this *Route) Delete(path string, handle func(req Request, res Response)) {
	this.Method("DELETE", path, handle)
}
func (this *Route) Put(path string, handle func(req Request, res Response)) {
	this.Method("PUT", path, handle)
}
func (this *Route) Patch(path string, handle func(req Request, res Response)) {
	this.Method("PATCH", path, handle)
}
func (this *Route) Method(method string, path string, handle func(req Request, res Response)) {
	this.MethodNext(method, path, func(req Request, res Response, next func(e error)) {
		handle(req, res)
	})
}

func (this *Route) GetNext(path string, handle func(req Request, res Response, next func(e error))) {
	this.MethodNext("GET", path, handle)
}
func (this *Route) PostNext(path string, handle func(req Request, res Response, next func(e error))) {
	this.MethodNext("POST", path, handle)
}
func (this *Route) DeleteNext(path string, handle func(req Request, res Response, next func(e error))) {
	this.MethodNext("DELETE", path, handle)
}
func (this *Route) PutNext(path string, handle func(req Request, res Response, next func(e error))) {
	this.MethodNext("PUT", path, handle)
}
func (this *Route) PatchNext(path string, handle func(req Request, res Response, next func(e error))) {
	this.MethodNext("PATCH", path, handle)
}
func (this *Route) MethodNext(method string, path string, handle func(req Request, res Response, next func(e error))) {
	this.handlers = append(this.handlers, Handler{method: strings.ToUpper(method), path: path, handle: handle})
}
