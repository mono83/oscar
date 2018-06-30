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

	// Building and attaching test runtime data
	rt := &RuntimeData{}
	ctx.Register(rt.BuildListener())

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
				suiteInitFailed = true
			} else {
				// Copying variables
				suiteContext.Import(caseContext.values)
			}
		}

		// Iterating over test cases
		for _, c := range suite.GetCases() {
			cid, cname := c.ID()

			caseContext := suiteContext.Fork(cid)

			caseContext.Emit(events.Start{Type: "TestCase", Name: cname})
			var err error
			if suiteInitFailed {
				err = Skip{
					Failed:  SuiteSetUp,
					Skipped: cname,
				}
				caseContext.Emit(events.AssertDone{Error: err})
			} else {
				if deps := c.GetDependsOn(); len(deps) > 0 {
					for _, d := range deps {
						if !rt.IsCompletedSuccessfully(d) {
							err = Skip{
								Failed:  rt.GetName(d),
								Skipped: cname,
							}
							break
						}
					}
				}
				if err == nil {
					err = c.Assert(caseContext)
				}
			}
			caseContext.Emit(events.Finish{Type: "TestCase", Name: cname, Error: err})
		}
		suiteContext.Emit(events.Finish{Type: "TestSuite", Name: sname})
	}

	ctx.Wait()

	if rt.TotalErrorsCount > 0 {
		return fmt.Errorf("%d error(s) encountered", rt.TotalErrorsCount)
	}

	return nil
}
