package middleware

import (
	"rest"
)

func Static(dir string, conf map[string]interface{}) func(request *rest.Reqeust, response *rest.Response, next func(e error)) {
	return func(request *rest.Reqeust, response *rest.Response, next func(e error)) {

		next(nil)
	}
}
