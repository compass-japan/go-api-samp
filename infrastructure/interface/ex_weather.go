package infrastructure

import "context"

type (
	MetaWeatherManager interface {
		GetSample(ctx context.Context) (string, error)
	}
)
