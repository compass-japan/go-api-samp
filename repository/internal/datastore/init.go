package datastore

import (
	"database/sql"
	"errors"
	"go-api-samp/model/entity"
)

type InitClient struct {
	Db *sql.DB
}

func (c *InitClient) FindAllLocations() (*[]entity.Location, error) {
	return nil, errors.New("")
}
