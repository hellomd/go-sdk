package events

import (
	"context"

	"github.com/hellomd/go-sdk/authentication"
	"github.com/hellomd/go-sdk/requestid"
)

func DefaultHeaders(ctx context.Context) (map[string]string, error) {
	auth := authentication.GetServiceTokenFromCtx(ctx)
	reqID, err := requestid.GetRequestIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		authentication.HeaderKey:     authentication.Scheme + " " + auth,
		requestid.RequestIDHeaderKey: reqID,
	}, nil
}
