package revel

import (
	"runtime/debug"
)

// PanicFilter wraps the action invocation in a protective defer blanket that
// converts panics into 500 error pages.
type PanicFilter struct{}

func (f PanicFilter) Call(c *Controller, fc FilterChain) {
	defer func() {
		if err := recover(); err != nil {
			handleInvocationPanic(c, err)
		}
	}()
	fc.Call(c)
}

// This function handles a panic in an action invocation.
// It cleans up the stack trace, logs it, and displays an error page.
func handleInvocationPanic(c *Controller, err interface{}) {
	error := NewErrorFromPanic(err)
	if error == nil {
		ERROR.Print(err, "\n", string(debug.Stack()))
		c.Response.Out.WriteHeader(500)
		c.Response.Out.Write(debug.Stack())
		return
	}

	ERROR.Print(err, "\n", error.Stack)
	c.RenderError(error).Apply(c.Request, c.Response)
}