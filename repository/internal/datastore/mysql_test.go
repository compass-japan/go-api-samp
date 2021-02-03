package datastore

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go-api-samp/model/entity"
	"go-api-samp/model/errors"
	"reflect"
	"regexp"
	"testing"
)

func TestMySQLSuccess(t *testing.T) {
	tests := []struct {
		name string
		test func() func(t *testing.T)
	}{
		{
			name: "AddWeather Test",
			test: testAddWeather,
		},
		{
			name: "UpdateWeather Test",
			test: testUpdateWeather,
		},
		{
			name: "GetWeather Test",
			test: testGetWeather,
		},
		{
			name: "FindLocation Test",
			test: testGetLocation,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, test.test())
	}
}

func testAddWeather() func(t *testing.T) {
	return func(t *testing.T) {
		tests := []struct {
			name            string
			regWeather      *entity.Weather
			mockCreaterFunc func(mock sqlmock.Sqlmock, w *entity.Weather)
			isErr           bool
		}{
			{
				name: "正常系",
				regWeather: &entity.Weather{
					Dat:     "20200101",
					Weather: 1,
					Location: &entity.Location{
						Id: 1,
					},
					Comment: "comment",
				},
				mockCreaterFunc: func(mock sqlmock.Sqlmock, w *entity.Weather) {
					mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO WEATHER")).
						ExpectQuery().WithArgs(w.Dat, w.Weather, w.Location.Id, w.Comment).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				},
				isErr: false,
			},
			{
				name: "prepare error",
				regWeather: &entity.Weather{
					Dat:     "20200101",
					Weather: 1,
					Location: &entity.Location{
						Id: 1,
					},
					Comment: "comment",
				},
				mockCreaterFunc: func(mock sqlmock.Sqlmock, w *entity.Weather) {
					mock.ExpectPrepare(regexp.QuoteMeta("UPDATE WEATHER")).
						ExpectQuery().WithArgs(w.Dat, w.Weather, w.Location.Id, w.Comment).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				},
				isErr: true,
			},
			{
				name: "execute error",
				regWeather: &entity.Weather{
					Dat:     "20200101",
					Weather: 1,
					Location: &entity.Location{
						Id: 1,
					},
					Comment: "comment",
				},
				mockCreaterFunc: func(mock sqlmock.Sqlmock, w *entity.Weather) {
					mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO WEATHER")).
						ExpectQuery().WithArgs(w.Dat, w.Weather, w.Location.Id).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				},
				isErr: true,
			},
		}

		for _, test := range tests {
			tp := test
			t.Run(tp.name, func(t *testing.T) {
				db, mock, _ := sqlmock.New()
				c := &MySQLClient{Db: db}
				tp.mockCreaterFunc(mock, tp.regWeather)
				err := c.AddWeather(context.Background(), tp.regWeather.Location.Id, tp.regWeather.Weather, tp.regWeather.Dat, tp.regWeather.Comment)
				if tp.isErr {
					if assert.Error(t, err) {
						e, ok := err.(errors.SystemError)
						re := errors.DataStoreSystemError(err)
						assert.True(t, ok)
						assert.Equal(t, re.Message(), e.Message())
					}
				} else {
					assert.NoError(t, err)
				}
			})
		}
	}
}

func testUpdateWeather() func(t *testing.T) {
	return func(t *testing.T) {
		tests := []struct {
			name            string
			regWeather      *entity.Weather
			mockCreaterFunc func(mock sqlmock.Sqlmock, w *entity.Weather)
			isErr           bool
		}{
			{
				name: "正常系",
				regWeather: &entity.Weather{
					Dat:     "20200101",
					Weather: 0,
					Location: &entity.Location{
						Id: 1,
					},
					Comment: "comment",
				},
				mockCreaterFunc: func(mock sqlmock.Sqlmock, w *entity.Weather) {
					mock.ExpectPrepare(regexp.QuoteMeta("UPDATE WEATHER")).
						ExpectQuery().WithArgs(w.Weather, w.Comment, w.Location.Id, w.Dat).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				},
				isErr: false,
			},
			{
				name: "prepare error",
				regWeather: &entity.Weather{
					Dat:     "20200101",
					Weather: 0,
					Location: &entity.Location{
						Id: 1,
					},
					Comment: "comment",
				},
				mockCreaterFunc: func(mock sqlmock.Sqlmock, w *entity.Weather) {
					mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO WEATHER")).
						ExpectQuery().WithArgs(w.Weather, w.Comment, w.Location.Id, w.Dat).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				},
				isErr: true,
			},
			{
				name: "execute error",
				regWeather: &entity.Weather{
					Dat:     "20200101",
					Weather: 0,
					Location: &entity.Location{
						Id: 1,
					},
					Comment: "comment",
				},
				mockCreaterFunc: func(mock sqlmock.Sqlmock, w *entity.Weather) {
					mock.ExpectPrepare(regexp.QuoteMeta("UPDATE WEATHER")).
						ExpectQuery().WithArgs(w.Weather, w.Comment).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				},
				isErr: true,
			},
		}

		for _, test := range tests {
			tp := test
			t.Run(tp.name, func(t *testing.T) {
				db, mock, _ := sqlmock.New()
				c := &MySQLClient{Db: db}
				tp.mockCreaterFunc(mock, tp.regWeather)
				err := c.UpdateWeather(context.Background(), tp.regWeather.Location.Id, tp.regWeather.Weather, tp.regWeather.Dat, tp.regWeather.Comment)
				if tp.isErr {
					if assert.Error(t, err) {
						e, ok := err.(errors.SystemError)
						re := errors.DataStoreSystemError(err)
						assert.True(t, ok)
						assert.Equal(t, re.Message(), e.Message())
					}
				} else {
					assert.NoError(t, err)
				}
			})
		}
	}
}

func testGetWeather() func(t *testing.T) {
	return func(t *testing.T) {
		tests := []struct {
			name       string
			locationId int
			dat        string
			mock       func(mock sqlmock.Sqlmock, locationId int, dat string)
			retWeather *entity.Weather
			isErr      bool
			systemErr  func(cause error) errors.SystemError
		}{
			{
				name:       "正常系",
				locationId: 1,
				dat:        "20200101",
				mock: func(mock sqlmock.Sqlmock, locationId int, dat string) {
					mock.ExpectPrepare(regexp.QuoteMeta("SELECT dat, weather, location_id, city, comment FROM WEATHER")).
						ExpectQuery().WithArgs(locationId, dat).
						WillReturnRows(sqlmock.NewRows([]string{"dat", "weather", "location_id", "city", "comment"}).
							AddRow("20200101", "1", "1", "新宿", "comment"))
				},
				retWeather: &entity.Weather{
					Dat:     "20200101",
					Weather: 1,
					Location: &entity.Location{
						Id:   1,
						City: "新宿",
					},
					Comment: "comment",
				},
				isErr: false,
			},
			{
				name:       "prepare error",
				locationId: 1,
				dat:        "20200101",
				mock: func(mock sqlmock.Sqlmock, locationId int, dat string) {
					mock.ExpectPrepare(regexp.QuoteMeta("SELECT dat FROM WEATHER")).
						ExpectQuery().WithArgs(locationId, dat).
						WillReturnRows(sqlmock.NewRows([]string{"dat", "weather", "location_id", "city", "comment"}).
							AddRow("20200101", "0", "1", "新宿", "comment"))
				},
				isErr:     true,
				systemErr: errors.DataStoreSystemError,
			},
			{
				name:       "execute error",
				locationId: 1,
				dat:        "20200101",
				mock: func(mock sqlmock.Sqlmock, locationId int, dat string) {
					mock.ExpectPrepare(regexp.QuoteMeta("SELECT dat, weather, location_id, city, comment FROM WEATHER")).
						ExpectQuery().WithArgs(0).
						WillReturnRows(sqlmock.NewRows([]string{"dat", "weather", "location_id", "city", "comment"}).
							AddRow("20200101", "0", "1", "新宿", "comment"))
				},
				isErr:     true,
				systemErr: errors.DataStoreSystemError,
			},
			{
				name:       "not found error",
				locationId: 1,
				dat:        "20200101",
				mock: func(mock sqlmock.Sqlmock, locationId int, dat string) {
					mock.ExpectPrepare(regexp.QuoteMeta("SELECT dat, weather, location_id, city, comment FROM WEATHER")).
						ExpectQuery().WithArgs(locationId, dat).
						WillReturnRows(sqlmock.NewRows([]string{"dat", "weather", "location_id", "city", "comment"}))
				},
				isErr:     true,
				systemErr: errors.DataStoreValueNotFoundSystemError,
			},
			{
				name:       "scan error",
				locationId: 1,
				dat:        "20200101",
				mock: func(mock sqlmock.Sqlmock, locationId int, dat string) {
					mock.ExpectPrepare(regexp.QuoteMeta("SELECT dat, weather, location_id, comment FROM WEATHER")).
						ExpectQuery().WithArgs(locationId, dat).
						WillReturnRows(sqlmock.NewRows([]string{"dat", "weather", "location_id"}).
							AddRow("20200101", "0", "1"))
				},
				isErr:     true,
				systemErr: errors.DataStoreSystemError,
			},
		}

		for _, test := range tests {
			tp := test
			t.Run(tp.name, func(t *testing.T) {
				db, mock, _ := sqlmock.New()
				c := &MySQLClient{Db: db}
				tp.mock(mock, tp.locationId, tp.dat)
				result, err := c.GetWeather(context.Background(), tp.locationId, tp.dat)
				if tp.isErr {
					if assert.Error(t, err) {
						e, ok := err.(errors.SystemError)
						re := tp.systemErr(err)
						assert.True(t, ok)
						assert.Equal(t, re.Message(), e.Message())
					}
				} else {
					if assert.NoError(t, err) {
						assert.True(t, reflect.DeepEqual(result, tp.retWeather))
					}
				}
			})
		}
	}
}

func testGetLocation() func(t *testing.T) {
	return func(t *testing.T) {
		tests := []struct {
			name       string
			locationId int
			value      int
			mock       func(mock sqlmock.Sqlmock, locationId, value int)
			isErr      bool
			systemErr  func(cause error) errors.SystemError
		}{
			{
				name:       "正常系",
				locationId: 1,
				value:      1,
				mock: func(mock sqlmock.Sqlmock, locationId, value int) {
					mock.ExpectPrepare(regexp.QuoteMeta("SELECT count(*)")).
						ExpectQuery().WithArgs(locationId).
						WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(value))
				},
				isErr: false,
			},
			{
				name:       "prepare error",
				locationId: 1,
				value:      1,
				mock: func(mock sqlmock.Sqlmock, locationId, value int) {
					mock.ExpectPrepare(regexp.QuoteMeta("SELEC count(*)")).
						ExpectQuery().WithArgs(locationId).
						WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(value))
				},
				isErr:     true,
				systemErr: errors.DataStoreSystemError,
			},
			{
				name:       "execute error",
				locationId: 1,
				value:      1,
				mock: func(mock sqlmock.Sqlmock, locationId, value int) {
					mock.ExpectPrepare(regexp.QuoteMeta("SELECT count(*)")).
						ExpectQuery().WithArgs(0).
						WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(value))
				},
				isErr:     true,
				systemErr: errors.DataStoreSystemError,
			},
			{
				name:       "scan convert error",
				locationId: 1,
				value:      1,
				mock: func(mock sqlmock.Sqlmock, locationId, value int) {
					mock.ExpectPrepare(regexp.QuoteMeta("SELECT count(*)")).
						ExpectQuery().WithArgs(locationId).
						WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("aaa"))
				},
				isErr:     true,
				systemErr: errors.DataStoreSystemError,
			},
			{
				name:       "value == 0",
				locationId: 1,
				value:      1,
				mock: func(mock sqlmock.Sqlmock, locationId, value int) {
					mock.ExpectPrepare(regexp.QuoteMeta("SELECT count(*)")).
						ExpectQuery().WithArgs(locationId).
						WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
				},
				isErr:     true,
				systemErr: errors.DataStoreValueNotFoundSystemError,
			},
		}

		for _, test := range tests {
			tp := test
			t.Run(tp.name, func(t *testing.T) {
				db, mock, _ := sqlmock.New()
				c := &MySQLClient{Db: db}
				tp.mock(mock, tp.locationId, tp.value)
				err := c.FindLocation(context.Background(), tp.locationId)
				if tp.isErr {
					assert.Error(t, err)
					e, ok := err.(errors.SystemError)
					re := tp.systemErr(err)
					assert.True(t, ok)
					assert.Equal(t, re.Message(), e.Message())
				} else {
					assert.NoError(t, err)
				}
			})
		}
	}
}
