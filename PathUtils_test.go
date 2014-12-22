package rest

import (
	"fmt"
	"testing"
)

func TestNamedMatches(t *testing.T) {

	base, params := NamedMatches(PathToReg("^/"), "/api/hello")
	fmt.Println("result", base, params)
}
