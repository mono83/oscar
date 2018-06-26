package lua

import (
	"github.com/yuin/gopher-lua"
)

// lAssertEquals performs equality assertions
func lAssertEquals(L *lua.LState) int {
	tc := lContext(L)
	actual := tc.Interpolate(L.ToString(2))
	expected := tc.Interpolate(L.ToString(3))
	doc := L.OptString(4, "")

	tc.Tracef(`Assert "%s" (actual, left) equals "%s"`, actual, expected)
	success := actual == expected
	if !success {
		L.RaiseError(
			`Assertion failed. "%s" (actual, left) != "%s".%s`,
			actual,
			expected,
			doc,
		)
	} else {
		tc.AssertFinished(nil)
	}

	return 0
}
