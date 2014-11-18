package rest

import (
	"log"
	"strings"
)

type Mount struct {
	base     string
	handlers []Handler
}

func (this *Mount) Get(path string, handle func(req Request, res Response)) {
	this.MethodNext("GET", path, func(req Request, res Response, next func(e error)) {
		handle(req, res)
	})
}
func (this *Mount) Post(path string, handle func(req Request, res Response)) {
	this.Method("POST", path, handle)
}
func (this *Mount) Delete(path string, handle func(req Request, res Response)) {
	this.Method("DELETE", path, handle)
}
func (this *Mount) Put(path string, handle func(req Request, res Response)) {
	this.Method("PUT", path, handle)
}
func (this *Mount) Patch(path string, handle func(req Request, res Response)) {
	this.Method("PATCH", path, handle)
}
func (this *Mount) Method(method string, path string, handle func(req Request, res Response)) {
}

func (this *Mount) GetNext(path string, handle func(req Request, res Response, next func(e error))) {
	this.MethodNext("GET", path, handle)
}
func (this *Mount) PostNext(path string, handle func(req Request, res Response, next func(e error))) {
	this.MethodNext("POST", path, handle)
}
func (this *Mount) DeleteNext(path string, handle func(req Request, res Response, next func(e error))) {
	this.MethodNext("DELETE", path, handle)
}
func (this *Mount) PutNext(path string, handle func(req Request, res Response, next func(e error))) {
	this.MethodNext("PUT", path, handle)
}
func (this *Mount) PatchNext(path string, handle func(req Request, res Response, next func(e error))) {
	this.MethodNext("PATCH", path, handle)
}
func (this *Mount) MethodNext(method string, path string, handle func(req Request, res Response, next func(e error))) {
	this.handlers = append(this.handlers, Handler{method: strings.ToUpper(method), path: PathToReg(path), handle: handle})
	log.Printf("methodNext handles len:%d\n", len(this.handlers))
}
