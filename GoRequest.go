package rest

import (
	"net/http"
	"net/url"
)

type GoRequest struct {
	app    App
	req    *http.Request
	method string
}

func (this *GoRequest) Param(name string) string {
	return ""
}
func (this *GoRequest) Query(name string) string {
	return ""

}
func (this *GoRequest) Body(name string) string {
	return ""
}
func (this *GoRequest) File(name string) *FormFile {
	return nil
}
func (this *GoRequest) GetParam(name string) string {
	return ""
}

func (this *GoRequest) Params(name string) []string {
	return nil
}
func (this *GoRequest) Queries(name string) []string {
	return nil
}
func (this *GoRequest) Bodies(name string) []string {
	return nil
}
func (this *GoRequest) Files(name string) []*FormFile {
	return nil
}
func (this *GoRequest) GetParams(name string) []string {
	return nil
}

func (this *GoRequest) AllFiles() []*FormFile {
	return nil
}

//headers
func (this *GoRequest) Get(name string) string {
	return ""
}
func (this *GoRequest) Cookie(name string) string {
	return ""
}
func (this *GoRequest) GetCookie(name string) *Cookie {
	return nil
}
func (this *GoRequest) Cookies() []*Cookie {
	return nil
}
func (this *GoRequest) Host() string {
	return ""
}
func (this *GoRequest) Ip() string {
	return ""
}
func (this *GoRequest) Ips() string {
	return ""
}
func (this *GoRequest) xhr() string {
	return ""
}
func (this *GoRequest) Path() string {
	return this.req.URL.Path
}
func (this *GoRequest) OriginUrl() string {
	return ""
}
func (this *GoRequest) Url() *url.URL {
	return this.req.URL
}
func (this *GoRequest) Protocol() string {
	return ""
}
func (this *GoRequest) IsSecure() string {
	return ""
}

func (this *GoRequest) Method() string {
	return this.method
}
