package mongo

import (
	"context"
	"testing"

	mgo "gopkg.in/mgo.v2"
)

func TestContext(t *testing.T) {
	ctx := context.Background()
	_, err := GetFromCtx(ctx)
	if err != errNotInCtx {
		t.Error("Expected errNotInCtx, got: ", err)
	}

	expectedDB := &mgo.Database{}
	ctx = SetInCtx(ctx, expectedDB)
	db, _ := GetFromCtx(ctx)
	if db != expectedDB {
		t.Errorf("Expected %v, got: %v", expectedDB, db)
	}
}
