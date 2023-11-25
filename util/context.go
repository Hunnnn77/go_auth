package util

import (
	"context"
	"net/http"
)

func HandleContext(r *http.Request, payload any) (*http.Request, any) {
	ctx := r.Context()
	senderCtx := context.WithValue(ctx, JWT, payload)
	return r.WithContext(senderCtx), ctx.Value(JWT)
}
