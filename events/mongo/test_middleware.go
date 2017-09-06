package mongo

import (
	"context"

	"github.com/hellomd/go-sdk/events"
	"github.com/hellomd/go-sdk/mongo"
	mgo "gopkg.in/mgo.v2"
)

// NewTestMiddleware -
func NewTestMiddleware(db *mgo.Database) func(context.Context, *events.Event, events.HandlerFunc) {
	return func(ctx context.Context, event *events.Event, next events.HandlerFunc) {
		next(mongo.SetInCtx(ctx, db), event)
	}
}
