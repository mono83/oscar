package oscar

import (
	"fmt"
	"github.com/mono83/oscar/util/jsonPath"
	"github.com/yuin/gopher-lua"
)

func lTestCaseAssert(L *lua.LState) int {
	tc := luaToTestCase(L)
	actual := tc.Interpolate(L.ToString(2))
	expected := tc.Interpolate(L.ToString(3))
	doc := L.OptString(4, "")

	tc.oscar.tracef(`Assert "%s" (actual, left) equals "%s"`, actual, expected)
	success := actual == expected
	tc.assertDone(success)
	if !success {
		tc.logError(fmt.Sprintf(
			`Assertion failed. "%s" (actual, left) != "%s".%s`,
			actual,
			expected,
			doc,
		))
		L.RaiseError("Assertion failed")
	} else {
		tc.oscar.tracef(`Assertion OK. "%s" == "%s"`, actual, expected)
	}

	return 0
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
