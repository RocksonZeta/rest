package rest

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"
)

type Response struct {
	Resp http.ResponseWriter
	App  *App
}

func (this *Response) Send(body string) (int, error) {
	return io.WriteString(this.Resp, body)
}
func (this *Response) Write(body []byte) (int, error) {
	return this.Resp.Write(body)
}
func (this *Response) SendFile(file string) {
	of, e := os.Open(file)
	if nil != e {
		panic(e)
	}
	defer of.Close()
	io.Copy(this, of)
}
func (this *Response) Download(file string) {
	name := path.Base(file)
	this.Set("Content-Disposition", "attachment; filename=\""+name+"\"")
	this.SendFile(file)
}
func (this *Response) Json(obj interface{}) {
	this.ContentType("application/json")
	r, e := json.Marshal(obj)
	if nil != e {
		panic(e)
	}
	this.Write(r)
}
func (this *Response) Jsonp(obj interface{}, callbackName ...string) {
	this.ContentType("text/javascript")
	r, e := json.Marshal(obj)
	if nil != e {
		panic(e)
	} else {
		cb := "callback"
		if len(callbackName) > 0 {
			cb = callbackName[0]
		}
		cb += "(" + string(r) + ")"
		this.Send(cb)
	}
}
func (this *Response) Render(tpl string, data map[string]interface{}) {
	temp, e := template.ParseFiles(tpl)
	if nil != e {
		panic(e)
	}
	temp.Execute(this, data)
}
func (this *Response) Redirect(url string) {
	this.Set("Location", url)
	this.Status(302)
}
func (this *Response) Status(status int) {
	this.Resp.WriteHeader(status)
}
func (this *Response) SetCookie(cookie *http.Cookie) {
	http.SetCookie(this.Resp, cookie)
}
func (this *Response) Cookie(name, value string) {
	this.SetCookie(&http.Cookie{Name: name, Value: value})
}
func (this *Response) CookieMaxAge(name, value string, maxAge int) {
	cookie := &http.Cookie{Name: name, Value: value, MaxAge: maxAge}
	this.SetCookie(cookie)
}
func (this *Response) ClearCookie(name string) {
	this.CookieMaxAge(name, "", -1)
}
func (this *Response) ContentType(contentType string) {
	this.Set("Content-Type", contentType)
}
func (this *Response) Set(name string, value string) {
	this.Resp.Header().Set(name, value)
}
func (this *Response) Get(name string) string {
	return this.Resp.Header().Get(name)
}
