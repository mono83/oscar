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

	tc.Trace(`Assert "%s" (actual, left) equals "%s"`, actual, expected)
	success := actual == expected
	if !success {
		err := fmt.Errorf(
			`assertion failed. "%s" (actual, left) != "%s".%s`,
			actual,
			expected,
			doc,
		)
		tc.assertDone(err)
		L.RaiseError("Assertion failed")
	} else {
		tc.Trace(`Assertion OK. "%s" == "%s"`, actual, expected)
		tc.assertDone(nil)
	}

	return 0
}
