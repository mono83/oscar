package out

import (
	"fmt"
	"github.com/fatih/color"
	"io"
)

// PrintTestCaseErrorsSummary outputs errors summary
func PrintTestCaseErrorsSummary(stream io.Writer, report *Report) {
	i := 0
	for _, node := range report.Flatten() {
		if node.Type == "TestCase" && node.Error != nil {
			i++
			colorErrorSummaryErrorLine.Fprintf(stream, "%d. Error in \"%s\"\n", i, node.Name)
			fmt.Fprintln(stream, *node.Error)
			if len(node.Variables) > 0 {
				colorErrorSummaryVars.Fprintln(stream, "Variables dump:")
				for k, v := range node.Variables {
					colorErrorSummaryVars.Fprintf(stream, "%s := ", k)
					fmt.Fprintln(stream, v)
				}
			}
			fmt.Fprintln(stream, "")
		}
	}
}

var colorErrorSummaryErrorLine = color.New(color.FgHiRed)
var colorErrorSummaryVars = color.New(color.FgHiBlack)
