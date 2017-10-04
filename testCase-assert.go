package oscar

import (
	"fmt"
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
