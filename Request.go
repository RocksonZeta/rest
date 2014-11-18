package rest

import "net/url"

type Request interface {
	Param(name string) string
	Query(name string) string
	Body(name string) string
	File(name string) *FormFile
	GetParam(name string) string

	Params(name string) []string
	Queries(name string) []string
	Bodies(name string) []string
	Files(name string) []*FormFile
	GetParams(name string) []string

	AllFiles() []*FormFile

	//headers
	Get(name string) string
	Cookie(name string) string
	GetCookie(name string) *Cookie
	Cookies() []*Cookie
	Host() string
	Ip() string
	Ips() string
	xhr() string
	Path() string
	OriginUrl() string
	Url() *url.URL
	Protocol() string
	IsSecure() string

	Method() string
}
