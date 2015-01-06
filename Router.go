package rest

import (
	"path"
	"strings"
)

//Routable mean this object can used by Router Mount method
//We can construct ours Routable unit ,such as UserRouter,xxxServiceRouter .etc
type Routable interface {
	//return if the app completed the request,true completed,false uncompleted
	Exec(req Request, res Response, i int) bool
}

//Basic Routable implement.
type Router struct {
	Handlers []*Handler
}

//mount a Routable object at specifed mounting point.
//base is the mount point
//router is the router unit,
//if the router mount at base , the path of the router will be the trimed.
func (this *Router) Mount(base string, router Routable) {
	this.UsePath(base, func(req Request, res Response, next func()) {
		if !router.Exec(req, res, 0) {
			next()
		}
	})
}

//Execute the handler of the Router one by one , if the request handled , execute terminated.
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

//add middleware to router, it will get all request at this point
func (this *Router) Use(handle HandleFn) {
	this.UsePath("", handle)
}

//add middleware at specified point to router,it will dispose the request will path prefix
func (this *Router) UsePath(path string, handle HandleFn) {
	if len(path) > 0 && !strings.HasPrefix(path, "^") {
		path = "^" + path
	}
	this.RouteNext("", path, handle)
}

//add "Get" method request handler
func (this *Router) Get(path string, handle DoneFn) {
	this.Route("GET", path, handle)
}

//add "Post" method request handler
func (this *Router) Post(path string, handle DoneFn) {
	this.Route("POST", path, handle)
}

//add "Delete" method request handler
func (this *Router) Delete(path string, handle DoneFn) {
	this.Route("DELETE", path, handle)
}

//add "Put" method request handler
func (this *Router) Put(path string, handle DoneFn) {
	this.Route("PUT", path, handle)
}

//add "Patch" method request handler
func (this *Router) Patch(path string, handle DoneFn) {
	this.Route("PATCH", path, handle)
}

//add request handler
//method - request method
//path - request path,
//handler - a handler
func (this *Router) Route(method string, path string, handle DoneFn) {
	this.RouteNext(method, path, func(req Request, res Response, next func()) {
		handle(req, res)
	})
}

//add "Get" request middle at spcified path
func (this *Router) GetNext(path string, handle HandleFn) {
	this.RouteNext("GET", path, handle)
}

//add "Post" request middle at spcified path
func (this *Router) PostNext(path string, handle HandleFn) {
	this.RouteNext("POST", path, handle)
}

//add "Delete" request middle at spcified path
func (this *Router) DeleteNext(path string, handle HandleFn) {
	this.RouteNext("DELETE", path, handle)
}

//add "Put" request middle at spcified path
func (this *Router) PutNext(path string, handle HandleFn) {
	this.RouteNext("PUT", path, handle)
}

//add "Patch" request middle at spcified path
func (this *Router) PatchNext(path string, handle HandleFn) {
	this.RouteNext("PATCH", path, handle)
}

//add a user defined middleware
func (this *Router) RouteNext(method string, path string, handle HandleFn) {
	this.Handlers = append(this.Handlers, &Handler{Method: strings.ToUpper(method), PathReg: PathToReg(path), Handle: handle})
}
