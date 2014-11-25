package rest

import (
	"net/http"
	"net/url"
)

type Request struct {
	app     *App
	req     *http.Request
	method  string
	params  map[string]string
	queries map[string][]string
	fields  map[string][]string
	files   map[string][]FormFile
}

func (this *Request) Param(name string) string {
	return ""
}
func (this *Request) Query(name string) string {
	return ""

}
func (this *Request) Field(name string) string {
	return this.req.PostFormValue(name)
}
func (this *Request) File(name string) *FormFile {
	return nil
}
func (this *Request) GetParam(name string) string {
	return ""
}

func (this *Request) Params(name string) []string {
	return nil
}
func (this *Request) Queries(name string) []string {
	return nil
}
func (this *Request) Fields(name string) []string {
	return nil
}
func (this *Request) Files(name string) []*FormFile {
	return nil
}
func (this *Request) GetParams(name string) []string {
	return nil
}

func (this *Request) AllFiles() []*FormFile {
	return nil
}

//headers
func (this *Request) Get(name string) string {
	return ""
}
func (this *Request) Cookie(name string) string {
	return ""
}
func (this *Request) GetCookie(name string) *Cookie {
	return nil
}
func (this *Request) Cookies() []*Cookie {
	return nil
}
func (this *Request) Host() string {
	return ""
}
func (this *Request) Ip() string {
	return ""
}
func (this *Request) Ips() string {
	return ""
}
func (this *Request) xhr() string {
	return ""
}
func (this *Request) Path() string {
	return this.req.URL.Path
}
func (this *Request) OriginUrl() string {
	return ""
}
func (this *Request) Url() *url.URL {
	return this.req.URL
}
func (this *Request) Protocol() string {
	return ""
}
func (this *Request) IsSecure() string {
	return ""
}

func (this *Request) Method() string {
	return this.method
}
