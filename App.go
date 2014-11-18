package rest

type App interface {
	GetEnv(name string) string
	SetEnv(name string, value interface{})
	Enable(name string)
	Enabled(name string) bool
	Disable(name string)
	Disabled(name string) bool
	Use(handle func(req Request, res Response, next func(e error)))
	UsePath(path string, handle func(req Request, res Response, next func(e error)))
	Get(path string, handle func(req Request, res Response))
	Post(path string, handle func(req Request, res Response))
	Delete(path string, handle func(req Request, res Response))
	Put(path string, handle func(req Request, res Response))
	Patch(path string, handle func(req Request, res Response))
	Method(method string, path string, handle func(req Request, res Response))

	GetNext(path string, handle func(req Request, res Response, next func(e error)))
	PostNext(path string, handle func(req Request, res Response, next func(e error)))
	DeleteNext(path string, handle func(req Request, res Response, next func(e error)))
	PutNext(path string, handle func(req Request, res Response, next func(e error)))
	PatchNext(path string, handle func(req Request, res Response, next func(e error)))
	MethodNext(method string, path string, handle func(req Request, res Response, next func(e error)))
	//Append(base string, route *Route)

	Listen(port int)
}
