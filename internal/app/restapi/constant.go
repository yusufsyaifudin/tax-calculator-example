package restapi

import (
	"context"
)

// content types
const (
	ContentTypeJSON     = "application/json"
	ContentTypePostForm = "application/x-www-form-urlencoded"
)

// Handler represents an api handler
type Handler func(context.Context, Request) Response
type Middleware func(Handler) Handler

// Implementing this idea https://hackernoon.com/simple-http-middleware-with-go-79a4ad62889b
func ChainMiddleware(mw ...Middleware) Middleware {
	return func(final Handler) Handler {
		return func(ctx context.Context, req Request) Response {
			last := final
			for i := len(mw) - 1; i >= 0; i-- {
				last = mw[i](last)
			}

			// last middleware
			return last(ctx, req)
		}
	}
}
