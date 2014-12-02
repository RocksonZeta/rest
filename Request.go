package rest

import (
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"

	"rest/utils"
)

type Request struct {
	Req     *http.Request
	App     *App
	Host    string
	Method  string
	Path    string
	Params  map[string]string
	Queries map[string][]string
	Fields  map[string][]string
	Files   map[string][]*FormFile
}

func (this *Request) Init() {
	this.Path = this.Req.URL.Path
	this.Queries = utils.ParseQueryString(this.Req.URL.RawQuery)
	this.Fields = utils.ParseQueryString(this.Req.PostForm.Encode())
	if(this.ContentType())
}

func (this *Request) parseMultiparts() {
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
		defer file.Close()
		path := this.genRandomFile(path.Ext(item.Filename))
		of, e := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
		if nil != e {
			panic(e.Error())
		}
		defer of.Close()
		writeLen, e := io.Copy(of, file)
		if nil != e {
			panic(e.Error())
		}
		formFile.Path = path
		formFile.Size = writeLen
		result[i] = formFile
	}
	return result
}

func (this *Request) genRandomFile(suffix string) string {
	rint := rand.New(rand.NewSource(time.Now().UnixNano())).Int63()
	name := strconv.FormatInt(rint, 16)
	if dir, ok := this.App.Env[UPLOAD_DIR]; ok {
		if dir, ok := dir.(string); ok {
			if fileInfo, e := os.Stat(dir); nil != e || !fileInfo.IsDir() {
				os.MkdirAll(dir, 0755)
			}
			return path.Join(dir, name+"."+suffix)
		}
	}
	return path.Join(os.TempDir(), name+"."+suffix)

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

//headers
func (this *Request) Cookie(name string) string {
	return ""
}
func (this *Request) GetCookie(name string) *Cookie {
	return nil
}
func (this *Request) Cookies() []*Cookie {
	return nil
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
	return ""
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
	return this.Req.Header.Get("Content-type");
}

