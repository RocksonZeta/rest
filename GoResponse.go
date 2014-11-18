package rest

import (
	"io"
	"net/http"
)

type GoResponse struct {
	app App
	res http.ResponseWriter
}

func (this *GoResponse) Send(body string) {
	io.WriteString(this.res, body)
}
func (this *GoResponse) SendFile(file string) {

}
func (this *GoResponse) Download(path string) {

}
func (this *GoResponse) Json(obj interface{})                           {}
func (this *GoResponse) Jsonp(obj interface{})                          {}
func (this *GoResponse) Render(tpl string, data map[string]interface{}) {}
func (this *GoResponse) Redirect(url string)                            {}
func (this *GoResponse) Status(status int)                              {}
func (this *GoResponse) Location(location string)                       {}
func (this *GoResponse) Cookie(cookie Cookie)                           {}
func (this *GoResponse) ClearCookie(name string)                        {}
func (this *GoResponse) ContentType(contentType string)                 {}
func (this *GoResponse) Set(name string, value string)                  {}
func (this *GoResponse) Get(name string)                                {}
