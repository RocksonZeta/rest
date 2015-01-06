rest
====

go rest!<br/>
To build restful application like expressjs.

[![Build Status](https://travis-ci.org/RocksonZeta/rest.svg)](https://travis-ci.org/RocksonZeta/rest)

## Installation
```
$ go get github.com/RocksonZeta/rest
```

## Documentation
- [API Reference](http://godoc.org/github.com/RocksonZeta/rest)

## Features
- Easy to use.
- Easy to extends.
- Parameter validator with its own.

## Examples

### begin
```go
package main

import "github.com/RocksonZeta/rest"
import "log"

func main() {
	app := rest.NewApp()
	app.Get("/", func(req rest.Request, res rest.Response) {
		res.Json(map[string]string{"hello": "world"})
	})
	log.Println("server start at:6161")
	app.Listen(6161)
}
``` 

### normal 
```go
package main

import (
	"log"
	"runtime"
	"strconv"
	"time"
	"github.com/RocksonZeta/rest"
)

type User struct {
	Name string
}

type UserRouter struct {
}

func (this *UserRouter) Router() *rest.Router {
	app := &rest.Router{}
	//the request path should be /user
	app.Get("/", this.info)
	return app
}

func (this *UserRouter) info(req rest.Request, res rest.Response) {
	res.Json(&User{"jim"})
}

func main() {
	app := rest.NewApp()
	app.Use(func(req rest.Request, res rest.Response, next func()) {
		defer func() {
			e := recover()
			if nil != e {
				buf := make([]byte, 1024)
				runtime.Stack(buf, false)
				log.Println(string(buf))
				res.Json(e)
			}
		}()
		next()
	})
	app.Use(func(req rest.Request, res rest.Response, next func()) {
		begin := time.Now()
		next()
		cost := (time.Now().UnixNano() - begin.UnixNano()) / 1e6
		log.Printf("%s %s %dms\n", req.Method, req.Path, cost)
	})
	//use app as static file server
	app.UsePath("/", rest.Static("./public"))
	//enable app has session ability
	app.Use(rest.LocalSession())

	user := UserRouter{}
	//the base path of user router  is /user
	app.Mount("/user", user.Router())

	app.Get("/", func(req rest.Request, res rest.Response) {
		res.Send("hello");
	})

	app.Post("/upload", func(req rest.Request, res rest.Response) {
		fileType := req.Field("type").Int()
		file1 := req.File("file1").Limit(0, 1024).Save("/upload")
		req.Panic()
		log.Printf("upload %s success,type:%d\n", file1.FileName, fileType)
		res.Json(map[string]interface{}{"file": file1.FileName, "type": fileType})
	})
	app.Post("/signin", func(req rest.Request, res rest.Response) {
		email := req.Field("email").Empty().IsEmail().String()
		name := req.Field("name").Len(1, 20).String()
		password := req.Field("password").NotEmpty().Md5()
		addr := req.Field("addr").Empty().Len(3, 100).String()
		homepage := req.Field("homepage").Optional().IsUrl().String()
		age := req.Field("age").Ge(7).Lt(100).Int()
		req.Panic() //panic if params have errors
		res.Json(map[string]interface{}{"name": name, "age": age, "password": password, "addr": addr, "email": email, "homepage": homepage})
	})

	app.Get("/api/user/:id", func(req rest.Request, res rest.Response) {
		res.Json(req.Params)
	})
	app.Get("/tpl", func(req rest.Request, res rest.Response) {
		data := map[string]interface{}{"name": "jim", "age": 12, "user": User{"tom"}}
		res.Render("view/hello.tpl", data)
	})

	log.Println("server listen at:6161")
	app.Listen(6161)
}

```

## License
[MIT License](https://github.com/RocksonZeta/rest/blob/master/LICENSE)