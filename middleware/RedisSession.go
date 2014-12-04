package middleware

import "rest"

type RedisSessionConf struct {
	maxAge int
	sid    string
}

func RedisSession(dir string, conf ...RedisSessionConf) func(request *rest.Request, response *rest.Response, next func(e error)) {
	return nil
}
