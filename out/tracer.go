package out

import (
	"fmt"
	"github.com/mono83/oscar"
	"io"
)

// GetTracer returns event receiver, used to print tracing information
func GetTracer(stream io.Writer) func(interface{}) {
	return func(e interface{}) {
		if t, ok := e.(oscar.TraceEvent); ok {
			fmt.Fprintln(stream, string(t))
		}
	}
}
