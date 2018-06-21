package events

// SetVar is TestContext event, emitted when variables changes it's value
type SetVar struct {
	Key      string
	Value    string
	Previous *string
}
