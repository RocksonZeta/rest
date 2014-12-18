package rest

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
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

}
func (this *Response) Download(path string) {

}
func (this *Response) Json(obj interface{}) {
	this.ContentType("application/json")
	r, e := json.Marshal(obj)
	if nil != e {
		panic(e)
	} else {
		_, e := this.Resp.Write(r)
		if nil != e {
			log.Panic(e)
		}
	}
}
func (this *Response) Jsonp(obj interface{}) {}
func (this *Response) Render(tpl string, data map[string]interface{}) {
	temp, e := template.ParseFiles(tpl)
	if nil != e {
		panic(e)
	}
	temp.Execute(this, data)
}
func (this *Response) Redirect(url string) {}
func (this *Response) Status(status int) {
	this.Resp.WriteHeader(status)
}
func (this *Response) Location(location string) {}
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
