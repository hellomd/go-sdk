package recovery

import (
	"context"
	"errors"
	"testing"

	"github.com/hellomd/go-sdk/events"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRecoveryMiddleware(t *testing.T) {
	Convey("RecoveryMiddleware", t, func() {
		var (
			receivedErr error
			receivedCtx context.Context
		)

		handleError := func(ctx context.Context, err error) {
			receivedErr = err
			receivedCtx = ctx
		}

		ctx := context.Background()
		fakeAck := new(fakeAcknowledger)
		e := &events.Event{
			Acknowledger: fakeAck,
			Key:          "foo.bar",
		}
		m := NewMiddleware("")
		m.handleError = handleError
		expectedErr := errors.New("test error")
		m.Process(ctx, e, func(ctx context.Context, e *events.Event) {
			panic(expectedErr)
		})
		Convey("it recovers on panic", func() {
			So(receivedErr, ShouldNotBeNil)
			So(receivedErr.Error(), ShouldEqual, "panic recovered: "+expectedErr.Error())
		})

		Convey("event should be rejected with error", func() {
			So(fakeAck.rejected, ShouldBeTrue)
			So(fakeAck.requeued, ShouldBeFalse)
			So(fakeAck.err.Error(), ShouldEqual, "panic recovered: "+expectedErr.Error())
		})
	})
}

type fakeAcknowledger struct {
	acked, rejected, requeued bool
	err                       error
}

func (a *fakeAcknowledger) Reject(requeue bool, err error) {
	a.rejected = true
	a.requeued = requeue
	a.err = err
}

func (a *fakeAcknowledger) Ack() {
	a.acked = true
}
