package db

import (
	"database/sql"
	db "github.com/Mgeorg1/go-todo-list/db/sqlc"
)

type Store struct {
	*db.Queries
	db *sql.DB
}

func NewStore(conn *sql.DB) Store {
	return Store{
		db:      conn,
		Queries: db.New(conn),
	}
}
