package db

import (
	"context"
	"errors"
)

type Remind struct {
	Id      string `json:"id"`
	From    string `json:"from"`
	Where   string `json:"where"`
	Message string `json:"message"`
	Date    string `json:"date"`
}

func (db *PgxCon) AddRemind(rem Remind) error {
	connCtx, cancel := context.WithTimeout(context.Background(), waitTimeout)
	defer cancel()
	var id int

	tx, _ := db.pgConn.Begin(connCtx)
	err := tx.QueryRow(connCtx,
		"INSERT INTO reminder (from_to,where_to,message,date) VALUES ($1,$2,$3,$4) returning id",
		rem.From, rem.Where, rem.Message, rem.Date).Scan(&id)

	if err != nil {
		tx.Rollback(connCtx)
		return err
	}

	tx.Commit(connCtx)
	if id == 0 {
		return errors.New("не удалось добавить новое напоминание")
	}

	return nil
}
