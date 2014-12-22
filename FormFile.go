package rest

import (
	"io"
	"mime/multipart"
	"os"
)

type FormFile struct {
	FileName    string
	ContentType string
	File        multipart.File
	limit       int64
}

func (this *FormFile) Limit(len int64) *FormFile {
	this.limit = len
	return this
}

func (this *FormFile) Save(path string) {
	defer this.File.Close()
	of, e := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if nil != e {
		panic(e.Error())
	}
	defer of.Close()
	buffer := make([]byte, 4*1024)
	var totalLen int64 = 0
	for {
		rlen, re := this.File.Read(buffer)
		if re != nil && io.EOF != re {
			panic(re)
		}
		if 0 == rlen {
			break
		}

		if _, we := of.Write(buffer); we != nil {
			panic(we)
		}
		totalLen += int64(rlen)
	}
	return
}
