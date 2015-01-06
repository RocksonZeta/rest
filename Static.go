package rest

import (
	"io"
	"os"
	"path"
	"time"
)

//static middleaware configuration
type StaticConf struct {
	CacheControl string
}

//rest static file server middleware
func Static(dir string, conf ...StaticConf) func(request Request, response Response, next func()) {
	stat, e := os.Stat(dir)
	if nil != e {
		panic(e)
	}
	if !stat.IsDir() {
		panic(&RestError{Reason: dir + " directory not exists!"})
	}
	return func(request Request, response Response, next func()) {
		file := path.Join(dir, request.Path)
		fileInfo, e := os.Stat(file)
		if nil != e || fileInfo.IsDir() {
			next()
			return
		}
		since := request.Get("If-Modified-Since")
		if 0 != len(since) {
			sinceTime, e := time.Parse(GMT_FORMAT, since)
			if nil == e && (sinceTime.Unix()-fileInfo.ModTime().Unix() >= 0) {
				response.Status(304)
				return
			}
		}
		response.Set("Last-Modified", fileInfo.ModTime().UTC().Format(GMT_FORMAT))
		openedFile, e := os.Open(file)
		io.Copy(&response, openedFile)
	}
}
