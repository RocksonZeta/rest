package rest

type Cookie struct {
	name, value string //name=value
	path        string //;Path=value
	domain      string //;Domain=value
	maxAge      int    //;Max-Age=value
	secure      bool   //;Secure
	isHttpOnly  bool   //
	comment     string //;Comment=value
}
