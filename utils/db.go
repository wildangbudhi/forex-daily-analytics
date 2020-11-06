package utils

import (
	"fmt"

	"github.com/go-pg/pg"
)

func NewDbConnection(address, user, password, database string) (*pg.DB, error) {

	opts := &pg.Options{
		User:     user,
		Password: password,
		Addr:     address,
		Database: database,
	}

	var db *pg.DB = pg.Connect(opts)

	if db == nil {
		return nil, fmt.Errorf("Database Connection Failed")
	}

	return db, nil

}
