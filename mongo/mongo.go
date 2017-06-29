package mongo

import (
	"context"
	"errors"
	"net/http"

	mgo "gopkg.in/mgo.v2"
)

// ErrNoMongoInCtx -
var ErrNoMongoInCtx = errors.New("No mongo in context")

type mongoContextKey int

// MongoContextKey -
const MongoContextKey mongoContextKey = 0

// Mongo -
type Mongo interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type mongo struct {
	session *mgo.Session
	dbName  string
}

// NewMongo -
func NewMongo(mongoURL string, dbName string) Mongo {
	s, _ := mgo.Dial(mongoURL)
	return &mongo{s, dbName}
}

// GetMongoFromContext -
func GetMongoFromContext(ctx context.Context) (*mgo.Database, error) {
	db, ok := ctx.Value(MongoContextKey).(*mgo.Database)
	if !ok {
		return nil, ErrNoMongoInCtx
	}
	return db, nil
}

// ServeHTTP copies the db session, adds it to the request context
// Closes the db session on defer
func (mw *mongo) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	s := mw.session.Copy()
	next(w, r.WithContext(context.WithValue(r.Context(), MongoContextKey, s.DB(mw.dbName))))
	defer s.Close()
}
