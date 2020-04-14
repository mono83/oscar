package lua

import (
	"errors"

	"github.com/mono83/oscar"
	"github.com/mono83/oscar/events"
	"github.com/yuin/gopher-lua"
)

// TestCaseMeta contains metatable name for userdata structure TestCase in lua
const TestCaseMeta = "TestCaseType"

// SuiteFromFiles builds suite using Lua sources file
func SuiteFromFiles(c *oscar.Context, files ...string) (oscar.Suite, error) {
	if len(files) == 0 {
		return nil, errors.New("empty files list to load")
	}

	// Building Lua state
	L := lua.NewState()

	// Building test suite
	s := &testSuite{
		id:    id(),
		name:  files[len(files)-1],
		state: L,
	}

	registered := map[string]int{}

	ctx := c.Fork(s.id)
	ctx.Register(func(emitted *events.Emitted) {
		if b, ok := emitted.Data.(events.RegistrationBegin); ok {
			registered[b.Name] = b.ID
		}
	})

	// Injecting module into Lua runtime
	injectModule(s, ctx, L)

	// Emitting registration start event
	ctx.Emit(events.RegistrationBegin{Type: "TestSuite", ID: s.id, Name: s.name})

	// Reading files sequentially
	for _, file := range files {
		if err := L.DoFile(file); err != nil {
			return nil, err
		}
	}

	// Resolving dependencies
	for _, tc := range s.cases {
		if len(tc.deps) > 0 {
			for _, dep := range tc.deps {
				f, ok := registered[dep]
				if !ok {
					return nil, errors.New("unable to find dependency " + dep + " for " + tc.name)
				}
				tc.dep = append(tc.dep, f)
			}
		}
	}

	// Emitting registration done event
	ctx.Emit(events.RegistrationEnd{Type: "TestSuite", Name: s.name})

	return s, nil
}
