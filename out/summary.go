package out

import (
	"io"
	"time"

	"github.com/mono83/table"
	"github.com/mono83/table/cells"
)

// PrintSummary outputs summary table into provided stream using provided report data
func PrintSummary(stream io.Writer, report *Report) {
	table.Print(summaryTable{report}, table.PrintOptions{ColumnSeparator: "  ", HeaderBorder: "-", Writer: stream})
}

type summaryTable struct {
	*Report
}

func (summaryTable) Headers() []string {
	return []string{
		"",
		"Name",
		"Fail",
		"Success",
		"Requests",
		"Time total",
		"HTTP",
		"Sleep",
	}
}

func (s summaryTable) EachRow(f func(...table.Cell)) {
	for _, node := range s.Flatten() {
		success, failed := node.CountAssertionsRecursive()

		var typeCell table.Cell = cells.Empty{}
		if "TestSuite" == node.Type || "TestSuiteInit" == node.Type {
			typeCell = cells.ColoredWhite(cells.AlignCenter(cells.String("Suite")))
		}

		var nameCell table.Cell
		nameCell = cells.String(node.Name)
		if "TestSuiteInit" == node.Type {
			nameCell = cells.String(" ~ SetUp ~ ")
		}

		if failed > 0 && node.IsSkip {
			nameCell = cells.ColoredMagentaHi(nameCell)
		} else if failed > 0 {
			nameCell = cells.ColoredRedHi(nameCell)
		} else if "TestSuiteInit" == node.Type {
			nameCell = cells.ColoredWhite(nameCell)
		} else {
			nameCell = cells.ColoredGreenHi(nameCell)
		}

		f(
			typeCell,
			nameCell,
			zeroEmptyCell(failed),
			zeroEmptyCell(success),
			zeroEmptyCell(node.CountRemoteRequestsRecursive()),
			elapsedCell(node.Elapsed()),
			elapsedCell(node.ElapsedRemoteRecursive()),
			elapsedCell(node.ElapsedSleepRecursive()),
		)
	}
}

func zeroEmptyCell(i int) table.Cell {
	if i == 0 {
		return cells.ColoredWhite(cells.AlignRight(cells.String("-")))
	}

	return cells.Int(i)
}

func elapsedCell(t time.Duration) table.Cell {
	if t.Nanoseconds() == 0 {
		return cells.Empty{}
	}
	if t.Seconds() < 0.001 {
		return cells.ColoredWhite(cells.AlignRight(cells.String("< 1ms")))
	}

	return cells.Duration(t)
}
