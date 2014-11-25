package rest

import (
	"log"
	"net/http"
	"strconv"
)

type App struct {
	Route
	env map[string]interface{}
	//handlers []Handler
}

func (this *App) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	request := &Request{app: this, req: req, method: req.Method}
	response := &Response{app: this, res: res}
	this.exec(request, response, 0)
}

func (this *App) exec(request *Request, response *Response, i int) {
	if len(this.handlers) <= i {
		return
	}
	handler := this.handlers[i]
	if handler.Matches(request.Method(), request.Path()) {
		handler.handle(request, response, func(e error) {
			if nil != e {
				panic(e.Error())
				return
			}
			this.exec(request, response, i+1)
		})
	} else {
		this.exec(request, response, i+1)
	}
}

func (this *App) Listen(port int) {
	log.Printf("server listen at:%d\n", port)

	http.ListenAndServe(":"+strconv.Itoa(port), this)
}

func (this *App) Use(handle func(req *Request, res *Response, next func(e error))) {
	this.UsePath("", handle)
}

func (this *App) UsePath(path string, handle func(req *Request, res *Response, next func(e error))) {
	this.MethodNext("", path, handle)
}

func (this *App) GetEnv(name string) string {
	return ""
}
func (this *App) SetEnv(name string, value interface{}) {

}
func (this *App) Enable(name string) {

}
func (this *App) Enabled(name string) bool {
	return false
}
func (this *App) Disable(name string) {

}
func (this *App) Disabled(name string) bool {
	return false
}
func (this *App) Get(path string, handle func(req *Request, res *Response)) {
	this.MethodNext("GET", path, func(req *Request, res *Response, next func(e error)) {
		handle(req, res)
	})
}

//func (this *App) Post(path string, handle func(req Request, res Response)) {
//	this.Method("POST", path, handle)
//}
//func (this *App) Delete(path string, handle func(req Request, res Response)) {
//	this.Method("DELETE", path, handle)
//}
//func (this *App) Put(path string, handle func(req Request, res Response)) {
//	this.Method("PUT", path, handle)
//}
//func (this *App) Patch(path string, handle func(req Request, res Response)) {
//	this.Method("PATCH", path, handle)
//}
//func (this *App) Method(method string, path string, handle func(req Request, res Response)) {
//}

//func (this *App) GetNext(path string, handle func(req Request, res Response, next func(e error))) {
//	this.MethodNext("GET", path, handle)
//}
//func (this *App) PostNext(path string, handle func(req Request, res Response, next func(e error))) {
//	this.MethodNext("POST", path, handle)
//}
//func (this *App) DeleteNext(path string, handle func(req Request, res Response, next func(e error))) {
//	this.MethodNext("DELETE", path, handle)
//}
//func (this *App) PutNext(path string, handle func(req Request, res Response, next func(e error))) {
//	this.MethodNext("PUT", path, handle)
//}
//func (this *App) PatchNext(path string, handle func(req Request, res Response, next func(e error))) {
//	this.MethodNext("PATCH", path, handle)
//}
//func (this *App) MethodNext(method string, path string, handle func(req Request, res Response, next func(e error))) {
//	this.handlers = append(this.handlers, Handler{method: strings.ToUpper(method), path: PathToReg(path), handle: handle})
//	log.Printf("methodNext handles len:%d\n", len(this.handlers))
//}

//func (this *App) Append(base string, mount *Route) {
//	this.Mount(base, mount)
//}
