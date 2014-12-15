package rest

import (
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
)

type App struct {
	Env      map[string]interface{}
	Handlers []*Handler
}

func (this *App) Mount(base string, app *App) {
	this.UsePath(base, func(req Request, res Response, next func()) {
		if !app.Exec(req, res, 0) {
			next()
		}
	})
}

func (this *App) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	request := Request{Req: req, App: this}
	request.Init()
	response := Response{Resp: res, App: this}
	this.Exec(request, response, 0)
}

func (this *App) Exec(request Request, response Response, i int) bool {
	log.Printf("Exec method:%s,path:%s,OriginPath:%s", request.Method, request.Path, request.OriginUrl())
	if len(this.Handlers) <= i {
		return false //no completed
	}
	handler := this.Handlers[i]
	base, params := handler.Matches(request.Method, request.Path)
	if 0 < len(base) {
		log.Printf("match ok ,base:%s,path:%s\n", request.Base, request.Path)
		request.Params = params
		if 1 < len(base) {
			request.Base = base
			request.Path = strings.TrimPrefix(path.Clean(request.Path), path.Clean(base))
			if 0 == len(request.Path) {
				request.Path = "/"
			}
		}
		log.Printf("matched ok ,base:%s,path:%s\n", request.Base, request.Path)
		handler.Handle(request, response, func() {
			if 1 < len(request.Base) {
				request.Path = path.Join(request.Base, request.Path)
				request.Base = "/"
			}
			//to do
			this.Exec(request, response, i+1)

		})

		return true
	} else {
		return this.Exec(request, response, i+1)
	}
}

func (this *App) Listen(port int) {
	log.Printf("server listen at:%d\n", port)
	http.ListenAndServe(":"+strconv.Itoa(port), this)
}

func (this *App) GetEnv(name string) interface{} {
	return this.Env[name]
}
func (this *App) SetEnv(name string, value interface{}) {
	this.Env[name] = value
}
func (this *App) Enable(name string) {
	this.Env[name] = true
}
func (this *App) Enabled(name string) bool {
	if r, ok := this.Env[name].(bool); ok {
		return r
	} else {
		return false
	}
}
func (this *App) Disable(name string) {
	this.Env[name] = false
}
func (this *App) Disabled(name string) bool {
	if r, ok := this.Env[name].(bool); ok {
		return r
	} else {
		return true
	}
}

func (this *App) Use(handle HandleFn) {
	this.UsePath("", handle)
}

func (this *App) UsePath(path string, handle HandleFn) {
	if len(path) > 0 && !strings.HasPrefix(path, "^") {
		path = "^" + path
	}
	this.RouteNext("", path, handle)
}

func (this *App) Get(path string, handle DoneFn) {
	this.Route("GET", path, handle)
}
func (this *App) Post(path string, handle DoneFn) {
	this.Route("POST", path, handle)
}
func (this *App) Delete(path string, handle DoneFn) {
	this.Route("DELETE", path, handle)
}
func (this *App) Put(path string, handle DoneFn) {
	this.Route("PUT", path, handle)
}
func (this *App) Patch(path string, handle DoneFn) {
	this.Route("PATCH", path, handle)
}
func (this *App) Route(method string, path string, handle DoneFn) {
	this.RouteNext(method, path, func(req Request, res Response, next func()) {
		handle(req, res)
	})
}

func (this *App) GetNext(path string, handle HandleFn) {
	this.RouteNext("GET", path, handle)
}
func (this *App) PostNext(path string, handle HandleFn) {
	this.RouteNext("POST", path, handle)
}
func (this *App) DeleteNext(path string, handle HandleFn) {
	this.RouteNext("DELETE", path, handle)
}
func (this *App) PutNext(path string, handle HandleFn) {
	this.RouteNext("PUT", path, handle)
}
func (this *App) PatchNext(path string, handle HandleFn) {
	this.RouteNext("PATCH", path, handle)
}
func (this *App) RouteNext(method string, path string, handle HandleFn) {
	log.Printf("method:%s,path:%s\n", method, path)
	this.Handlers = append(this.Handlers, &Handler{Method: strings.ToUpper(method), PathReg: PathToReg(path), Handle: handle})
}
