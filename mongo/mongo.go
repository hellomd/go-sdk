package mongo

import (
	"context"
	"net/http"

	mgo "gopkg.in/mgo.v2"
)

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
func GetMongoFromContext(ctx context.Context) *mgo.Database {
	db, ok := ctx.Value(MongoContextKey).(*mgo.Database)
	if !ok {
		panic("Could not lookup mongo from context")
	}
	return db
}

// ServeHTTP copies the db session, adds it to the request context
// Closes the db session on defer
func (mw *mongo) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	s := mw.session.Copy()
	next(w, r.WithContext(context.WithValue(r.Context(), MongoContextKey, s.DB(mw.dbName))))
	defer s.Close()
}
