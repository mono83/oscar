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

	_ = assertion{
		Actual:    actual,
		Expected:  expected,
		Qualifier: "",
		Doc:       doc,
	}.Equals(L, tc)

	return 0
}
