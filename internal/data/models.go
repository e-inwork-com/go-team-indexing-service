// Copyright 2023, e-inwork.com. All rights reserved.

package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Teams TeamModel
}

func InitModels(db *sql.DB) Models {
	return Models{
		Teams: TeamModel{DB: db},
	}
}
