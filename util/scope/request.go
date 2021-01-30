package scope

import (
	"context"
)

const (
	RequestIDContextKey string = "rid-ctx-key"
)

func SetRequestID(parent context.Context, id string) context.Context {
	return context.WithValue(parent, RequestIDContextKey, id)
}

func GetRequestID(ctx context.Context) string {
	v := ctx.Value(RequestIDContextKey)
	id, ok := v.(string)
	if !ok {
		return ""
	}
	return id
}
