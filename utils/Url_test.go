package utils

import "testing"

func TestParseQueryString(t *testing.T) {
	params := ParseQueryString("k=v&k1&k=v1&h=%E5%93%88%E5%93%88")
	if params["h"][0] != "哈哈" {
		t.Errorf("expected 哈哈,but got " + params["h"][0])
	}
}
