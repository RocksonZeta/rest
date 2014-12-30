package rest

import (
	"path"
	"strings"
)

type Routable interface {
	//return if the app completed the request,true completed,false uncompleted
	Exec(req Request, res Response, i int) bool
}

type Router struct {
	Handlers []*Handler
}

func (this *Router) Mount(base string, app Routable) {
	this.UsePath(base, func(req Request, res Response, next func()) {
		if !app.Exec(req, res, 0) {
			next()
		}
	})
}

func (this *Router) Exec(request Request, response Response, i int) bool {
	//log.Printf("Exec method:%s,path:%s,OriginPath:%s", request.Method, request.Path, request.OriginUrl())
	if len(this.Handlers) <= i {
		return false //no completed
	}
	handler := this.Handlers[i]
	base, params := handler.Matches(request.Method, request.Path)
	if 0 < len(base) {
		//log.Printf("match ok ,base:%s,path:%s\n", request.Base, request.Path)
		request.Params = params
		var oldPath = request.Path
		if 1 < len(base) {
			request.Base = base
			request.Path = strings.TrimPrefix(path.Clean(request.Path), path.Clean(base))
			if 0 == len(request.Path) {
				request.Path = "/"
			}
		}
		//var hp string
		//if nil != handler.PathReg {
		//	hp = handler.PathReg.String()
		//}
		//log.Printf("matched ok ,base:%s,path:%s,handler path:%s\n", request.Base, request.Path, handler.PathReg)
		handler.Handle(request, response, func() {
			if 1 < len(request.Base) {
				request.Path = oldPath
				request.Base = "/"
			}
			this.Exec(request, response, i+1)

		})

		return true
	} else {
		return this.Exec(request, response, i+1)
	}
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
	this.RouteNext(method, path, func(req Request, res Response, next func()) {
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
	this.Handlers = append(this.Handlers, &Handler{Method: strings.ToUpper(method), PathReg: PathToReg(path), Handle: handle})
}
