package middleware

import (
	"io"
	"log"
	"os"
	"path"
	"rest"
)

func Static(dir string, conf map[string]interface{}) func(request *rest.Request, response *rest.Response, next func(e error)) {
	stat, e := os.Stat(dir)
	if nil != e {
		panic(e)
	}
	if !stat.IsDir() {
		panic(&rest.RestError{Reason: dir + " directory not exists!"})
	}
	return func(request *rest.Request, response *rest.Response, next func(e error)) {
		file := path.Join(dir, request.Path)
		log.Printf("static file:%s\n", file)
		fileInfo, e := os.Stat(file)
		if nil != e || fileInfo.IsDir() {
			next(nil)
			return
		}
		openedFile, e := os.OpenFile(file, os.O_RDONLY, 0666)
		io.Copy(response, openedFile)
	}
}
