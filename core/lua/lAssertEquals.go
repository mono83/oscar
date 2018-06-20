package lua

import (
	"github.com/yuin/gopher-lua"
	"fmt"
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
		err := fmt.Errorf(
			`assertion failed. "%s" (actual, left) != "%s".%s`,
			actual,
			expected,
			doc,
		)
		tc.AssertFinished(err)
		L.RaiseError("Assertion failed")
	} else {
		tc.AssertFinished(nil)
	}

	return 0
}
