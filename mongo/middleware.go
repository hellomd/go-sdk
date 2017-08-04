package mongo

import (
	"net/http"

	mgo "gopkg.in/mgo.v2"
)

// Middleware -
type Middleware interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type middleware struct {
	session *mgo.Session
	dbName  string
}

// NewMiddleware -
func NewMiddleware(mongoURL string, useSSL bool) (Middleware, error) {
	s, dbName, err := CreateSession(mongoURL, useSSL)
	if err != nil {
		return nil, err
	}

	return &middleware{s, dbName}, nil
}

func (mw *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	s := mw.session.Copy()
	context := SetInCtx(r.Context(), s.DB(mw.dbName))
	next(w, r.WithContext(context))
	defer s.Close()
}
