package events

// RegistrationBegin event emitted on test startup, when test/case/suite is scheduled
type RegistrationBegin struct {
	ID   int
	Type string
	Name string
}

// RegistrationEnd event emitted on test startup, when test/case/suite is scheduled
type RegistrationEnd struct {
	Type string
	Name string
}
