package oscar

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/yuin/gopher-lua"
	"time"
)

var colorLogTime = color.New(color.FgBlue)
var colorLogName = color.New(color.FgWhite)
var colorLogMessage = color.New(color.FgHiBlack)
var colorLogInfo = color.New(color.FgHiCyan)
var colorLogError = color.New(color.FgHiYellow)

func lTestCaseLog(L *lua.LState) int {
	if t := luaToTestCase(L); t != nil {
		t.log(L.ToString(2))
	}
	return 0
}

func (t TestCase) log(message string) {
	fmt.Fprintf(
		t.oscar.output(),
		"%s %s %s\n",
		colorLogTime.Sprint(time.Now().Format("15:04:05")),
		colorLogName.Sprint("["+t.Name+"]"),
		colorLogMessage.Sprint(t.Interpolate(message)),
	)
}

func (t TestCase) logInfo(message string) {
	fmt.Fprintf(
		t.oscar.output(),
		"%s %s %s\n",
		colorLogTime.Sprint(time.Now().Format("15:04:05")),
		colorLogName.Sprint("["+t.Name+"]"),
		colorLogInfo.Sprint(message),
	)
}

func (t TestCase) logError(message string) {
	fmt.Fprintf(
		t.oscar.output(),
		"%s %s %s\n",
		colorLogTime.Sprint(time.Now().Format("15:04:05")),
		colorLogName.Sprint("["+t.Name+"]"),
		colorLogError.Sprint(message),
	)
}
