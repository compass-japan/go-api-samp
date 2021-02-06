package exapi

import (
	"context"
	"encoding/json"
	"fmt"
	"go-api-samp/model/dto"
	"go-api-samp/model/errors"
	"go-api-samp/util/config"
	"go-api-samp/util/log"
	"io"
	"io/ioutil"
	"net/http"
)

/*
 * 外部APIへのアクセス実装
 * https://golang.org/pkg/net/http/
 * response bodyは読み切ることが必要。読み切らないとcloseされないかもしれない
 */

type MetaWeatherClient struct {
	Config *config.ExAPIConfig
}

func (c *MetaWeatherClient) GetExWeather(ctx context.Context) (*dto.ExApiResponse, error) {
	logger := log.GetLogger()

	ctx, cancel := context.WithTimeout(ctx, c.Config.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.Config.URL, nil)
	if err != nil {
		logger.Error(ctx, "failed to make request", err)
		return nil, errors.ExAPISystemError(err)
	}

	httpClient := http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		logger.Error(ctx, "failed to request", err)
		return nil, errors.ExAPISystemError(err)
	}

	defer func() {
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		m := fmt.Sprintf("error status code %d", res.StatusCode)
		logger.Error(ctx, m)
		return nil, errors.ExAPISystemError(err)
	}

	var w dto.ExApiResponse
	if err := json.NewDecoder(res.Body).Decode(&w); err != nil {
		logger.Error(ctx, "json decode error. response body: %s", res.Body, err)
		return nil, errors.ExAPISystemError(err)
	}

	return &w, nil
}
