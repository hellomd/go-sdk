package recovery

import (
	"github.com/getsentry/raven-go"
	"github.com/sirupsen/logrus"
)

// Reporter reports errors to multiple destinations
type Reporter struct {
	raven   *raven.Client
	logger  *logrus.Logger
	httpCtx *raven.Http
}

// Error logs error on logger and to sentry if its available
func (r *Reporter) Error(err error) {
	if r.raven != nil {
		r.raven.CaptureError(err, nil)
	}
	logrus.Error(err)
}
