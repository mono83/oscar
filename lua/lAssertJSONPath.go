package lua

import (
	"github.com/mono83/oscar/util/jsonpath"
	"github.com/yuin/gopher-lua"
)

func lAssertJSONPath(L *lua.LState) int {
	tc := lContext(L)
	xpath := tc.Interpolate(L.ToString(2))
	expected := tc.Interpolate(L.ToString(3))
	doc := L.OptString(4, "")

	// Reading response body
	body := tc.Get("http.response.body")

	// Extracting json path
	actual, err := jsonpath.Extract([]byte(body), xpath)

	if err != nil {
		throwLua(L, tc, "JSON XPATH error on %s - %s", xpath, err.Error())
	} else {
		_ = assertion{
			Actual:    actual,
			Expected:  expected,
			Qualifier: xpath,
			Doc:       doc,
		}.Equals(L, tc)
	}

	return 0
}
