package contracts

import (
	"context"
	"net/http"
)

type WebFrameworkContext interface {
	// SetSameSite with cookie
	SetSameSite(sameSite http.SameSite)

	// SetCookie adds a Set-Cookie header to the ResponseWriter's headers.
	// The provided cookie must have a valid Name.
	// Invalid cookies may be silently dropped.
	SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool)

	// Error attaches an error to the current context.
	// The error is pushed to a list of errors.
	// It's a good idea to call Error for each error that occurred during the resolution of a request.
	// A middleware can be used to collect all the errors and push them to a database together, print a log, or append it in the HTTP response.
	// Error will panic if err is nil.
	Error(err error) error

	// Abort prevents pending handlers from being called.
	// Note that this will not stop the current handler.
	// Let's say you have an authorization middleware that validates that the current request is authorized.
	// If the authorization fails (ex: the password does not match), call Abort to ensure the remaining handlers for this request are not called.
	Abort()

	// Cookie returns the named cookie provided in the request or ErrNoCookie if not found.
	// And return the named cookie is unescaped.
	// If multiple cookies match the given name, only one cookie will be returned.
	Cookie(name string) (string, error)

	// Set is used to store a new key/value pair exclusively for this context.
	// It also lazy initializes c.Keys if it was not used previously.
	Set(key string, value any)

	// Get returns the value for the given key, ie: (value, true).
	// If the value does not exist it returns (nil, false)
	Get(key string) (value any, exists bool)

	// Next should be used only inside middleware. It executes the pending handlers in the chain inside the calling handler. See example in GitHub.
	Next()

	// Request returns the underlying *http.Request object.
	Request() *http.Request

	context.Context
}
