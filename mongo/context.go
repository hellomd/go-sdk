package mongo

import (
	"context"

	mgo "gopkg.in/mgo.v2"
)

type ctxKey struct{}

// GetFromCtx -
func GetFromCtx(ctx context.Context) (*mgo.Database, error) {
	db, ok := ctx.Value(ctxKey{}).(*mgo.Database)
	if !ok {
		return nil, errNotInCtx
	}
	return db, nil
}

// SetInCtx -
func SetInCtx(ctx context.Context, db *mgo.Database) context.Context {
	return context.WithValue(ctx, ctxKey{}, db)
}
