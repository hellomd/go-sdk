package recovery

import (
	"context"
	"fmt"

	"github.com/hellomd/go-sdk/events"
	"github.com/sirupsen/logrus"
)

func NewMiddleware() events.PipelineHandlerFunc {
	return func(ctx context.Context, event *events.Event, next events.HandlerFunc) {
		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("panic recovered: %v", r)
				logrus.Error(err)
				event.Reject(false, err)
			}
		}()

		next(ctx, event)
	}
}
