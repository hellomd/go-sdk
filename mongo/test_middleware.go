package mongo

import (
	"net/http"

	mgo "gopkg.in/mgo.v2"
)

// NewTestMiddleware -
func NewTestMiddleware(db *mgo.Database) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		ctx := SetInCtx(r.Context(), db)
		next(w, r.WithContext(ctx))
	}
}
