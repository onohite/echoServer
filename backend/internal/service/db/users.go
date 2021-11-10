package db

import (
	"context"
)

type User struct {
	Name       string
	Email      string
	AvatarLink string
	Sex        int
	Bdate      string
	Unique     string
}

func (db *PgxCon) AddUser(u User) (string, error) {
	connCtx, cancel := context.WithTimeout(context.Background(), waitTimeout)
	defer cancel()
	var id string

	_ = db.pgConn.QueryRow(connCtx, "SELECT id from users WHERE unique_identificator=$1", u.Unique).Scan(&id)
	if id != "" {
		return id, nil
	}

	tx, _ := db.pgConn.Begin(connCtx)
	err := tx.QueryRow(connCtx,
		"INSERT INTO Users (user_name,email,avatar_link,sex,bdate,unique_identificator) VALUES ($1,$2,$3,$4,$5,$6) returning id",
		u.Name, u.Email, u.AvatarLink, u.Sex, u.Bdate, u.Unique).Scan(&id)

	if err != nil {
		return "", err
	}

	if err != nil {
		tx.Rollback(connCtx)
		return "", err
	}

	tx.Commit(connCtx)
	return id, nil
}
