package datastore

import (
	errors2 "errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go-api-samp/model/entity"
	"go-api-samp/model/errors"
	"regexp"
	"testing"
)

var (
	e = errors2.New("test error")
)

func TestFindAllLocations(t *testing.T) {
	tests := []struct {
		name            string
		mockCreaterFunc func(mock sqlmock.Sqlmock)
		retLocations    []entity.Location
		isErr           bool
		sysErr          func(cause error) errors.SystemError
	}{
		{
			name: "正常系",
			mockCreaterFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(regexp.QuoteMeta("SELECT id, city FROM LOCATION")).
					ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"id", "city"}).
					AddRow(1, "新宿").AddRow(2, "中野"))
			},
			retLocations: []entity.Location{
				entity.Location{Id: 1, City: "新宿"},
				entity.Location{Id: 2, City: "中野"},
			},
			isErr: false,
		},
		{
			name: "prepare error",
			mockCreaterFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(regexp.QuoteMeta("SELECT id, city FROM LOCATION")).
					WillReturnError(e)
			},
			isErr:  true,
			sysErr: errors.DataStoreSystemError,
		},
		{
			name: "execute error",
			mockCreaterFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(regexp.QuoteMeta("SELECT id, city FROM LOCATION")).
					ExpectQuery().WillReturnError(e)
			},
			isErr:  true,
			sysErr: errors.DataStoreSystemError,
		},
		{
			name: "scan error",
			mockCreaterFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(regexp.QuoteMeta("SELECT id, city FROM LOCATION")).
					ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"id"}).
					AddRow(1).AddRow(2))
			},
			isErr:  true,
			sysErr: errors.DataStoreSystemError,
		},
		{
			name: "not found error",
			mockCreaterFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(regexp.QuoteMeta("SELECT id, city FROM LOCATION")).
					ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"id", "city"}))
			},
			isErr:  true,
			sysErr: errors.DataStoreValueNotFoundSystemError,
		},
	}

	for _, test := range tests {
		tp := test
		t.Run(tp.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			c := &InitClient{Db: db}
			tp.mockCreaterFunc(mock)
			locations, err := c.FindAllLocations()
			if tp.isErr {
				if assert.Error(t, err) {
					e, ok := err.(errors.SystemError)
					re := tp.sysErr(err)
					assert.True(t, ok)
					assert.Equal(t, re.Message(), e.Message())
				}
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, tp.retLocations, locations)
				}
			}
		})
	}
}
