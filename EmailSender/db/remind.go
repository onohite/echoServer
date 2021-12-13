package db

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"
)

type Remind struct {
	Id      string    `json:"id"`
	From    string    `json:"from"`
	Where   string    `json:"where"`
	Message string    `json:"message"`
	Date    time.Time `json:"date"`
}

func (db *PgxCon) GetReminds() ([]Remind, error) {
	var id int
	var where, message, from string
	var date time.Time
	connCtx, cancel := context.WithTimeout(context.Background(), waitTimeout)
	defer cancel()

	// '2021-12-13' date template
	timeNow := time.Now().String()
	result := strings.Split(timeNow, " ")
	rows, err := db.pgConn.Query(connCtx, "SELECT id,from_to,where_to,message,date FROM reminder WHERE status=true and date=$1", result[0])

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	list := make([]Remind, 0, 0)
	for rows.Next() {
		err = rows.Scan(&id, &from, &where, &message, &date)
		if err != nil {
			return nil, err
		}
		list = append(list, Remind{strconv.Itoa(id), from, where, message, date})
	}

	if len(list) <= 0 {
		return nil, errors.New("не удалось найти напоминания")
	}

	return list, nil
}

func (db *PgxCon) UpdateStatusReminds(id string) error {
	connCtx, cancel := context.WithTimeout(db.pgConnCtx, waitTimeout)
	defer cancel()

	tx, _ := db.pgConn.Begin(connCtx)
	_, err := tx.Exec(connCtx, "UPDATE reminder SET status=false WHERE id=$1",
		id,
	)

	// Rollback transaction if any error happened on insertion
	if err != nil {
		tx.Rollback(connCtx)
		return err
	}

	tx.Commit(connCtx)
	return nil
}
