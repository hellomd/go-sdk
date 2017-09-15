package logger

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/hellomd/go-sdk/events"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

type testFormatter struct {
	logCalls int
	entry    *logrus.Entry
}

func (f *testFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	f.logCalls++
	f.entry = entry
	return []byte{}, nil
}

func ShouldLastWithin(actual interface{}, expected ...interface{}) string {
	const wrongArgs = "You must provide two durations as arguments for this assertion."
	if len(expected) != 2 {
		return wrongArgs
	}

	actualDuration, okActual := actual.(time.Duration)
	expectedDuration, okExpected := expected[0].(time.Duration)
	margin, okMargin := expected[1].(time.Duration)

	if !okActual || !okExpected || !okMargin {
		return wrongArgs
	}

	if actualDuration < expectedDuration-margin || actualDuration > expectedDuration+margin {
		return fmt.Sprintf("Expected '%v' to last '%v' give or take '%v' (but it didn't)!",
			actualDuration, expectedDuration, margin)
	}

	return ""
}

func TestMiddleware(t *testing.T) {
	Convey("logger middleware", t, func() {
		formatter := new(testFormatter)

		m := NewMiddleware("myapp", "someenv", func() *logrus.Logger {
			logger := logrus.New()
			logger.Formatter = formatter
			logger.Out = ioutil.Discard
			return logger
		}())

		ctx := context.Background()
		e := &events.Event{
			Acknowledger: new(fakeAcknowledger),
			Key:          "foo.bar",
		}

		Convey("any handling", func() {
			m(ctx, e, func(ctx context.Context, e *events.Event) {
				time.Sleep(15 * time.Millisecond)
			})

			Convey("should log with correct message and fields", func() {
				So(formatter.logCalls, ShouldEqual, 1)
				latency := formatter.entry.Data["took"].(time.Duration)
				So(latency, ShouldLastWithin, 15*time.Millisecond, 2*time.Millisecond)

				So(formatter.entry.Message, ShouldStartWith, fmt.Sprintf("evt foo.bar | %v | ", latency))
				So(formatter.entry.Data, ShouldResemble, logrus.Fields{
					"request_id":       nil,
					"event_key":        "foo.bar",
					"application_name": "myapp",
					"environment":      "someenv",
					"took":             latency,
					"status":           "",
				})
			})
		})

		Convey("acknowledge event", func() {
			m(ctx, e, func(ctx context.Context, e *events.Event) {
				e.Ack()
			})

			So(formatter.logCalls, ShouldEqual, 1)
			So(formatter.entry.Message, ShouldEndWith, "| "+statusAcked)
			So(formatter.entry.Data["status"], ShouldEqual, statusAcked)
		})

		Convey("reject event", func() {
			m(ctx, e, func(ctx context.Context, e *events.Event) {
				e.Reject(false, errors.New("some error"))
			})

			So(formatter.logCalls, ShouldEqual, 1)
			So(formatter.entry.Message, ShouldEndWith, "| "+statusRejected+` "some error"`)
			So(formatter.entry.Data["status"], ShouldEqual, statusRejected)
			So(formatter.entry.Data["error"], ShouldEqual, "some error")
		})

		Convey("reject and requeue event", func() {
			m(ctx, e, func(ctx context.Context, e *events.Event) {
				e.Reject(true, errors.New("some error"))
			})

			So(formatter.logCalls, ShouldEqual, 1)
			So(formatter.entry.Message, ShouldEndWith, "| "+statusRequeued+` "some error"`)
			So(formatter.entry.Data["status"], ShouldEqual, statusRequeued)
			So(formatter.entry.Data["error"], ShouldEqual, "some error")
		})
	})
}

func TestAcknowledger(t *testing.T) {
	Convey("logger acknowledger", t, func() {
		inner := new(fakeAcknowledger)
		a := &loggerAcknowledger{inner: inner}

		Convey("reject and requeue", func() {
			err := errors.New("nope")
			a.Reject(true, err)

			Convey("status should be requeued", func() {
				So(a.status, ShouldEqual, statusRequeued)
				So(a.err, ShouldEqual, err)
			})

			Convey("underlying acknowledger should be called", func() {
				So(inner.rejected, ShouldBeTrue)
				So(inner.requeued, ShouldBeTrue)
				So(inner.err, ShouldEqual, err)
			})
		})

		Convey("reject without requeuing", func() {
			err := errors.New("nope")
			a.Reject(false, err)

			Convey("status should be rejected", func() {
				So(a.status, ShouldEqual, statusRejected)
				So(a.err, ShouldEqual, err)
			})

			Convey("underlying acknowledger should be called", func() {
				So(inner.rejected, ShouldBeTrue)
				So(inner.requeued, ShouldBeFalse)
				So(inner.err, ShouldEqual, err)
			})
		})

		Convey("acknowledge", func() {
			a.Ack()

			Convey("status should be acked", func() {
				So(a.status, ShouldEqual, statusAcked)
			})

			Convey("underlying acknowledger should be called", func() {
				So(inner.acked, ShouldBeTrue)
			})
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
