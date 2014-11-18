package rest

type Restable interface{}
type Getable interface {
	Get(path string, handle func(req Request, res Response))
}
type Postable interface {
	Post(path string, handle func(req Request, res Response))
}
type Putable interface {
	Put(path string, handle func(req Request, res Response))
}
type Patchable interface {
	Patch(path string, handle func(req Request, res Response))
}
type Deleteable interface {
	Delete(path string, handle func(req Request, res Response))
}
type Methodable interface {
	Method(method string, path string, handle func(req Request, res Response))
}

type GetNextable interface {
	Get(path string, handle func(req Request, res Response, next func(e error)))
}
type PostNextable interface {
	Post(path string, handle func(req Request, res Response, next func(e error)))
}
type PutNextable interface {
	Put(path string, handle func(req Request, res Response, next func(e error)))
}
type PatchNextable interface {
	Patch(path string, handle func(req Request, res Response, next func(e error)))
}
type DeleteNextable interface {
	Delete(path string, handle func(req Request, res Response, next func(e error)))
}
type MethodNextable interface {
	Method(method string, path string, handle func(req Request, res Response, next func(e error)))
}
