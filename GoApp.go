package rest

import (
	"log"
	"net/http"
	"strconv"
)

type GoApp struct {
	Mount
	env map[string]interface{}
	//handlers []Handler
}

func (this *GoApp) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	request := &GoRequest{app: this, req: req, method: req.Method}
	response := &GoResponse{app: this, res: res}
	this.exec(request, response, 0)
}

func (this *GoApp) exec(request Request, response Response, i int) {
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

func (this *GoApp) Listen(port int) {
	log.Printf("server listen at:%d\n", port)

	http.ListenAndServe(":"+strconv.Itoa(port), this)
}

func (this *GoApp) Use(handle func(req Request, res Response, next func(e error))) {
	this.UsePath("", handle)
}

func (this *GoApp) UsePath(path string, handle func(req Request, res Response, next func(e error))) {
	this.MethodNext("", path, handle)
}

func (this *GoApp) GetEnv(name string) string {
	return ""
}
func (this *GoApp) SetEnv(name string, value interface{}) {

}
func (this *GoApp) Enable(name string) {

}
func (this *GoApp) Enabled(name string) bool {
	return false
}
func (this *GoApp) Disable(name string) {

}
func (this *GoApp) Disabled(name string) bool {
	return false
}
func (this *GoApp) Get(path string, handle func(req Request, res Response)) {
	this.MethodNext("GET", path, func(req Request, res Response, next func(e error)) {
		handle(req, res)
	})
}

//func (this *GoApp) Post(path string, handle func(req Request, res Response)) {
//	this.Method("POST", path, handle)
//}
//func (this *GoApp) Delete(path string, handle func(req Request, res Response)) {
//	this.Method("DELETE", path, handle)
//}
//func (this *GoApp) Put(path string, handle func(req Request, res Response)) {
//	this.Method("PUT", path, handle)
//}
//func (this *GoApp) Patch(path string, handle func(req Request, res Response)) {
//	this.Method("PATCH", path, handle)
//}
//func (this *GoApp) Method(method string, path string, handle func(req Request, res Response)) {
//}

//func (this *GoApp) GetNext(path string, handle func(req Request, res Response, next func(e error))) {
//	this.MethodNext("GET", path, handle)
//}
//func (this *GoApp) PostNext(path string, handle func(req Request, res Response, next func(e error))) {
//	this.MethodNext("POST", path, handle)
//}
//func (this *GoApp) DeleteNext(path string, handle func(req Request, res Response, next func(e error))) {
//	this.MethodNext("DELETE", path, handle)
//}
//func (this *GoApp) PutNext(path string, handle func(req Request, res Response, next func(e error))) {
//	this.MethodNext("PUT", path, handle)
//}
//func (this *GoApp) PatchNext(path string, handle func(req Request, res Response, next func(e error))) {
//	this.MethodNext("PATCH", path, handle)
//}
//func (this *GoApp) MethodNext(method string, path string, handle func(req Request, res Response, next func(e error))) {
//	this.handlers = append(this.handlers, Handler{method: strings.ToUpper(method), path: PathToReg(path), handle: handle})
//	log.Printf("methodNext handles len:%d\n", len(this.handlers))
//}

func (this *GoApp) Attach(base string, mount *Mount) {
	this.Append(base, mount)
}
