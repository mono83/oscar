package out

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/mono83/oscar"
	"io"
	"time"
)

// GetTestCasePrinter returns events receiver, used to print test case flow
func GetTestCasePrinter(stream io.Writer) func(interface{}) {
	return func(i interface{}) {
		if s, ok := i.(oscar.StartEvent); ok {
			if t, ok := s.Owner.(*oscar.TestCase); ok {
				if t.Name == oscar.InitFuncName {
					print(stream, t, "Running test case initialization func", colorLogTestCase)
				} else {
					print(stream, t, "Running test case "+t.Name, colorLogTestCase)
				}
			}
		} else if s, ok := i.(oscar.FinishEvent); ok {
			if t, ok := s.Owner.(*oscar.TestCase); ok {
				if t.Error != nil {
					print(stream, t, fmt.Sprintf("Test case failed. Success: %d, failed 1", t.CountAssertSuccess), colorLogError)
				} else {
					elapsed, _, _ := t.Elapsed()
					print(stream, t, fmt.Sprintf("Test case done in %.2f sec. Assertions: %d", elapsed.Seconds(), t.CountAssertSuccess), colorLogTestCase)
				}
			}
		} else if e, ok := i.(oscar.TestLogEvent); ok {
			c := colorLogDebug
			if e.Level == 1 {
				c = colorLogInfo
			}
			print(stream, e.Owner, e.Message, c)
		}
	}
}

var colorLogTime = color.New(color.FgWhite)

var colorLogDebug = color.New(color.FgHiBlack)
var colorLogInfo = color.New(color.FgCyan)
var colorLogTestCase = color.New(color.FgHiGreen)
var colorLogError = color.New(color.FgHiYellow)

func print(stream io.Writer, t *oscar.TestCase, message string, c *color.Color) {
	fmt.Fprintf(
		stream,
		"%s %s\n",
		colorLogTime.Sprint(time.Now().Format("15:04:05")),
		c.Sprint(t.Interpolate(message)),
	)
}
