package middleware

import "rest"

type RedisSessionConf struct {
	MaxAge int
	Sid    string
}

func RedisSession(dir string, conf ...RedisSessionConf) func(request rest.Request, response rest.Response, next func()) {
	return nil
}
