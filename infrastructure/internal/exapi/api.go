package exapi

import "context"

type MetaWeatherClient struct {
}

func (c *MetaWeatherClient) GetSample(ctx context.Context) (string, error) {
	return "", nil
}
