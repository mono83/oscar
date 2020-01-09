package oscar

import (
	"fmt"
	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/mono83/oscar/events"
)

// RuntimeData contains data about current test run
type RuntimeData struct {
	m sync.Mutex

	// All errors been caught
	Errors []error
	// Names contains map of test case/suite names
	Names map[int]string
	// Invocations contains flags for success/failure invocations
	Invocations map[int]bool
}

// BuildListener builds and returns event listener
func (r *RuntimeData) BuildListener() func(emitted *events.Emitted) {
	r.Names = map[int]string{}
	r.Invocations = map[int]bool{}

	er := events.EventRouter{
		Assert: func(done events.AssertDone, emitted *events.Emitted) {
			if done.Error != nil {
				r.m.Lock()
				defer r.m.Unlock()
				r.Errors = append(r.Errors, done.Error)
			}
		},
		Start: func(reg events.Start, em *events.Emitted) {
			r.m.Lock()
			defer r.m.Unlock()
			r.Names[em.OwnerID] = reg.Name
		},
		Finish: func(f events.Finish, em *events.Emitted) {
			r.m.Lock()
			defer r.m.Unlock()
			r.Invocations[em.OwnerID] = f.Error == nil
		},
	}

	return er.OnEvent
}

// GetName resolves and returns name by ID
func (r *RuntimeData) GetName(id int) string {
	if len(r.Names) > 0 {
		r.m.Lock()
		name, ok := r.Names[id]
		r.m.Unlock()

		if ok {
			return name
		}
	}

	return "undefined"
}

// IsCompletedSuccessfully returns true only if test entity with requested ID was invoked and
// actually completed successfully
func (r *RuntimeData) IsCompletedSuccessfully(id int) bool {
	if len(r.Invocations) > 0 {
		r.m.Lock()
		success, ok := r.Invocations[id]
		r.m.Unlock()

		return ok && success
	}

	return false
}

// GetErrors returns all errors as single wrapper
func (r *RuntimeData) GetErrors() error {
	// Building multierror
	me := multierror.Error{
		Errors: r.Errors,
		ErrorFormat: func(errors []error) string {
			return fmt.Sprintf("%d error(s) encountered", len(errors))
		},
	}

	return me.ErrorOrNil()
}
