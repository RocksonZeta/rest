package rest

type Response interface {
	Send(body string)
	SendFile(file string)
	Download(path string)
	Json(obj interface{})
	Jsonp(obj interface{})
	Render(tpl string, data map[string]interface{})
	Redirect(url string)
	Status(status int)
	Location(location string)

	Cookie(cookie Cookie)
	ClearCookie(name string)
	ContentType(contentType string)
	Set(name string, value string)
	Get(name string)
}
