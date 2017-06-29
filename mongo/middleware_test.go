package mongo

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/urfave/negroni"

	"gopkg.in/mgo.v2"
)

func TestMiddleware(t *testing.T) {
	dbName := "test"
	mw := NewMiddleware(os.Getenv("MONGO_URL"), dbName).(*middleware)
	req := httptest.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	a := negroni.New()
	a.Use(mw)
	a.Use(negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		v, ok := r.Context().Value(MongoCtxKey).(*mgo.Database)
		if !ok {
			t.Error("Expected database on context")
		}
		if v.Name != dbName {
			t.Errorf("Expected %s database name, got: %s", dbName, v.Name)
		}
		if v.Session == mw.session {
			t.Errorf("Expected session to be different from stored session")
		}
		next(w, r)
	}))

	a.ServeHTTP(response, req)

}

func TestGetMongoFromContext(t *testing.T) {
	ctx := context.Background()
	_, err := GetMongoFromCtx(ctx)
	if err != ErrNoMongoInCtx {
		t.Error("Expected ErrNoMongoInCtx, got: ", err)
	}

	expectedDB := &mgo.Database{}
	ctx = context.WithValue(ctx, MongoCtxKey, expectedDB)
	db, _ := GetMongoFromCtx(ctx)
	if db != expectedDB {
		t.Errorf("Expected %v, got: %v", expectedDB, db)
	}
}
