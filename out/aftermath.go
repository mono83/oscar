package out

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/mono83/oscar"
	"io"
	"time"
)

// GetAftermath returns aftermath event printer
func GetAftermath(stream io.Writer) func(interface{}) {
	return func(i interface{}) {
		if s, ok := i.(oscar.FinishEvent); ok {
			if o, ok := s.Owner.(*oscar.Oscar); ok {
				err := o.GetError()
				if err != nil {
					// Printing error details
					fmt.Fprintln(stream)
					fmt.Fprintln(stream, " Errors:")
					i := 1
					o.IterateErrors(func(c *oscar.TestContext, name string, err error) {
						fmt.Fprintf(stream, "  %d. %s\n", i, name)
						fmt.Fprintln(stream, "     ", err)
						for k, v := range c.Vars {
							fmt.Fprintln(stream, "      ", k, ":=", v)
						}
						fmt.Fprintln(stream)
						i++
					})

					fmt.Fprintln(stream)
				}

				// Building global aftermath
				longest := len("Test suite")
				o.IterateResults(func(name string, success int, err int, remote int, elapsedTotal time.Duration, elapsedHTTP time.Duration, elapsedSleep time.Duration) {
					if l := len(name); l > longest {
						longest = l
					}
				})

				namePattern := fmt.Sprintf(" %%-%ds", longest)
				fullPattern := "%s" + namePattern + "  %5d   %5d     %5d   %7.1fms  %7.1fms  %7.1fms\n"

				fmt.Fprintln(stream)
				fmt.Fprintln(stream)
				fmt.Fprintf(
					stream,
					"      "+namePattern+" Success  Failed  Requests  Total time     HTTP       Sleep\n",
					"Test suite",
				)
				fmt.Fprintln(stream)

				o.IterateResults(func(name string, success int, err int, remote int, elapsedTotal time.Duration, elapsedHTTP time.Duration, elapsedSleep time.Duration) {
					status := colorOscarSummarySuccess.Sprint("  OK  ")
					if err > 0 {
						status = colorOscarSummaryFailed.Sprint(" FAIL ")
					}

					fmt.Fprintf(
						stream,
						fullPattern,
						status,
						name,
						success,
						err,
						remote,
						elapsedTotal.Seconds()*1000,
						elapsedHTTP.Seconds()*1000,
						elapsedSleep.Seconds()*1000,
					)
				})

				fmt.Fprintln(stream)
				fmt.Fprintln(stream)
			}
		}
	}
}

var colorOscarSummarySuccess = color.New(color.FgHiGreen)
var colorOscarSummaryFailed = color.New(color.FgHiRed)
