package api

import (
	"context"
	errors2 "errors"
	"github.com/stretchr/testify/assert"
	"go-api-samp/model/dto"
	"go-api-samp/model/entity"
	"go-api-samp/model/errors"
	"testing"
)

/*
 * Infrastructure,Repositoryをモック化してテスト
 */

var (
	e = errors2.New("")
)

func TestService(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "Test Register",
			test: testRegister,
		},
		{
			name: "Test GetWeather",
			test: testGetWeather,
		},
		{
			name: "Test GetAPIData",
			test: testGetAPIData,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, test.test)
	}
}

func testRegister(t *testing.T) {
	payload := &dto.RegisterRequest{}
	tests := []struct {
		name    string
		payload *dto.RegisterRequest
		mock    *mockStore
		err     error
	}{
		{
			name: "正常系 add",
			mock: &mockStore{
				addErr:    nil,
				updateErr: nil,
				weather:   &entity.Weather{},
				getErr:    nil,
				findErr:   nil,
			},
			err: nil,
		},
		{
			name: "正常系 update",
			mock: &mockStore{
				addErr:    nil,
				updateErr: nil,
				weather:   nil,
				getErr:    errors.DataStoreValueNotFoundSystemError(e),
				findErr:   nil,
			},
			err: nil,
		},
		{
			name: "find error",
			mock: &mockStore{
				findErr: e,
			},
			err: e,
		},
		{
			name: "get weather error",
			mock: &mockStore{
				findErr: nil,
				weather: &entity.Weather{},
				getErr:  errors.DataStoreSystemError(e),
			},
			err: e,
		},
		{
			name: "update error",
			mock: &mockStore{
				findErr:   nil,
				updateErr: e,
				weather:   nil,
				getErr:    nil,
			},
			err: e,
		},
		{
			name: "add error",
			mock: &mockStore{
				findErr: nil,
				getErr:  errors.DataStoreValueNotFoundSystemError(e),
				addErr:  e,
			},
			err: e,
		},
	}

	for _, test := range tests {
		tp := test
		t.Run(tp.name, func(t *testing.T) {
			s := &API{
				Store: tp.mock,
			}
			err := s.Register(context.Background(), payload)
			if tp.err == nil {
				assert.NoError(t, err)
			}
			if tp.err != nil {
				assert.Error(t, err)
			}
		})
	}
}

func testGetWeather(t *testing.T) {
	payload := &dto.GetWeatherRequest{}
	tests := []struct {
		name     string
		mock     *mockStore
		response *dto.GetWeatherResponse
		err      error
	}{
		{
			name: "正常系",
			mock: &mockStore{
				weather: &entity.Weather{
					Dat:     "20200101",
					Weather: 1,
					Location: &entity.Location{
						Id:   1,
						City: "新宿",
					},
					Comment: "test comment",
				},
				getErr: nil,
			},
			response: &dto.GetWeatherResponse{
				Location: "新宿",
				Date:     "20200101",
				Weather:  "Sunny",
				Comment:  "test comment",
			},
			err: nil,
		},
		{
			name: "error",
			mock: &mockStore{
				weather: nil,
				getErr:  e,
			},
			err: e,
		},
	}

	for _, test := range tests {
		tp := test
		s := &API{
			Store: tp.mock,
		}
		ety, err := s.GetWeather(context.Background(), payload)
		if tp.err == nil {
			assert.NoError(t, err)
			assert.Equal(t, tp.response, ety)
		}
		if tp.err != nil {
			assert.Error(t, err)
		}
	}
}

func testGetAPIData(t *testing.T) {
	tests := []struct {
		name string
		mock *mockInfra
	}{
		{
			name: "正常系",
			mock: &mockInfra{
				res: &dto.ExApiResponse{},
				err: nil,
			},
		},
		{
			name: "異常系",
			mock: &mockInfra{
				res: nil,
				err: e,
			},
		},
	}

	for _, test := range tests {
		tp := test
		s := &API{
			Infra: tp.mock,
		}
		w, err := s.GetAPIData(context.Background())
		if tp.mock.err == nil {
			assert.NoError(t, err)
			assert.Equal(t, tp.mock.res, w)
		}
		if tp.mock.err != nil {
			var res *dto.ExApiResponse
			assert.Equal(t, res, w)
		}
	}
}

type mockStore struct {
	addErr    error
	updateErr error
	weather   *entity.Weather
	getErr    error
	findErr   error
}

func (m *mockStore) AddWeather(ctx context.Context, locationId, weather int, date, comment string) error {
	return m.addErr
}

func (m *mockStore) UpdateWeather(ctx context.Context, locationId, weather int, date, comment string) error {
	return m.updateErr
}

func (m *mockStore) GetWeather(ctx context.Context, locationId int, date string) (*entity.Weather, error) {
	return m.weather, m.getErr
}

func (m *mockStore) FindLocation(ctx context.Context, locationId int) error {
	return m.findErr
}

type mockInfra struct {
	res *dto.ExApiResponse
	err error
}

func (m *mockInfra) GetExWeather(ctx context.Context) (*dto.ExApiResponse, error) {
	return m.res, m.err
}
