package authentication

import (
	"net/http"
)

// NewMiddleware -
func NewMiddleware(secret []byte) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ctxAuth := NewContextAuthenticator(secret)

	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		ctx, err := ctxAuth(r.Context(), r.Header.Get(HeaderKey))
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		next(w, r.WithContext(ctx))
	}
}
