package oscar

import (
	"errors"
	"fmt"
	"github.com/mono83/oscar/events"
)

// RunSequential runs all provided cases sequentially
func RunSequential(ctx *Context, suites []Suite) error {
	if len(suites) == 0 {
		return errors.New("empty suites list")
	}

	var errorsCnt int

	for _, suite := range suites {
		sid, sname := suite.ID()
		suiteContext := ctx.Fork(sid)
		suiteContext.Emit(events.Start{Type: "TestSuite", Name: sname})

		// Running INIT func
		suiteInitFailed := false
		if c := suite.GetSetUp(); c != nil {
			cid, cname := c.ID()

			caseContext := suiteContext.Fork(cid)

			caseContext.Emit(events.Start{Type: "TestSuiteInit", Name: cname})
			err := c.Assert(caseContext)
			caseContext.Emit(events.Finish{Type: "TestSuiteInit", Name: cname, Error: err})
			if err != nil {
				errorsCnt++
				suiteInitFailed = true
			} else {
				// Copying variables
				suiteContext.Import(caseContext.values)
			}
		}

		// Iterating over test cases
		if !suiteInitFailed {
			for _, c := range suite.GetCases() {
				cid, cname := c.ID()

				caseContext := suiteContext.Fork(cid)

				caseContext.Emit(events.Start{Type: "TestCase", Name: cname})
				err := c.Assert(caseContext)
				caseContext.Emit(events.Finish{Type: "TestCase", Name: cname, Error: err})
				if err != nil {
					errorsCnt++
				}
			}
		}
		suiteContext.Emit(events.Finish{Type: "TestSuite", Name: sname})
	}

	ctx.Wait()

	if errorsCnt > 0 {
		return fmt.Errorf("%d error(s) encountered", errorsCnt)
	}

	return nil
}
