package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

//request context,such as request session,request parameter errors
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

//http request object wrapper
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
	this.Queries = parseQueryString(this.Req.URL.RawQuery)
	this.Context = &RequestContext{}
	if strings.Contains(this.ContentType(), "application/x-www-form-urlencoded") {
		body := &bytes.Buffer{}
		this.Context.Body = body
		io.Copy(body, this.Req.Body)
		this.Fields = parseQueryString(body.String())
	}
	if strings.Contains(this.ContentType(), "multipart/form-data") {
		this.parseMultiparts()
	}
}

func parseQueryString(queryString string) map[string][]string {
	var kvStrings = strings.Split(queryString, "&")

	result := make(map[string][]string, len(kvStrings))
	for _, kvString := range kvStrings {
		kvs := strings.Split(kvString, "=")
		if nil == result[kvs[0]] {
			if 2 <= len(kvs) {
				v, _ := url.QueryUnescape(kvs[1])
				result[kvs[0]] = []string{v}
			} else {
				result[kvs[0]] = []string{""}
			}
		} else {
			if 2 <= len(kvs) {
				v, _ := url.QueryUnescape(kvs[1])
				result[kvs[0]] = append(result[kvs[0]], v)
			} else {
				result[kvs[0]] = append(result[kvs[0]], "")
			}
		}
	}
	return result
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

//get param from url as a fieldValidator
//filed value is from url eg. /user/:id
func (this *Request) Param(name string) *FieldValidator {
	value, exists := this.Params[name]
	return &FieldValidator{Validator: Validator{Key: name, Exists: exists, GoOn: true, Req: this}, Value: value}

}

//get param from url's query as a fieldValidator
//filed value is from url eg. /user?id=xx
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

//get param from request body as a fieldValidator
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

//get fileValidator from multipart/form-data
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

//get param from requets,it will test url param,query,body
func (this *Request) GetParam(name string) *FieldValidator {
	ps := this.GetParams(name)
	if 0 == len(ps) {
		return &FieldValidator{Validator: Validator{Key: name, Exists: false, GoOn: true, Req: this}, Value: ""}
	} else {
		return &FieldValidator{Validator: Validator{Key: name, Exists: true, GoOn: true, Req: this}, Value: ps[0]}
	}
}

//get params from requets,it will test url param,query,body
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

//get session object
func (this *Request) Session() ISession {
	return this.Context.Session
}

//get cookie object
func (this *Request) GetCookie(name string) *http.Cookie {
	cookie, e := this.Req.Cookie(name)
	if nil != e {
		return nil
	} else {
		return cookie
	}
}

//get request cookie value
func (this *Request) Cookie(name string) string {
	cookie, e := this.Req.Cookie(name)
	if nil != e {
		return ""
	} else {
		return cookie.Value
	}
}

//check request header if it accept the specified type(t)
func (this *Request) Accept(t string) bool {
	acceptArr := strings.Split(this.Get("Accept"), ",")
	for _, a := range acceptArr {
		if a == t || a == t[(strings.LastIndex(a, "/")+1):] {
			return true
		}
	}
	return false
}

//get remote host ip address
func (this *Request) Ip() string {
	return this.Ips()[0]
}

//get remote host ip addresses , from "X-Forwarded-For" header if TRUST_PROXY enabled
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

//check the request is a ajax request
//only check the "X-Requested-With" if equals "XMLHttpRequest"
func (this *Request) Xhr() bool {
	return "XMLHttpRequest" == this.Get("X-Requested-With")
}

//get request origin url
func (this *Request) OriginUrl() string {
	return this.Req.URL.String()
}

//get request origin url object
func (this *Request) Url() *url.URL {
	return this.Req.URL
}

//get request protocol
func (this *Request) Protocol() string {

	if this.App.Enabled(TRUST_PROXY) {
		proto := this.Get("X-Forwarded-Proto")
		if 0 < len(proto) {
			return proto
		}
	}
	return this.Req.URL.Scheme
}

//check request if secure .
//check the request protocol if equals "https"
func (this *Request) IsSecure() bool {
	return "https" == this.Protocol()
}
func (this *Request) ContentType() string {
	return this.Req.Header.Get("Content-type")
}

//get header from http request
func (this *Request) Get(head string) string {
	return this.Req.Header.Get(head)
}

//if request param validate has errors , Panic method will throw ParamError
func (this *Request) Panic() {
	if 0 < len(this.Context.ParamErrors) {
		panic(&ParamError{this.Context.ParamErrors})
	}
}
