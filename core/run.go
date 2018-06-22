package core

import (
	"errors"
	"fmt"
	"github.com/mono83/oscar/core/events"
)

// RunSequential runs all provided cases sequentially
func RunSequential(ctx *Context, suites []Suite) error {
	if len(suites) == 0 {
		return errors.New("empty suites list")
	}

	var errorsCnt int

	for _, suite := range suites {
		suiteContext := ctx.Fork()
		sid, sname := suite.ID()
		suiteContext.Emit(events.Start{Type: "TestSuite", ID: sid, Name: sname})
		for _, c := range suite.GetCases() {
			// Forking context for test case
			caseContext := suiteContext.Fork()

			cid, cname := c.ID()

			caseContext.Emit(events.Start{Type: "TestCase", ID: cid, Name: cname})
			err := c.Assert(suiteContext.Fork())
			caseContext.Emit(events.Finish{Type: "TestCase", ID: cid, Name: cname, Error: err})
			if err != nil {
				errorsCnt++
			}
		}
		suiteContext.Emit(events.Finish{Type: "TestSuite", ID: sid, Name: sname})
	}

	ctx.Wait()

	if errorsCnt > 0 {
		return fmt.Errorf("%d error(s) encountered", errorsCnt)
	}

	return nil
}
