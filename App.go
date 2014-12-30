package rest

import (
	"net/http"
	"strconv"
)

type App struct {
	Router
	Env map[string]interface{}
}

func NewApp() *App {
	app := &App{Env: map[string]interface{}{}}
	app.Enable(TRUST_PROXY)
	return app
}

func (this *App) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	request := Request{Req: req, App: this}
	request.Init()
	response := Response{Resp: res, App: this}
	this.Exec(request, response, 0)
}

func (this *App) Listen(port int) {
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
