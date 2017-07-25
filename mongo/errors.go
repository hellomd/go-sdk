package mongo

import "errors"

var errNotInCtx = errors.New("no mongo in context")

// IsNotInCtxError -
func IsNotInCtxError(err error) bool {
	return err == errNotInCtx
}
