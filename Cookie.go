package rest

import (
	"net/url"
	"time"
)

type Cookie struct {
	Name, Value string     //name=value
	Expires     *time.Time //;Expires=Tue, 15 Jan 2013 21:47:38 GMT
	Domain      string     //;Domain=value
	Path        string     //;Path=value
	Secure      bool       //;Secure
	HttpOnly    bool       //; HttpOnly
	Comment     string     //;Comment=value
}

func (this *Cookie) Encode() string {
	if 0 == len(this.Name) {
		return ""
	}
	var str = url.QueryEscape(this.Name) + "=" + url.QueryEscape(this.Value)
	if nil != this.Expires {
		str += "; Expires=" + this.Expires.UTC().Format(GMT_FORMAT)
	}
	if 0 != len(this.Domain) {
		str += "; Domain=" + this.Domain
	}
	if 0 != len(this.Path) {
		str += "; Path=" + this.Path
	}
	if this.Secure {
		str += "; Secure"
	}
	if this.HttpOnly {
		str += "; HttpOnly"
	}
	return str
}

func (this *Cookie) SetMaxAge(seconds int) {
	t := time.Now().Add(time.Duration(int64(time.Second * time.Duration(seconds))))
	this.Expires = &t
}
