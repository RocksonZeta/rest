package rest

import (
	"net/http"
	"strconv"
)

//represent ours application,
type App struct {
	Router
	Env map[string]interface{}
}

//new an Application
func NewApp() *App {
	app := &App{Env: map[string]interface{}{}}
	app.Enable(TRUST_PROXY)
	return app
}

//implements the http.Handler
func (this *App) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	request := Request{Req: req, App: this}
	request.Init()
	response := Response{Resp: res, App: this}
	this.Exec(request, response, 0)
}

func (this *App) Listen(port int) {
	http.ListenAndServe(":"+strconv.Itoa(port), this)
}

//get environment variable from the app
func (this *App) GetEnv(name string) interface{} {
	return this.Env[name]
}

//set environment variable
func (this *App) SetEnv(name string, value interface{}) {
	this.Env[name] = value
}

//set environment variable to true
func (this *App) Enable(name string) {
	this.Env[name] = true
}

//peek environment variable weather to be true
func (this *App) Enabled(name string) bool {
	if r, ok := this.Env[name].(bool); ok {
		return r
	} else {
		return false
	}
}

//set environment variable to false
func (this *App) Disable(name string) {
	this.Env[name] = false
}

//peek environment variable weather to be false
func (this *App) Disabled(name string) bool {
	if r, ok := this.Env[name].(bool); ok {
		return r
	} else {
		return true
	}
}
