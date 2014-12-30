package rest

import (
	"math/rand"
	"net/http"
)

type LocalSessionConf struct {
	MaxAge     int
	SessionKey string
}

func initConf(conf *LocalSessionConf) {
	if 0 >= conf.MaxAge {
		conf.MaxAge = 3600
	}
	if 0 == len(conf.SessionKey) {
		conf.SessionKey = "sid"
	}
}

func generateSessionId(length int) string {
	bs := []byte("0123456789abcedfghijklmnopqrstuvwxyzABCEDFGHIJKLMNOPQRSTUVWXYZ")
	cs := make([]byte, length)
	for i := 0; i < length; i++ {
		cs[i] = bs[rand.Intn(len(bs))]
	}
	return string(cs[:])
}

func LocalSession(confs ...LocalSessionConf) func(request Request, response Response, next func()) {
	var conf LocalSessionConf
	if 0 < len(confs) {
		conf = confs[0]
	}
	initConf(&conf)
	sessions := map[string]*LocalSessionStore{}

	return func(request Request, response Response, next func()) {
		key := request.Cookie(conf.SessionKey)
		session := sessions[key]
		if nil != session {
			response.SetCookie(&http.Cookie{Name: conf.SessionKey, Value: key, MaxAge: conf.MaxAge, HttpOnly: true})
		} else {
			newKey := generateSessionId(32)
			session = &LocalSessionStore{SessionId: newKey, Store: map[string]interface{}{}}
			sessions[newKey] = session
			response.SetCookie(&http.Cookie{Name: conf.SessionKey, Value: newKey, MaxAge: conf.MaxAge, HttpOnly: true})
		}
		request.Context.Session = session
		next()
	}
}

/**
Get(key string) interface{}
Set(key string, value interface{})
Delete(key string)
Length() int
Destroy()
Save()
Reload()
Regenerate()
Has(key string) bool
GetInt(key string) int
GetString(key string) string
*/
type LocalSessionStore struct {
	SessionId string
	Store     map[string]interface{}
}

func (this *LocalSessionStore) Has(key string) bool {
	_, ok := this.Store[key]
	return ok
}
func (this *LocalSessionStore) Get(key string) interface{} {
	return this.Store[key]
}
func (this *LocalSessionStore) GetInt(key string) int {
	if v, ok := this.Store[key].(int); ok {
		return v
	}
	return 0
}
func (this *LocalSessionStore) GetString(key string) string {
	if v, ok := this.Store[key].(string); ok {
		return v
	}
	return ""
}
func (this *LocalSessionStore) Set(key string, value interface{}) {
	this.Store[key] = value
}
func (this *LocalSessionStore) Delete(key string) {
	delete(this.Store, key)
}
func (this *LocalSessionStore) Length() int {
	return len(this.Store)
}
func (this *LocalSessionStore) Destroy() {
	this.Store = map[string]interface{}{}
}
func (this *LocalSessionStore) Save() {

}
func (this *LocalSessionStore) Reload() {

}
func (this *LocalSessionStore) Regenerate() {

}
