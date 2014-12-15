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

func (this *FormFile) save(path string, limit string) int64 {
	defer this.File.Close()
	of, e := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if nil != e {
		panic(e.Error())
	}
	defer of.Close()
	bof := bufio.NewWriterSize(of, 8*1024)
	writeLen, e := io.Copy(bof, this.File)
	if nil != e {
		panic(e.Error())
	}
	return writeLen
}
