package datastore

import (
	"database/sql"
	"go-api-samp/model/entity"
)

type InitClient struct {
	Db *sql.DB
}

func (c *InitClient) FindAllLocations() (*[]entity.Location, error) {
	return nil, nil
}
