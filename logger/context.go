package logger

import (
	"context"
	"errors"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

var errNotInCtx = errors.New("no logrus in context")
var noOpLogger *logrus.Entry

type ctxKey struct{}

func init() {
	l := logrus.New()
	l.Out = ioutil.Discard
	noOpLogger = l.WithFields(logrus.Fields{})
}

// GetFromCtx get the logrus instance from the context
func GetFromCtx(ctx context.Context) (*logrus.Entry, error) {
	logger, ok := ctx.Value(ctxKey{}).(*logrus.Entry)
	if !ok {
		return nil, errNotInCtx
	}
	return logger, nil
}

// SetInCtx sets given logrus instance to given context
func SetInCtx(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

func Get(ctx context.Context) *logrus.Entry {
	logger, ok := ctx.Value(ctxKey{}).(*logrus.Entry)
	if !ok {
		return noOpLogger
	}

	return logger
}
