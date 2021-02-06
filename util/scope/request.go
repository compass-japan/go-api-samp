package scope

import (
	"context"
)

/*
 * request_idを扱う(X-Request-ID)
 * リクエストスコープにidを持たせることでアクセスのログを追える
 */

const (
	RequestIDContextKey string = "rid-ctx-key"
)

func SetRequestID(parent context.Context, id string) context.Context {
	return context.WithValue(parent, RequestIDContextKey, id)
}

func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	v := ctx.Value(RequestIDContextKey)
	id, ok := v.(string)
	if !ok {
		return ""
	}
	return id
}
