package database

import "github.com/tittuvarghese/ss-go-core/storage"

type RelationalDatabase struct {
	Instance *storage.RelationalDB
}

var DbTxn storage.AtomicTransaction
var DbOps storage.Operation
var DbExpr storage.Expr

var CreateCommand = storage.CreateCommand
var UpdateCommand = storage.UpdateCommand

func NewRelationalDatabase(conn string) (*RelationalDatabase, error) {
	handler, err := storage.NewRelationalDbHandler(conn)
	if err != nil {
		return nil, err
	}
	return &RelationalDatabase{Instance: handler}, nil
}
