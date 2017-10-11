package oscar

import (
	"fmt"
	"github.com/mono83/oscar/util/jsonPath"
	"github.com/yuin/gopher-lua"
)

func lJSONPathExtract(L *lua.LState) int {
	tc := luaToTestCase(L)
	path := tc.Interpolate(L.CheckString(2))
	source := tc.Interpolate(L.CheckString(3))

	value, err := jsonPath.Extract([]byte(source), path)
	if err != nil {
		L.RaiseError(err.Error())
		return 0
	}

	L.Push(lua.LString(value))
	return 1
}

func lTestCaseAssertJSONXPath(L *lua.LState) int {
	tc := luaToTestCase(L)
	xpath := tc.Interpolate(L.ToString(2))
	expected := tc.Interpolate(L.ToString(3))
	doc := L.OptString(4, "")

	// Reading response body
	body := tc.Get("http.response.body")

	tc.Trace(`Reading JSON XPath "%s"`, xpath)

	// Extracting json path
	actual, err := jsonPath.Extract([]byte(body), xpath)
	if err != nil {
		tc.assertDone(err)
		L.RaiseError(err.Error())
	} else {
		tc.Trace(`Assert "%s" (actual, left) equals "%s"`, actual, expected)
		success := actual == expected
		if !success {
			err := fmt.Errorf(
				`JSON XPath assertion failed. "%s" (actual, left) != "%s".%s`,
				actual,
				expected,
				doc,
			)
			tc.assertDone(err)
			L.RaiseError("Assertion failed")
		} else {
			tc.assertDone(nil)
			tc.Trace(`Assertion OK. "%s" == "%s"`, xpath, expected)
		}
	}

	return 0
}
