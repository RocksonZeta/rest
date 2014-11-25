package rest

import (
	"net/http"
	"net/url"
)

type Request struct {
	*http.Request
	app     *App
	Path    string
	params  map[string]string
	queries map[string][]string
	fields  map[string][]string
	files   map[string][]*FormFile
}

func (this *Request) Init() {
	this.path = this.URL.Path
	var querySTring = this.URL.RawQuery
}
func parse


func (this *Request) Param(name string) string {
	if 0 >= len(this.params) {
		return ""
	} else {
		return this.params[name][0]
	}
}
func (this *Request) Query(name string) string {
	if 0 >= len(this.queries) {
		return ""
	} else {
		return this.queries[name][0]
	}
}
func (this *Request) Field(name string) string {
	if 0 >= len(this.fields) {
		return ""
	} else {
		return this.fields[name][0]
	}
}
func (this *Request) File(name string) *FormFile {
	if 0 >= len(this.files) {
		return ""
	} else {
		return this.files[name][0]
	}
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
//func (this *Request) Cookie(name string) string {
//	return ""
//}
//func (this *Request) GetCookie(name string) *Cookie {
//	return nil
//}
//func (this *Request) Cookies() []*Cookie {
//	return nil
//}
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
