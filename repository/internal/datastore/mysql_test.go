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

func TestSample(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	c := &MySQLClient{
		Db: db,
	}

	mock.ExpectPrepare(regexp.QuoteMeta("SELECT count(*)")).
		ExpectQuery().WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	b, err := c.GetLocation(context.Background(), 1)
	assert.NoError(t, err)
	assert.True(t, b, mock)
}

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
			name: "GetWeather Test",
			test: testGetWeather,
		},
		{
			name: "GetLocation Test",
			test: testGetLocation,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, test.test())
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
			systemErr  errors.SystemErrorBuilder
			result     bool
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
				isErr:  false,
				result: true,
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
				systemErr: errors.System.DataStoreError,
				result:    false,
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
				systemErr: errors.System.DataStoreError,
				result:    false,
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
				systemErr: errors.System.DataStoreError,
				result:    false,
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
				systemErr: errors.System.DataStoreValueNotFoundError,
				result:    false,
			},
		}

		for _, test := range tests {
			tp := test
			t.Run(tp.name, func(t *testing.T) {
				db, mock, _ := sqlmock.New()
				c := &MySQLClient{Db: db}
				tp.mock(mock, tp.locationId, tp.value)
				b, err := c.GetLocation(context.Background(), tp.locationId)
				if tp.isErr {
					assert.Error(t, err)
					e, ok := err.(errors.SystemError)
					re := tp.systemErr(err)
					assert.True(t, ok)
					assert.Equal(t, re.Message(), e.Message())
				} else {
					assert.NoError(t, err)
				}
				assert.Equal(t, tp.result, b)
			})
		}
	}
}

func testAddWeather() func(t *testing.T) {
	return func(t *testing.T) {
		tests := []struct {
			name    string
			weather *entity.Weather
			mock    func(mock sqlmock.Sqlmock, w *entity.Weather)
			isErr   bool
		}{
			{
				name: "正常系",
				weather: &entity.Weather{
					Dat:        "20200101",
					Weather:    0,
					LocationId: 1,
					Comment:    "comment",
				},
				mock: func(mock sqlmock.Sqlmock, w *entity.Weather) {
					mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO WEATHER")).
						ExpectQuery().WithArgs(w.Dat, w.Weather, w.LocationId, w.Comment).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				},
				isErr: false,
			},
			{
				name: "prepare error",
				weather: &entity.Weather{
					Dat:        "20200101",
					Weather:    0,
					LocationId: 1,
					Comment:    "comment",
				},
				mock: func(mock sqlmock.Sqlmock, w *entity.Weather) {
					mock.ExpectPrepare(regexp.QuoteMeta("UPDATE WEATHER")).
						ExpectQuery().WithArgs(w.Dat, w.Weather, w.LocationId, w.Comment).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				},
				isErr: true,
			},
			{
				name: "execute error",
				weather: &entity.Weather{
					Dat:        "20200101",
					Weather:    0,
					LocationId: 1,
					Comment:    "comment",
				},
				mock: func(mock sqlmock.Sqlmock, w *entity.Weather) {
					mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO WEATHER")).
						ExpectQuery().WithArgs(w.Dat, w.Weather, w.LocationId).
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
				tp.mock(mock, tp.weather)
				err := c.AddWeather(context.Background(), tp.weather.LocationId, tp.weather.Weather, tp.weather.Dat, tp.weather.Comment)
				if tp.isErr {
					if assert.Error(t, err) {
						e, ok := err.(errors.SystemError)
						re := errors.System.DataStoreError(err)
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
			weather    *entity.Weather
			isErr      bool
			systemErr  errors.SystemErrorBuilder
		}{
			{
				name:       "正常系",
				locationId: 1,
				dat:        "20200101",
				mock: func(mock sqlmock.Sqlmock, locationId int, dat string) {
					mock.ExpectPrepare(regexp.QuoteMeta("SELECT dat, weather, location_id, comment FROM WEATHER")).
						ExpectQuery().WithArgs(locationId, dat).
						WillReturnRows(sqlmock.NewRows([]string{"dat", "weather", "location_id", "comment"}).
							AddRow("20200101", "0", "1", "comment"))
				},
				weather: &entity.Weather{
					Dat:        "20200101",
					Weather:    0,
					LocationId: 1,
					Comment:    "comment",
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
						WillReturnRows(sqlmock.NewRows([]string{"dat", "weather", "location_id", "comment"}).
							AddRow("20200101", "0", "1", "comment"))
				},
				isErr:     true,
				systemErr: errors.System.DataStoreError,
			},
			{
				name:       "execute error",
				locationId: 1,
				dat:        "20200101",
				mock: func(mock sqlmock.Sqlmock, locationId int, dat string) {
					mock.ExpectPrepare(regexp.QuoteMeta("SELECT dat, weather, location_id, comment FROM WEATHER")).
						ExpectQuery().WithArgs(0).
						WillReturnRows(sqlmock.NewRows([]string{"dat", "weather", "location_id", "comment"}).
							AddRow("20200101", "0", "1", "comment"))
				},
				isErr:     true,
				systemErr: errors.System.DataStoreError,
			},
			{
				name:       "not found error",
				locationId: 1,
				dat:        "20200101",
				mock: func(mock sqlmock.Sqlmock, locationId int, dat string) {
					mock.ExpectPrepare(regexp.QuoteMeta("SELECT dat, weather, location_id, comment FROM WEATHER")).
						ExpectQuery().WithArgs(locationId, dat).
						WillReturnRows(sqlmock.NewRows([]string{"dat", "weather", "location_id", "comment"}))
				},
				isErr:     true,
				systemErr: errors.System.DataStoreValueNotFoundError,
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
				systemErr: errors.System.DataStoreError,
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
						assert.True(t, reflect.DeepEqual(result, tp.weather))
					}
				}
			})
		}
	}
}