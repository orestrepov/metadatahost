package db

import (
	"database/sql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"
)

type Database struct {
	*sql.DB
}

// Create new db connection
func New(config *Config) (*Database, error) {
	db, err := sql.Open("postgres", config.DatabaseURI)
	if err != nil {
		return nil, errors.Wrap(err, "unable to connect to database")
	}
	return &Database{db}, nil
}
