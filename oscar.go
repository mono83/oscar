package oscar

import (
	"errors"
	"path"
	"regexp"
	"time"
)

// Oscar is main test suite runner
type Oscar struct {
	*TestContext

	err error

	Suits            []*TestSuite
	suiteDefinitions []struct {
		Header       string
		FileName     string
		CaseSelector func(*TestCase) bool
	}
}

// AddTestSuiteFile registers new test suite
func (o *Oscar) AddTestSuiteFile(name, header, filter string) error {
	suite := struct {
		Header       string
		FileName     string
		CaseSelector func(*TestCase) bool
	}{
		Header:   header,
		FileName: name,
	}

	if len(filter) > 0 {
		nameMatcher, err := regexp.Compile("(?i)" + filter)
		if err != nil {
			return err
		}
		suite.CaseSelector = func(tc *TestCase) bool {
			return nameMatcher.MatchString(tc.Name)
		}
	}

	o.suiteDefinitions = append(o.suiteDefinitions, suite)
	return nil
}

// Start starts tests
func (o *Oscar) Start() error {
	if len(o.suiteDefinitions) == 0 {
		return errors.New("no suites configured")
	}

	// Building test suite objects
	o.Suits = make([]*TestSuite, len(o.suiteDefinitions))
	for i, def := range o.suiteDefinitions {
		o.Suits[i] = &TestSuite{
			TestContext: &TestContext{
				Parent: o.TestContext,
			},
			CaseSelector: def.CaseSelector,
		}
		if len(def.Header) > 0 {
			o.Suits[i].Include = []string{def.Header}
		}
	}

	// Starting test suits
	for i, suite := range o.Suits {
		if err := suite.StartFile(o.suiteDefinitions[i].FileName); err != nil && o.err == nil {
			o.err = err
		}
	}

	o.Emit(FinishEvent{Time: time.Now(), Owner: o})
	return o.err
}

// prefix generates test case prefix
func (o Oscar) prefix(i int) string {
	name := ""
	if len(o.Suits) > 1 {
		name = path.Base(o.suiteDefinitions[i].FileName) + ":"
	}

	return name
}

// GetError returns main error
func (o Oscar) GetError() error { return o.err }

// IterateResults iterates over all test cases, passing results to provided callback
func (o Oscar) IterateResults(f func(string, int, int, int, time.Duration, time.Duration, time.Duration)) {
	for i, ts := range o.Suits {
		for _, tc := range ts.GetCases() {
			cntErr := 0
			if tc.Error != nil {
				cntErr = 1
			}
			if cntErr == 0 && tc.CountAssertSuccess == 0 {
				continue
			}
			elapsedTotal, elapsedHTTP, elapsedSleep := tc.Elapsed()

			f(
				o.prefix(i)+tc.Name,
				tc.CountAssertSuccess,
				cntErr,
				tc.CountRemoteRequests,
				elapsedTotal,
				elapsedHTTP,
				elapsedSleep,
			)
		}
	}
}

// IterateErrors iterates over all test cases, passing error information to callback
func (o Oscar) IterateErrors(f func(*TestContext, string, error)) {
	for i, ts := range o.Suits {
		if ts.Error == nil {
			continue
		}

		for _, tc := range ts.GetCases() {
			if tc.Error != nil {
				f(tc.TestContext, o.prefix(i)+tc.Name, tc.Error)
			}
		}
	}
}
