package out

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/mono83/oscar/core/events"
	"io"
	"sync"
)

// DotRealTimePrinter returns events receiver, used to print test case flow
func DotRealTimePrinter(stream io.Writer) func(interface{}) {
	cnt := 0
	max := 40
	m := sync.Mutex{}

	print := func(s rune, c *color.Color) {
		str := string(s)
		if c != nil {
			str = c.Sprint(str)
		}

		m.Lock()
		fmt.Fprint(stream, str)
		cnt++
		if cnt == max {
			fmt.Fprintln(stream)
			cnt = 0
		}
		m.Unlock()
	}

	return func(i interface{}) {
		if a, ok := i.(events.AssertDone); ok {
			if a.Error == nil {
				print('.', colorDotOK)
			} else {
				print('E', colorDotErr)
			}
		} else if _, ok := i.(events.Start); ok {
			print('<', colorDotSF)
		} else if _, ok := i.(events.Finish); ok {
			print('>', colorDotSF)
		} else if _, ok := i.(events.Sleep); ok {
			print('z', colorDotSleep)
		} else if _, ok := i.(events.RemoteRequest); ok {
			print('^', colorDotRemote)
		}
	}
}

var colorDotSF = color.New(color.FgBlack)
var colorDotOK = color.New(color.FgGreen)
var colorDotErr = color.New(color.FgRed)
var colorDotSleep = color.New(color.FgBlue)
var colorDotRemote = color.New(color.FgHiCyan)
