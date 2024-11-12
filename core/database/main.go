package database

import "github.com/tittuvarghese/core/storage"

type RelationalDatabase struct {
	Instance *storage.RelationalDB
}

func NewRelationalDatabase(conn string) (*RelationalDatabase, error) {
	handler, err := storage.NewRelationalDbHandler(conn)
	if err != nil {
		return nil, err
	}
	return &RelationalDatabase{Instance: handler}, nil
}
