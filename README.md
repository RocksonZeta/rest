rest
====

go rest!


## Example
```go
package main

import (
	"log"
	"rest"
	"rest/middleware"
	"time"
)

type User struct {
	Name string
}

type UserRouter struct {
}

func (this *UserRouter) Router() *rest.App {
	app := &rest.App{}
	app.Get("/", this.info)
	app.Get("/info", this.info)
	app.Post("/", this.post)
	return app
}

func (this *UserRouter) info(req *rest.Request, res *rest.Response) {

	res.Json(&User{"jim"})
}

func (this *UserRouter) post(req *rest.Request, res *rest.Response) {
	res.Json(req.Fields)
}

func main() {
	app := &rest.App{}
	user := UserRouter{}
	app.Mount("/", user.Router())
	app.UsePath("/", middleware.Static("./public"))
	app.Use(func(req *rest.Request, res *rest.Response, next func(e error)) {
		begin := time.Now()
		next(nil)
		cost := (time.Now().UnixNano() - begin.UnixNano()) / 1000000
		log.Printf("%s %s %dms\n", req.Method, req.Path, cost)
	})
	app.Get("/hello", func(req *rest.Request, res *rest.Response) {
		res.Cookie("name", "jim")
		res.Cookie("name1", "jim1")
		res.CookieMaxAge("name1", "jim1", 30)
		res.Send("world")
	})
	app.Get("/api/user/:id", func(req *rest.Request, res *rest.Response) {
		res.Json(req.Params)
	})
	app.Get("/api/user/:name/first ", func(req *rest.Request, res *rest.Response) {
		res.Json(req.Params)
	})
	app.Listen(1000)
}
```