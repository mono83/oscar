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

	tc.oscar.tracef(`Reading JSON XPath "%s"`, xpath)

	// Extracting json path
	actual, err := jsonPath.Extract([]byte(body), xpath)
	if err != nil {
		tc.logError(fmt.Sprintf(
			"Unable to parse JSON XPath %s - %s",
			xpath,
			err.Error(),
		))
		tc.assertDone(false)
		L.RaiseError(err.Error())
	} else {
		tc.oscar.tracef(`Assert "%s" (actual, left) equals "%s"`, actual, expected)
		success := actual == expected
		tc.assertDone(success)
		if !success {
			tc.logError(fmt.Sprintf(
				`JSON XPath assertion failed. "%s" (actual, left) != "%s".%s`,
				actual,
				expected,
				doc,
			))
			L.RaiseError("Assertion failed")
		} else {
			tc.oscar.tracef(`Assertion OK. "%s" == "%s"`, xpath, expected)
		}
	}

	return 0
}
