package datastore

import (
	"database/sql"
	"go-api-samp/model/entity"
	"go-api-samp/model/errors"
	"go-api-samp/util/log"
)

/*
 * 初期化処理のRepository(datastore)層の実装
 */

type InitClient struct {
	Db *sql.DB
}

func (c *InitClient) FindAllLocations() ([]entity.Location, error) {
	logger := log.GetLogger()

	sql := "SELECT id, city FROM LOCATION"
	stmt, err := c.Db.Prepare(sql)
	if err != nil {
		logger.Error(nil, "failed to prepare statement.", err)
		return nil, errors.DataStoreSystemError(err)
	}
	defer stmt.Close()

	row, err := stmt.Query()
	if err != nil {
		logger.Error(nil, "failed to execute get all locations query.", err)
		return nil, errors.DataStoreSystemError(err)
	}

	locations := make([]entity.Location, 0)

	for row.Next() {
		var location entity.Location
		if err := row.Scan(&location.Id, &location.City); err != nil {
			logger.Error(nil, "failed to scan.", err)
			return nil, errors.DataStoreSystemError(err)
		}
		locations = append(locations, location)
	}

	if len(locations) == 0 {
		logger.Error(nil, "no location.", err)
		return nil, errors.DataStoreValueNotFoundSystemError(nil)
	}

	return locations, nil
}
