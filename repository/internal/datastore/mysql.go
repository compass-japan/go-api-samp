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

	sql := "INSERT INTO WEATHER(dat, regWeather, location_id, comment) VALUES(?, ?, ?, ?)"
	stmt, err := c.Db.Prepare(sql)
	if err != nil {
		logger.Error(ctx, "failed to prepare statement.", err)
		return errors.DataStoreSystemError(err)
	}
	defer stmt.Close()

	_, err = stmt.Query(date, weather, locationId, comment)
	if err != nil {
		logger.Error(ctx, "failed to execute add regWeather query.", err)
		return errors.DataStoreSystemError(err)
	}

	return nil
}

func (c *MySQLClient) GetWeather(ctx context.Context, locationId int, date string) (*entity.Weather, error) {
	logger := log.GetLogger()

	sql := "SELECT dat, weather, location_id, city, comment FROM WEATHER as w INNER JOIN LOCATION as l ON w.location_id = l.id WHERE l.id = ? AND dat = ?"
	stmt, err := c.Db.Prepare(sql)
	if err != nil {
		logger.Error(ctx, "failed to prepare statement.", err)
		return nil, errors.DataStoreSystemError(err)
	}
	defer stmt.Close()

	row, err := stmt.Query(locationId, date)
	if err != nil {
		logger.Error(ctx, "failed to execute get location query.", err)
		return nil, errors.DataStoreSystemError(err)
	}

	if !row.Next() {
		m := "regWeather not found"
		logger.Info(ctx, m)
		return nil, errors.DataStoreValueNotFoundSystemError(err)
	}

	w := &entity.Weather{
		Location: &entity.Location{},
	}
	if err := row.Scan(&w.Dat, &w.Weather, &w.Location.Id, &w.Location.City, &w.Comment); err != nil {
		logger.Error(ctx, "failed to scan.", err)
		return nil, errors.DataStoreSystemError(err)
	}

	return w, nil
}

func (c *MySQLClient) FindLocation(ctx context.Context, locationId int) error {
	logger := log.GetLogger()

	sql := `SELECT count(*) FROM "LOCATION" WHERE id = ?`
	stmt, err := c.Db.Prepare(sql)
	if err != nil {
		logger.Error(ctx, "failed to prepare statement.", err)
		return errors.DataStoreSystemError(err)
	}
	defer stmt.Close()

	row, err := stmt.Query(locationId)
	if err != nil {
		logger.Error(ctx, "failed to execute get location query.", err)
		return errors.DataStoreSystemError(err)
	}

	row.Next()

	var count int
	if err := row.Scan(&count); err != nil {
		logger.Error(ctx, "failed to scan.", err)
		return errors.DataStoreSystemError(err)
	}

	if count == 0 {
		m := "invalid location id"
		logger.Info(ctx, m)
		return errors.DataStoreValueNotFoundSystemError(err)
	}

	return nil
}
