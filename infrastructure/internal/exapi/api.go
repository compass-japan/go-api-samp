package exapi

import (
	"context"
	"errors"
	"fmt"
	"go-api-samp/util/config"
	"go-api-samp/util/log"
	"io"
	"io/ioutil"
	"net/http"
)

type MetaWeatherClient struct {
	Config *config.ExAPIConfig
}

func (c *MetaWeatherClient) GetSample(ctx context.Context) (string, error) {
	logger := log.GetLogger()

	ctx, cancel := context.WithTimeout(ctx, c.Config.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.Config.URL, nil)
	if err != nil {
		logger.Error(ctx, "failed to make request", err)
		return "", err
	}

	httpClient := http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		logger.Error(ctx, "failed to request", err)
		return "", err
	}

	defer func() {
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		m := fmt.Sprintf("error status code %d", res.StatusCode)
		logger.Error(ctx, m)
		return "", errors.New(m)
	}

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error(ctx, "failed to read response", err)
		return "", err
	}

	return string(result), nil
}
