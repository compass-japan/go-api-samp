package exapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-api-samp/model/dto"
	"go-api-samp/util/config"
	"go-api-samp/util/log"
	"io"
	"io/ioutil"
	"net/http"
)

type MetaWeatherClient struct {
	Config *config.ExAPIConfig
}

func (c *MetaWeatherClient) GetSample(ctx context.Context) (*dto.ExApiResponse, error) {
	logger := log.GetLogger()

	ctx, cancel := context.WithTimeout(ctx, c.Config.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.Config.URL, nil)
	if err != nil {
		logger.Error(ctx, "failed to make request", err)
		return nil, err
	}

	httpClient := http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		logger.Error(ctx, "failed to request", err)
		return nil, err
	}

	defer func() {
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		m := fmt.Sprintf("error status code %d", res.StatusCode)
		logger.Error(ctx, m)
		return nil, errors.New(m)
	}

	var w dto.ExApiResponse
	if err := json.NewDecoder(res.Body).Decode(&w); err != nil {
		logger.Error(ctx, "json decode error.", err)
		return nil, err
	}

	return &w, nil
}
