package lua

import (
	"fmt"

	"github.com/mono83/oscar"
	lua "github.com/yuin/gopher-lua"
)

type assertion struct {
	Actual, Expected interface{}
	Qualifier        string
	Doc              string
}

func (a assertion) Equals(L *lua.LState, ctx *oscar.Context) error {
	if len(a.Qualifier) > 0 {
		ctx.Tracef(`Assert %s "%v" (actual, left) equals "%v"`, a.Qualifier, a.Actual, a.Expected)
	} else {
		ctx.Tracef(`Assert "%v" (actual, left) equals "%v"`, a.Actual, a.Expected)
	}

	if a.Actual != a.Expected {
		ctx.AssertFinished(false, "==", a.Qualifier, a.Doc, a.Actual, a.Expected)
		throwLua(
			L,
			ctx,
			`Assertion failed. "%v" (actual, left) != "%v".%s`,
			a.Actual,
			a.Expected,
			a.Doc,
		)
		return assertEqualsError{assertion: a}
	}

	ctx.AssertFinished(true, "== ", a.Qualifier, a.Doc, a.Actual, a.Expected)
	return nil
}

type assertEqualsError struct {
	assertion
}

func (a assertEqualsError) Error() string {
	return fmt.Sprintf(
		`Assertion failed. "%v" (actual, left) != "%v".%s`,
		a.Actual,
		a.Expected,
		a.Doc,
	)
}
