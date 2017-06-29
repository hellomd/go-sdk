package mongo

import (
	"context"
	"errors"
	"net/http"

	mgo "gopkg.in/mgo.v2"
)

// ErrNoMongoInCtx -
var ErrNoMongoInCtx = errors.New("No mongo in context")

type mongoCtxKey int

// MongoCtxKey -
const MongoCtxKey mongoCtxKey = 0

// Middleware -
type Middleware interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type middleware struct {
	session *mgo.Session
	dbName  string
}

// NewMiddleware -
func NewMiddleware(mongoURL string, dbName string) Middleware {
	s, _ := mgo.Dial(mongoURL)
	return &middleware{s, dbName}
}

// GetMongoFromCtx -
func GetMongoFromCtx(ctx context.Context) (*mgo.Database, error) {
	db, ok := ctx.Value(MongoCtxKey).(*mgo.Database)
	if !ok {
		return nil, ErrNoMongoInCtx
	}
	return db, nil
}

// ServeHTTP copies the db session, adds it to the request context
// Closes the db session on defer
func (mw *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	s := mw.session.Copy()
	next(w, r.WithContext(context.WithValue(r.Context(), MongoCtxKey, s.DB(mw.dbName))))
	defer s.Close()
}
