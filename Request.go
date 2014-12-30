package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"rest/utils"
)

type RequestContext struct {
	Session     ISession
	Body        *bytes.Buffer
	ParamErrors []string
}

type ParamError struct {
	Errors []string
}

func (this *ParamError) Error() string {
	r, _ := json.Marshal(this.Errors)
	return string(r)
}

type Request struct {
	Req     *http.Request
	App     *App
	Host    string
	Method  string
	Path    string
	Base    string
	Params  map[string]string           //params in url path
	Queries map[string][]string         //the query params
	Fields  map[string][]string         //form field or upload fields
	Files   map[string][]*FileValidator //upload files
	Context *RequestContext
}

func (this *Request) Init() {
	this.Path = this.Req.URL.Path
	this.Base = "/"
	this.Method = this.Req.Method
	this.Host = this.Req.Host
	this.Queries = utils.ParseQueryString(this.Req.URL.RawQuery)
	this.Context = &RequestContext{}
	if strings.Contains(this.ContentType(), "application/x-www-form-urlencoded") {
		body := &bytes.Buffer{}
		this.Context.Body = body
		io.Copy(body, this.Req.Body)
		this.Fields = utils.ParseQueryString(body.String())
	}
	if strings.Contains(this.ContentType(), "multipart/form-data") {
		this.parseMultiparts()
	}
}

func (this *Request) parseMultiparts() {
	this.Req.ParseMultipartForm(1024 * 1024 * 32)
	this.Fields = this.Req.MultipartForm.Value
	this.Files = make(map[string][]*FileValidator)
	for k, v := range this.Req.MultipartForm.File {
		fileObjects := make([]*FileValidator, len(v))
		for i, item := range v {
			fileValidator := &FileValidator{Validator: Validator{Key: k, Exists: true, GoOn: true, Req: this}, ContentType: item.Header.Get("Content-type"), FileName: item.Filename, FileItem: item, Max: -1}
			fileObjects[i] = fileValidator
		}
		this.Files[k] = fileObjects
	}
}

func (this *Request) Param(name string) *FieldValidator {
	value, exists := this.Params[name]
	return &FieldValidator{Validator: Validator{Key: name, Exists: exists, GoOn: true, Req: this}, Value: value}

}

func (this *Request) Query(name string) *FieldValidator {
	var value string
	exists := false
	if nil != this.Queries {
		values, ex := this.Queries[name]
		if 0 < len(values) {
			value = values[0]
		}
		exists = ex
	}
	return &FieldValidator{Validator: Validator{Key: name, Exists: exists, GoOn: true, Req: this}, Value: value}
}
func (this *Request) Field(name string) *FieldValidator {
	var value string
	exists := false
	if nil != this.Fields {
		values, ex := this.Fields[name]
		if 0 < len(values) {
			value = values[0]
		}
		exists = ex
	}
	return &FieldValidator{Validator: Validator{Key: name, Exists: exists, GoOn: true, Req: this}, Value: value}
}
func (this *Request) File(name string) *FileValidator {
	if 0 >= len(this.Files) {
		return &FileValidator{Validator: Validator{Key: name, Exists: true, GoOn: true, Req: this}}
	} else {
		if files, ok := this.Files[name]; ok {
			if 0 == len(files) {
				return &FileValidator{Validator: Validator{Key: name, Exists: true, GoOn: true, Req: this}}
			}
			return files[0]
		} else {
			return &FileValidator{Validator: Validator{Key: name, Exists: true, GoOn: true, Req: this}}
		}
	}
}
func (this *Request) GetParam(name string) string {
	ps := this.GetParams(name)
	if 0 == len(ps) {
		return ""
	} else {
		return ps[0]
	}
}
func (this *Request) GetParams(name string) []string {
	param := this.Params[name]
	if 0 < len(param) {
		return []string{param}
	}
	if ps, ok := this.Queries[name]; ok && 0 > len(ps) {
		return ps
	}
	if ps, ok := this.Fields[name]; ok && 0 > len(ps) {
		return ps
	}
	return []string{}
}

func (this *Request) Session() ISession {
	return this.Context.Session
}

func (this *Request) GetCookie(name string) *http.Cookie {
	cookie, e := this.Req.Cookie(name)
	if nil != e {
		return nil
	} else {
		return cookie
	}
}

func (this *Request) Cookie(name string) string {
	cookie, e := this.Req.Cookie(name)
	if nil != e {
		return ""
	} else {
		return cookie.Value
	}
}

func (this *Request) Accept(t string) bool {
	acceptArr := strings.Split(this.Get("Accept"), ",")
	for _, a := range acceptArr {
		if a == t || a == t[(strings.LastIndex(a, "/")+1):] {
			return true
		}
	}
	return false
}

func (this *Request) Ip() string {
	return this.Ips()[0]
}
func (this *Request) Ips() []string {
	if this.App.Enabled(TRUST_PROXY) {
		forwords := this.Get("X-Forwarded-For")
		if 0 == len(forwords) {
			return []string{this.Req.RemoteAddr}
		} else {
			return strings.Split(forwords, ",")
		}
	} else {
		return []string{this.Req.RemoteAddr}
	}
}
func (this *Request) Xhr() bool {
	return "XMLHttpRequest" == this.Get("X-Requested-With")
}

func (this *Request) OriginUrl() string {
	return this.Req.URL.String()
}
func (this *Request) Url() *url.URL {
	return this.Req.URL
}
func (this *Request) Protocol() string {

	if this.App.Enabled(TRUST_PROXY) {
		proto := this.Get("X-Forwarded-Proto")
		if 0 < len(proto) {
			return proto
		}
	}
	return this.Req.URL.Scheme
}
func (this *Request) IsSecure() bool {
	return "https" == this.Protocol()
}
func (this *Request) ContentType() string {
	return this.Req.Header.Get("Content-type")
}

//get head from http request
func (this *Request) Get(head string) string {
	return this.Req.Header.Get(head)
}

func (this *Request) Panic() {
	if 0 < len(this.Context.ParamErrors) {
		panic(&ParamError{this.Context.ParamErrors})
	}
}
