package rest

//http request session interface
type ISession interface {
	//get value from session
	Get(key string) interface{}
	//set value to session
	Set(key string, value interface{})
	//delete specified sessoin key
	Delete(key string)
	//return the lenght of the keys in session
	Length() int
	//clear all keys in session
	Destroy()
	//save new key and value to session store
	Save()
	//reload session key and value to current session from session store
	Reload()
	//regenerate session,old session will be destroyed
	Regenerate()
	//check the session if has key
	Has(key string) bool
	//get session value as integer
	GetInt(key string) int
	//get session value as string
	GetString(key string) string
}
