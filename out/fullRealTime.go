package out

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/mono83/oscar/events"
	"io"
	"time"
)

// FullRealTimePrinter returns events receiver, used to print test case flow
func FullRealTimePrinter(stream io.Writer, showSetValue bool, showTrace bool) func(*events.Emitted) {
	switcher := events.EventRouter{
		Start: func(start events.Start, _ *events.Emitted) {
			if start.Type == "TestSuite" && !showTrace {
				return
			}
			print(
				stream,
				fmt.Sprintf("Starting %s %s", start.Type, start.Name),
				colorLogTestCase,
			)
		},
		Finish: func(finish events.Finish, _ *events.Emitted) {
			if finish.Type == "TestSuite" && !showTrace {
				return
			}
			if finish.Error == nil {
				print(
					stream,
					fmt.Sprintf("Sucessfully finished %s %s", finish.Type, finish.Name),
					colorLogTestCase,
				)
			} else {
				print(
					stream,
					fmt.Sprintf("%s \"%s\" completed with an error", finish.Type, finish.Name),
					colorLogError,
				)
			}
		},
		Log: func(log events.LogEvent, _ *events.Emitted) {
			if log.Level == 0 && !showTrace {
				return
			}
			c := colorLogDebug
			if log.Level == 2 {
				c = colorLogInfo
			}
			print(stream, log.Pattern, c)
		},
	}

	if showSetValue {
		switcher.Var = func(s events.SetVar, _ *events.Emitted) {
			prev := ""
			if s.Previous != nil && *s.Previous != s.Value {
				prev = " previous value was " + *s.Previous
			}

			print(stream, fmt.Sprintf("Setting %s := %s%s", s.Key, s.Value, prev), colorLogDebug)
		}
	}

	return switcher.OnEvent
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
