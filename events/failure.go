package events

// Failure is special event (compatible with error),
// that contains assertion or runtime error message
type Failure string

func (f Failure) Error() string {
	return string(f)
}

// IfFailure checks, if provided interface is instance of Failure event
// and if so, passes in into callback func
func IfFailure(e *Emitted, f func(Failure)) {
	if e != nil && e.Data != nil && f != nil {
		if a, ok := e.Data.(Failure); ok {
			f(a)
		}
	}
}
