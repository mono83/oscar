package oscar

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/yuin/gopher-lua"
	"time"
)

var colorLogTime = color.New(color.FgWhite)

//var colorLogName = color.New(color.FgWhite)
var colorLogDebug = color.New(color.FgHiBlack)
var colorLogInfo = color.New(color.FgCyan)
var colorLogTestCase = color.New(color.FgHiGreen)
var colorLogError = color.New(color.FgHiYellow)

func lTestCaseDebug(L *lua.LState) int {
	if t := luaToTestCase(L); t != nil {
		t.print("  "+L.ToString(2), colorLogDebug)
	}
	return 0
}

func lTestCaseInfo(L *lua.LState) int {
	if t := luaToTestCase(L); t != nil {
		t.print("  "+L.ToString(2), colorLogInfo)
	}
	return 0
}

func (t TestCase) print(message string, c *color.Color) {
	fmt.Fprintf(
		t.oscar.output(),
		"%s %s\n",
		colorLogTime.Sprint(time.Now().Format("15:04:05")),
		c.Sprint(t.Interpolate(message)),
	)
}

func (t TestCase) logDebug(message string) {
	t.print("  "+message, colorLogDebug)
}

func (t TestCase) logTestCase(message string) {
	t.print(message, colorLogTestCase)
}

func (t TestCase) logError(message string) {
	t.print(message, colorLogError)
}
