package rest

import (
	"io"
	"net/http"
)

type Response struct {
	Resp *http.ResponseWriter
	App  *App
}

func (this *Response) Send(body string) {
	io.WriteString(*this.Resp, body)
}
func (this *Response) SendFile(file string) {

}
func (this *Response) Download(path string) {

}
func (this *Response) Json(obj interface{})                           {}
func (this *Response) Jsonp(obj interface{})                          {}
func (this *Response) Render(tpl string, data map[string]interface{}) {}
func (this *Response) Redirect(url string)                            {}
func (this *Response) Status(status int)                              {}
func (this *Response) Location(location string)                       {}
func (this *Response) Cookie(cookie Cookie)                           {}
func (this *Response) ClearCookie(name string)                        {}
func (this *Response) ContentType(contentType string)                 {}
func (this *Response) Set(name string, value string)                  {}
func (this *Response) Get(name string)                                {}
