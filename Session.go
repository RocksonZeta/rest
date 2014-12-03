package rest

type ISession interface {
	Get(key string) interface{}
	Set(key string, value interface{})
	Delete(key string)
	Length() int
	Destroy()
	Save()
	Reload()
	Regenerate()
}
