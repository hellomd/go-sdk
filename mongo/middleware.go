package mongo

import (
	"context"
	"net/http"

	mgo "gopkg.in/mgo.v2"
)

// Middleware -
type Middleware struct {
	session *mgo.Session
	dbName  string
}

// NewMiddleware -
func NewMiddleware(mongoURL string, useSSL bool) (*Middleware, error) {
	s, dbName, err := createSession(mongoURL, useSSL)
	if err != nil {
		return nil, err
	}

	return &Middleware{s, dbName}, nil
}

func (mw *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	mw.UseSession(r.Context(), func(ctx context.Context) {
		next(w, r.WithContext(ctx))
	})
}

func (mw *Middleware) UseSession(ctx context.Context, next func(context.Context)) {
	s := mw.session.Copy()
	defer s.Close()
	ctx = SetInCtx(ctx, s.DB(mw.dbName))
	next(ctx)
}
