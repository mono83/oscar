package out

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/mono83/oscar/core/events"
	"io"
	"time"
)

// FullRealTimePrinter returns events receiver, used to print test case flow
func FullRealTimePrinter(stream io.Writer, showSetValue bool) func(interface{}) {
	return func(i interface{}) {
		if s, ok := i.(events.Start); ok {
			print(
				stream,
				fmt.Sprintf("Starting %s %s", s.Type, s.Name),
				colorLogTestCase,
			)
		} else if s, ok := i.(events.Finish); ok {
			if s.Error == nil {
				print(
					stream,
					fmt.Sprintf("Sucessfully finished %s %s", s.Type, s.Name),
					colorLogTestCase,
				)
			} else {
				print(
					stream,
					fmt.Sprintf("%s %s completed with error", s.Type, s.Name),
					colorLogError,
				)
			}
		} else if e, ok := i.(events.LogEvent); ok {
			c := colorLogDebug
			if e.Level == 2 {
				c = colorLogInfo
			}
			print(stream, e.Pattern, c)
		} else if s, ok := i.(events.SetVar); ok && showSetValue {
			prev := ""
			if s.Previous != nil && *s.Previous != s.Value {
				prev = " previous value was " + *s.Previous
			}

			print(stream, fmt.Sprintf("Setting %s := %s%s", s.Key, s.Value, prev), colorLogDebug)
		}
	}
}

var colorLogTime = color.New(color.FgWhite)

var colorLogDebug = color.New(color.FgHiBlack)
var colorLogInfo = color.New(color.FgCyan)
var colorLogTestCase = color.New(color.FgHiGreen)
var colorLogError = color.New(color.FgHiYellow)

func print(stream io.Writer, message string, c *color.Color) {
	fmt.Fprintf(
		stream,
		"%s %s\n",
		colorLogTime.Sprint(time.Now().Format("15:04:05")),
		c.Sprint(message),
	)
}
