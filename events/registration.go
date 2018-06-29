package events

// RegistrationIn event emitted on test startup, when test/case/suite is scheduled
type RegistrationIn struct {
	ID   int
	Type string
	Name string
}

// RegistrationOut event emitted on test startup, when test/case/suite is scheduled
type RegistrationOut struct {
	ID   int
	Type string
	Name string
}
