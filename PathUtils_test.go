package rest

import "testing"

func TestNamedMatches(t *testing.T) {
	path := "/user/123"
	base, params := NamedMatches(PathToReg("/user/:id"), path)
	if base != path || params["id"] != "123" {
		t.Errorf("/user/123 not matched to /uesr/:id")
	}
}
