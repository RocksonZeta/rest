package rest

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"rest/utils"
)

type RequestContext struct {
	Session ISession
	Body    *bytes.Buffer
}

type Request struct {
	Req         *http.Request
	App         *App
	Host        string
	Method      string
	Path        string
	Base        string
	Params      map[string]string      //params in url path
	Queries     map[string][]string    //the query params
	Fields      map[string][]string    //form field or upload fields
	Files       map[string][]*FormFile //upload files
	Context     *RequestContext
	ParamErrors []string
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
	this.Req.ParseMultipartForm(1024 * 1024 * 8)
	this.Fields = this.Req.MultipartForm.Value
	this.Files = make(map[string][]*FormFile)
	for k, v := range this.Req.MultipartForm.File {
		this.Files[k] = this.parseMultipartFile(v)
	}
}

func (this *Request) parseMultipartFile(fileHeaders []*multipart.FileHeader) []*FormFile {
	result := make([]*FormFile, len(fileHeaders))
	for i, item := range fileHeaders {
		formFile := &FormFile{FileName: item.Filename, ContentType: item.Header.Get("Content-type")}
		file, e := item.Open()
		if nil != e {
			panic(e.Error())
		}
		formFile.File = file

		result[i] = formFile
	}
	return result
}

func (this *Request) Query(name string) string {
	if 0 >= len(this.Queries) {
		return ""
	} else {
		return this.Queries[name][0]
	}
}
func (this *Request) Field(name string) string {
	if 0 >= len(this.Fields) {
		return ""
	} else {
		return this.Fields[name][0]
	}
}
func (this *Request) File(name string) *FormFile {
	if 0 >= len(this.Files) {
		return nil
	} else {
		return this.Files[name][0]
	}
}
func (this *Request) GetParam(name string) string {
	return ""
}
func (this *Request) GetParams(name string) string {
	return ""
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

func (this *Request) Ip() string {

	ip := ""
	if this.App.Enabled(TRUST_PROXY) {
		return ""
	}
	ip = this.Req.RemoteAddr
	return ip
}
func (this *Request) Ips() string {
	return ""
}
func (this *Request) Xhr() string {
	return ""
}

func (this *Request) OriginUrl() string {
	return this.Req.URL.String()
}
func (this *Request) Url() *url.URL {
	return this.Req.URL
}
func (this *Request) Protocol() string {
	return ""
}
func (this *Request) IsSecure() string {
	return ""
}
func (this *Request) ContentType() string {
	return this.Req.Header.Get("Content-type")
}
func (this *Request) Get(head string) string {
	return this.Req.Header.Get(head)
}
