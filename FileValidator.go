package rest

import (
	"io"
	"mime/multipart"
	"os"
	"path"
	"regexp"
	"strings"
)

type FileValidator struct {
	Validator
	FileItem     *multipart.FileHeader
	FileName     string
	ContentType  string
	Min, Max     int64
	fileLimitTip []string
}

func (this *FileValidator) Limit(min, max int64, tip ...string) *FileValidator {
	this.NotEmpty()
	if !this.GoOn {
		return this
	}
	this.Min = min
	this.Max = max
	this.fileLimitTip = tip
	return this
}
func (this *FileValidator) Empty() *FileValidator {
	if nil == this.FileItem {
		this.GoOn = false
	}
	return this
}
func (this *FileValidator) NotEmpty(tip ...string) *FileValidator {
	if this.GoOn && (nil == this.FileItem) {
		this.FireError(this.Key+" file can not be empty.", tip)
	}
	return this
}
func (this *FileValidator) ContentTypeMatch(reg string, tip ...string) *FileValidator {
	this.NotEmpty()
	if this.GoOn && !regexp.MustCompile(reg).MatchString(this.ContentType) {
		this.FireError("Bad content type.", tip)
	}
	return this
}
func (this *FileValidator) FileNameMatch(reg string, tip ...string) *FileValidator {
	this.NotEmpty()
	if this.GoOn && !regexp.MustCompile(reg).MatchString(this.FileName) {
		this.FireError("Bad content type.", tip)
	}
	return this
}
func (this *FileValidator) FileSuffixIn(suffixes []string, tip ...string) *FileValidator {
	this.NotEmpty()
	if 0 == len(suffixes) {
		return this
	}
	si := strings.LastIndex(this.FileName, ".")
	if this.GoOn && -1 == si {
		this.FireError("Bad file suffix.", tip)
	}
	suffix := this.FileName[(si + 1):]
	for _, s := range suffixes {
		if s == suffix {
			return this
		}
	}
	if this.GoOn {
		this.FireError("Bad file suffix.", tip)
	}
	return this
}

func (this *FileValidator) Save(dirPath string, fileName ...string) *FileValidator {
	this.NotEmpty()
	if !this.GoOn {
		return this
	}
	if fi, e := os.Stat(dirPath); os.IsNotExist(e) || !fi.IsDir() {
		os.MkdirAll(dirPath, 0755)
	}
	if nil == this.FileItem {
		this.FireError(this.Key+" file is empty.", []string{})
	}
	uploadFile, e := this.FileItem.Open()
	if nil != e {
		//	panic(e)
		this.FireError("parse　"+this.Key+" file too large.", []string{})
		return this
	}
	defer uploadFile.Close()
	fname := this.FileName
	if 0 < len(fileName) {
		fname = fileName[0]
	}
	filePath := path.Join(dirPath, fname)
	of, e := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if nil != e {
		panic(e)
	}
	defer of.Close()
	buffer := make([]byte, 8*1024)
	var totalLen int64 = 0
	for {
		if -1 != this.Max && this.GoOn && totalLen > this.Max {
			this.FireError(this.Key+" file too large.", this.fileLimitTip)
			of.Close()
			os.Remove(filePath)
			return this
		}
		rlen, re := uploadFile.Read(buffer)
		if re != nil && io.EOF != re {
			panic(re)
		}
		if io.EOF == re {
			break
		}

		if _, we := of.Write(buffer[0:rlen]); we != nil {
			panic(we)
		}
		totalLen += int64(rlen)
	}
	if this.GoOn && totalLen < this.Min {
		of.Close()
		os.Remove(filePath)
		this.FireError(this.Key+" file too small.", this.fileLimitTip)
		return this
	}
	return this
}
