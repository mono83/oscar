package oscar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext_Vars(t *testing.T) {
	// Creating context
	ctx := NewContext()
	// Creating child context
	cld := ctx.Fork(42)

	ctx.Set("foo", "bar")
	assert.Equal(t, "bar", ctx.Get("foo"))
	assert.Equal(t, "bar", cld.Get("foo"))

	cld.Set("foo", "notbar")
	assert.Equal(t, "bar", ctx.Get("foo"))
	assert.Equal(t, "notbar", cld.Get("foo"))
}

func TestContext_ExportVars(t *testing.T) {
	// Creating context
	ctx := NewContext()
	// Creating child context
	cld := ctx.Fork(42)
	cld2 := ctx.Fork(300)

	ctx.SetExport("foo", "bar")
	assert.Equal(t, "bar", ctx.Get("foo"))
	assert.Equal(t, "bar", cld.Get("foo"))
	assert.Equal(t, "bar", cld2.Get("foo"))

	cld2.SetExport("foo", "notbar")
	assert.Equal(t, "notbar", ctx.Get("foo"))
	assert.Equal(t, "notbar", cld.Get("foo"))
	assert.Equal(t, "notbar", cld2.Get("foo"))
}