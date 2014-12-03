package rest

import (
	"net/url"
	"time"
)

type Cookie struct {
	Name, Value string    //name=value
	Path        string    //;Path=value
	Domain      string    //;Domain=value
	Expires     time.Time //;Expires=Tue, 15 Jan 2013 21:47:38 GMT
	Secure      bool      //;Secure
	HttpOnly    bool      //; HttpOnly
	Comment     string    //;Comment=value
}

func (this *Cookie) encode() string {
	if 0 == len(this.Name) {
		return ""
	}
	var str = url.QueryEscape(this.Name) + "=" + url.QueryEscape(this.Value)
	if 0 != len(this.Domain) {
		str += "; Domain=" + this.Domain
	}
	if 0 != len(this.Path) {
		str += "; Path=" + this.Path
	}

	if nil != this.Expires {
		str += "; Expires=" + this.Expires.Format(GMT_FORMAT)
	}
	if this.Secure {
		str += "; Secure"
	}
	if this.HttpOnly {
		str += "; HttpOnly"
	}
}
