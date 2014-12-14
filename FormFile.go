package rest

import (
	"bufio"
	"io"
	"mime/multipart"
	"os"
)

type FormFile struct {
	FileName    string
	ContentType string
	File        multipart.File
}

func (this *FormFile) save(file string, limit string) int64 {
	defer this.File.Close()
	of, e := bufio.NewWriterSize(os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644), 8*1024)
	if nil != e {
		panic(e.Error())
	}
	defer of.Close()
	writeLen, e := io.Copy(of, file)
	if nil != e {
		panic(e.Error())
	}
	return writeLen
}
