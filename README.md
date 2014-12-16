rest
====

go rest!

## Coming Soon!

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

func (this *UserRouter) Router() rest.App {
	app := &rest.App{}
	app.Get("/", this.Info)
	app.Get("/info", this.Info)
	app.Post("/", this.Post)
	return app
}

func (this *UserRouter) Info(req rest.Request, res rest.Response) {

	res.Json(&User{"jim"})
}

func (this *UserRouter) Post(req rest.Request, res rest.Response) {
	res.Json(req.Fields)
}

func main() {
	app := &rest.App{}
	app.Use(func(req rest.Request, res rest.Response, next func()) {
		defer func() {
			e := recover()
			if nil != e {
				log.Println(e)
				res.Json(e)
			}
		}()
		next()
	})
	app.Use(func(req rest.Request, res rest.Response, next func()) {
		begin := time.Now()
		next(nil)
		cost := (time.Now().UnixNano() - begin.UnixNano()) / 1000000
		log.Printf("%s %s %dms\n", req.Method, req.Path, cost)
	})
	user := UserRouter{}
	app.Mount("/", user.Router())
	app.UsePath("/", middleware.Static("./public"))
	app.Get("/setcookie", func(req rest.Request, res rest.Response) {
		res.Cookie("name", "jim")
		res.Cookie("name1", "jim1")
		res.CookieMaxAge("name1", "jim1", 30)
		res.Send("world")
	})
	app.Get("/api/user/:id", func(req rest.Request, res rest.Response) {
		res.Json(req.Params)
	})
	app.Get("/api/user/:name/first ", func(req rest.Request, res rest.Response) {
		res.Json(req.Params)
	})
	app.Listen(1000)
}
```