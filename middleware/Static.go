package middleware

import (
	"io"
	"log"
	"os"
	"path"
	"rest"
	"time"
)

type StaticConf struct {
	autoIndex    bool
	cacheControl string
}

func Static(dir string, conf ...StaticConf) func(request *rest.Request, response *rest.Response, next func(e error)) {
	stat, e := os.Stat(dir)
	if nil != e {
		panic(e)
	}
	if !stat.IsDir() {
		panic(&rest.RestError{Reason: dir + " directory not exists!"})
	}
	return func(request *rest.Request, response *rest.Response, next func(e error)) {
		file := path.Join(dir, request.Path)
		fileInfo, e := os.Stat(file)
		if nil != e || fileInfo.IsDir() {
			next(nil)
			return
		}
		since := request.Get("If-Modified-Since")
		if 0 != len(since) {
			sinceTime, e := time.Parse(rest.GMT_FORMAT, since)
			log.Println(sinceTime.Unix(), fileInfo.ModTime().Unix())
			if nil == e && (sinceTime.Unix()-fileInfo.ModTime().Unix() >= 0) {
				response.Status(304)
				return
			}
		}
		response.Set("Last-Modified", fileInfo.ModTime().UTC().Format(rest.GMT_FORMAT))
		openedFile, e := os.Open(file)
		io.Copy(response, openedFile)
	}
}
