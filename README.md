rest
====

go rest!


## Example
```go
package main

import "rest"

type Users struct {
}

func (this *Users) install() *rest.Route {
	route := &rest.Route{}
	route.Get("/:hello", func(req *rest.Request, res *rest.Response) {
		res.Send("hello")
	})
	return route
}

func main() {
	app := new(rest.App)
	app.Mount("/api", (&Users{}).install())
	app.Get("/", func(req *rest.Request, res *rest.Response) {
		res.Send("haha")
	})
	app.Listen(2000)
}

```