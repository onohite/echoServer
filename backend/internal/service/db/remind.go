package db

import (
	"context"
	"errors"
	"strconv"
	"time"
)

type Remind struct {
	Id      string `json:"id"`
	From    string `json:"from"`
	Where   string `json:"where"`
	Message string `json:"message"`
	Date    string `json:"date"`
}

func (db *PgxCon) GetListReminds(from string) ([]Remind, error) {
	var id int
	var where, message string
	var date time.Time
	connCtx, cancel := context.WithTimeout(context.Background(), waitTimeout)
	defer cancel()
	rows, err := db.pgConn.Query(connCtx, "SELECT id,from_to,where_to,message,date FROM reminder WHERE from_to=$1", from)

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
		list = append(list, Remind{strconv.Itoa(id), from, where, message, date.Format(`2006-01-02`)})
	}

	if len(list) <= 0 {
		return nil, errors.New("не удалось найти напоминания")
	}

	return list, nil
}

func (db *PgxCon) UpdateRemindTo(id, to string) error {
	connCtx, cancel := context.WithTimeout(db.pgConnCtx, waitTimeout)
	defer cancel()

	tx, _ := db.pgConn.Begin(connCtx)
	_, err := tx.Exec(connCtx, "UPDATE reminder SET where_to=$1 WHERE id=$2",
		to,
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

func (db *PgxCon) UpdateRemindMessage(id, message string) error {
	connCtx, cancel := context.WithTimeout(db.pgConnCtx, waitTimeout)
	defer cancel()

	tx, _ := db.pgConn.Begin(connCtx)
	_, err := tx.Exec(connCtx, "UPDATE reminder SET message=$1 WHERE id=$2",
		message,
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

func (db *PgxCon) UpdateRemindDate(id string, date string) error {
	connCtx, cancel := context.WithTimeout(db.pgConnCtx, waitTimeout)
	defer cancel()
	t, _ := time.Parse(`2006-01-02`, date)
	tx, _ := db.pgConn.Begin(connCtx)
	_, err := tx.Exec(connCtx, "UPDATE reminder SET date=$1 WHERE id=$2",
		t,
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
