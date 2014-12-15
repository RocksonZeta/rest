package rest

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"rest/utils"
)

type Request struct {
	Req     *http.Request
	App     *App
	Host    string
	Method  string
	Path    string
	Base    string
	Params  map[string]string      //params in url path
	Queries map[string][]string    //the query params
	Fields  map[string][]string    //form field or upload fields
	Files   map[string][]*FormFile //upload files
	Session ISession
}

func (this *Request) Init() {
	this.Path = this.Req.URL.Path
	this.Base = "/"
	this.Method = this.Req.Method
	this.Host = this.Req.Host
	this.Queries = utils.ParseQueryString(this.Req.URL.RawQuery)
	if strings.Contains(this.ContentType(), "application/x-www-form-urlencoded") {
		body := &bytes.Buffer{}
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
		//defer file.Close()
		//path := this.genRandomFile(path.Ext(item.Filename))
		//of, e := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
		//if nil != e {
		//	panic(e.Error())
		//}
		//defer of.Close()
		//writeLen, e := io.Copy(of, file)
		//if nil != e {
		//	panic(e.Error())
		//}
		result[i] = formFile
	}
	return result
}

func (this *Request) genRandomFile(suffix string) string {
	fmt.Println("gen random file", suffix)
	rint := rand.New(rand.NewSource(time.Now().UnixNano())).Int63()
	name := strconv.FormatInt(rint, 16)
	if dir, ok := this.App.Env[UPLOAD_DIR]; ok {
		if dir, ok := dir.(string); ok {
			if fileInfo, e := os.Stat(dir); nil != e || !fileInfo.IsDir() {
				os.MkdirAll(dir, 0755)
			}
			return path.Join(dir, name+suffix)
		}
	}
	return path.Join(os.TempDir(), name+suffix)

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
