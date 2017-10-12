package out

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/mono83/oscar"
	"io"
)

// GetAftermath returns aftermath event printer
func GetAftermath(stream io.Writer) func(interface{}) {
	return func(i interface{}) {
		if s, ok := i.(oscar.FinishEvent); ok {
			if o, ok := s.Owner.(*oscar.TestSuite); ok {
				err := o.GetError()
				if err != nil {
					// Printing error details
					fmt.Fprintln(stream)
					fmt.Fprintln(stream, " Errors:")
					i := 1
					for _, s := range o.Cases {
						if s.Error != nil {
							fmt.Fprintf(stream, "  %d. %s\n", i, s.Name)
							fmt.Fprintln(stream, "     ", s.Error)
							for k, v := range s.Vars {
								fmt.Fprintln(stream, "      ", k, ":=", v)
							}
							fmt.Fprintln(stream)
							i++
						}
					}
					fmt.Fprintln(stream)
				}

				// Building global aftermath
				longest := len("Test suite")
				for _, s := range o.Cases {
					if s.Error != nil || s.CountAssertSuccess > 0 {
						if l := len(s.Name); l > longest {
							longest = l
						}
					}
				}

				namePattern := fmt.Sprintf(" %%-%ds", longest)
				fullPattern := "%s" + namePattern + "  %5d   %5d     %5d   %7.1fms\n"

				fmt.Fprintln(stream)
				fmt.Fprintln(stream)
				fmt.Fprintf(
					stream,
					"      "+namePattern+" Success  Failed  Requests  Time spent\n",
					"Test suite",
				)
				fmt.Fprintln(stream)

				for _, s := range o.Cases {
					if s.Error != nil || s.CountAssertSuccess > 0 {
						status := colorOscarSummarySuccess.Sprint("  OK  ")
						errorCnt := 0
						if s.Error != nil {
							errorCnt = 1
							status = colorOscarSummaryFailed.Sprint(" FAIL ")
						}

						fmt.Fprintf(
							stream,
							fullPattern,
							status,
							s.Name,
							s.CountAssertSuccess,
							errorCnt,
							s.CountRemoteRequests,
							s.Elapsed().Seconds()*1000,
						)
					}
				}

				fmt.Fprintln(stream)
				fmt.Fprintln(stream)
			}
		}
	}
}

var colorOscarSummarySuccess = color.New(color.FgHiGreen)
var colorOscarSummaryFailed = color.New(color.FgHiRed)
