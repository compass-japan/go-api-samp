package datastore

import (
	"context"
	"database/sql"
	"go-api-samp/model/entity"
	"go-api-samp/model/errors"
	"go-api-samp/util/log"
)

type MySQLClient struct {
	Db *sql.DB
}

func (c *MySQLClient) AddWeather(ctx context.Context, locationId, weather int, date, comment string) error {
	logger := log.GetLogger()

	sql := "INSERT INTO WEATHER(dat, weather, location_id, comment) VALUES(?, ?, ?, ?)"
	stmt, err := c.Db.Prepare(sql)
	if err != nil {
		logger.Error(ctx, "failed to prepare statement.", err)
		return errors.System.DataStoreError(err)
	}
	defer stmt.Close()

	_, err = stmt.Query(date, weather, locationId, comment)
	if err != nil {
		logger.Error(ctx, "failed to execute add weather query.", err)
		return errors.System.DataStoreError(err)
	}

	return nil
}

func (c *MySQLClient) GetWeather(ctx context.Context, locationId int, date string) (*entity.Weather, error) {
	logger := log.GetLogger()

	sql := "SELECT dat, weather, location_id, comment FROM WEATHER WHERE location_id = ? AND dat = ?"
	stmt, err := c.Db.Prepare(sql)
	if err != nil {
		logger.Error(ctx, "failed to prepare statement.", err)
		return nil, errors.System.DataStoreError(err)
	}
	defer stmt.Close()

	row, err := stmt.Query(locationId, date)
	if err != nil {
		logger.Error(ctx, "failed to execute get location query.", err)
		return nil, errors.System.DataStoreError(err)
	}

	if !row.Next() {
		m := "weather not found"
		logger.Info(ctx, m)
		return nil, errors.System.DataStoreValueNotFoundError(err)
	}

	w := &entity.Weather{}
	if err := row.Scan(&w.Dat, &w.Weather, &w.LocationId, &w.Comment); err != nil {
		logger.Error(ctx, "failed to scan.", err)
		return nil, errors.System.DataStoreError(err)
	}

	return w, nil
}

func (c *MySQLClient) GetLocation(ctx context.Context, locationId int) (bool, error) {
	logger := log.GetLogger()

	sql := `SELECT count(*) FROM "LOCATION" WHERE id = ?`
	stmt, err := c.Db.Prepare(sql)
	if err != nil {
		logger.Error(ctx, "failed to prepare statement.", err)
		return false, errors.System.DataStoreError(err)
	}
	defer stmt.Close()

	row, err := stmt.Query(locationId)
	if err != nil {
		logger.Error(ctx, "failed to execute get location query.", err)
		return false, errors.System.DataStoreError(err)
	}

	row.Next()

	var count int
	if err := row.Scan(&count); err != nil {
		logger.Error(ctx, "failed to scan.", err)
		return false, errors.System.DataStoreError(err)
	}

	if count == 0 {
		m := "invalid location id"
		logger.Info(ctx, m)
		return false, errors.System.DataStoreValueNotFoundError(err)
	}

	return true, nil
}
