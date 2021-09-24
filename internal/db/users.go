package db

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (db *PgxCon) GetAllUsers() (*[]User, error) {
	var rows pgx.Rows
	var id int
	var name string

	connCtx, cancel := context.WithTimeout(context.Background(), waitTimeout)
	defer cancel()
	rows, err := db.pgConn.Query(connCtx,
		"SELECT id,name FROM app_user")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make([]User, 0, 10)
	for rows.Next() {
		err = rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}
		users = append(users, User{id, name})
	}

	return &users, nil
}
