package rest

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"
)

//http response object wrapper
type Response struct {
	Resp http.ResponseWriter
	App  *App
}

//reponse body to client
func (this *Response) Send(body string) (int, error) {
	return io.WriteString(this.Resp, body)
}

//send bytes to client
func (this *Response) Write(body []byte) (int, error) {
	return this.Resp.Write(body)
}

//send file to client
func (this *Response) SendFile(file string) {
	of, e := os.Open(file)
	if nil != e {
		panic(e)
	}
	defer of.Close()
	io.Copy(this, of)
}

//download file as attachment
func (this *Response) Download(file string) {
	name := path.Base(file)
	this.Set("Content-Disposition", "attachment; filename=\""+name+"\"")
	this.SendFile(file)
}

//send object as json to client
func (this *Response) Json(obj interface{}) {
	this.ContentType("application/json")
	r, e := json.Marshal(obj)
	if nil != e {
		panic(e)
	}
	this.Write(r)
}

//reponse jsonp,default callback function name is "callback"
//we can change callback function name using callbackName
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

//render template with data to client
//tpl - the path of the template file
func (this *Response) Render(tpl string, data map[string]interface{}) {
	temp, e := template.ParseFiles(tpl)
	if nil != e {
		panic(e)
	}
	temp.Execute(this, data)
}

//redirect to url
func (this *Response) Redirect(url string) {
	this.Set("Location", url)
	this.Status(302)
}

//set response status
func (this *Response) Status(status int) {
	this.Resp.WriteHeader(status)
}

//set response cookie
func (this *Response) SetCookie(cookie *http.Cookie) {
	http.SetCookie(this.Resp, cookie)
}

//set reponse cookie with name and value
func (this *Response) Cookie(name, value string) {
	this.SetCookie(&http.Cookie{Name: name, Value: value})
}

//set reponse cookie with name, value and maxAge
func (this *Response) CookieMaxAge(name, value string, maxAge int) {
	cookie := &http.Cookie{Name: name, Value: value, MaxAge: maxAge}
	this.SetCookie(cookie)
}

//clear specified cookie
func (this *Response) ClearCookie(name string) {
	this.CookieMaxAge(name, "", -1)
}

//set response content type
func (this *Response) ContentType(contentType string) {
	this.Set("Content-Type", contentType)
}

//set respnose header
func (this *Response) Set(name string, value string) {
	this.Resp.Header().Set(name, value)
}

//get response header
func (this *Response) Get(name string) string {
	return this.Resp.Header().Get(name)
}
