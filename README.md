rest
====

go rest!<br/>
To build restful application like expressjs.

## Installation
```
$ go get github.com/RocksonZeta/rest
```

## Features
- Easy to use.
- Easy to extends.
- Parameter validator with its own.

## Examples

### begin
```go
package main

import "github.com/RocksonZeta/rest"

func main() {
	app := rest.NewApp()
	app.Get("/", func(req rest.Request, res rest.Response) {
		res.Json(map[string]string{"hello": "world"})
	})
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
	//the request path should be /user/info
	app.Get("/info", this.info)
	app.Post("/", this.post)
	return app
}

func (this *UserRouter) info(req rest.Request, res rest.Response) {
	res.Json(&User{"jim"})
}

func (this *UserRouter) post(req rest.Request, res rest.Response) {
	res.Json(req.Fields)
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
		if req.Session().Has("count") {
			req.Session().Set("count", req.Session().GetInt("count")+1)
		} else {
			req.Session().Set("count", 1)
		}
		res.Send("session:" + strconv.Itoa(req.Session().GetInt("count")))
	})
	app.UsePath("/api", func(req rest.Request, res rest.Response, next func()) {
		log.Println("we get /api")
		next()
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
	app.Get("/api/user/:name/first", func(req rest.Request, res rest.Response) {
		res.Json(req.Params)
	})
	app.Get("/cookie", func(req rest.Request, res rest.Response) {
		res.Cookie("name", "jim")
		res.Cookie("name1", "jim1")
		res.CookieMaxAge("name1", "jim1", 30)
		res.Send("world")
	})
	app.Get("/tpl", func(req rest.Request, res rest.Response) {
		data := map[string]interface{}{"name": "jim", "age": 12, "user": User{"tom"}}
		res.Render("view/hello.tpl", data)
	})
	app.Get("/download", func(req rest.Request, res rest.Response) {
		res.Download("/hello.txt")
	})
	app.Get("/Sendfile", func(req rest.Request, res rest.Response) {
		res.SendFile("public/user.html")
	})
	app.Get("/jsonp", func(req rest.Request, res rest.Response) {
		res.Jsonp(User{"jim"})
	})
	app.Get("/redirect", func(req rest.Request, res rest.Response) {
		log.Println("redirect to ")
		res.Redirect("http://www.baidu.com")
	})
	app.Post("/form", func(req rest.Request, res rest.Response) {
		log.Println("post form ", req.Fields)
		res.Json(req.Fields)
	})

	log.Println("server listen at:6161")
	app.Listen(6161)
}

```

## License
[MIT License](https://github.com/RocksonZeta/rest/blob/master/LICENSE)