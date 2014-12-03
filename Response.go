package rest

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Response struct {
	Resp       http.ResponseWriter
	App        *App
	SetCookies map[string]Cookie
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
func (this *Response) Jsonp(obj interface{})                          {}
func (this *Response) Render(tpl string, data map[string]interface{}) {}
func (this *Response) Redirect(url string)                            {}
func (this *Response) Status(status int)                              {}
func (this *Response) Location(location string)                       {}
func (this *Response) SetCookie(cookie Cookie) {
	if nil == this.SetCookies {
		this.SetCookies = map[string]Cookie{cookie.Name: cookie}
	} else {
		this.SetCookies[cookie.Name] = cookie
	}
	s := ""
	for _, item := range this.SetCookies {
		s += item.Encode() + "; "
	}
	if 0 < len(s) {
		s = s[0 : len(s)-2]
		this.Resp.Header().Set("Set-Cookie", s)
	} else {
		this.Resp.Header().Del("Set-Cookie")
	}
}
func (this *Response) Cookie(name, value string) {
	this.SetCookie(Cookie{Name: name, Value: value})
}
func (this *Response) CookieMaxAge(name, value string, maxAge int) {
	cookie := Cookie{Name: name, Value: value}
	cookie.SetMaxAge(maxAge)
	this.SetCookie(cookie)
}
func (this *Response) ClearCookie(name string) {
	this.CookieMaxAge(name, "", -100000)
}
func (this *Response) ContentType(contentType string) {}
func (this *Response) Set(name string, value string) {
	this.Resp.Header().Set(name, value)
}
func (this *Response) Get(name string) string {
	return this.Resp.Header().Get(name)
}
