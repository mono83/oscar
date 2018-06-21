package events

// AssertDone is an event, emitted when assertion is done. Error is optional
type AssertDone struct {
	Error error
}
