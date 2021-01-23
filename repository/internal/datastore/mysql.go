package datastore

import "context"

type MysqlClient struct {
}

func (c *MysqlClient) Add(ctx context.Context, locationId int, date, weather, comment string) error {
	return nil
}

func (c *MysqlClient) Get(ctx context.Context, locationId int, date string) (string, error) {
	return "", nil
}
