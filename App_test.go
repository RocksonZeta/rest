package rest

import (
	. "fmt"
	"log"
	"net/http"
	"strings"
	"testing"
)

type MyHandler struct {
}

func (this MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Print("get", r.RequestURI)

	if 0 == strings.Index(r.RequestURI, "/") {
		http.FileServer(http.Dir("./")).ServeHTTP(w, r)
		return
	} else {
		Fprint(w, "hello")
	}
}

func TestListen(t *testing.T) {
	log.Println("hello")
}
