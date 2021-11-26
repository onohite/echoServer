package db

import (
	"context"
)

type GameProfile struct {
	ID          int    `json:"id"`
	UserID      string `json:"user_id"`
	GameID      int    `json:"game_id"`
	RankID      int    `json:"rank_id"`
	Sex         int    `json:"sex"`
	Age         int    `json:"age"`
	Contact     string `json:"contact"`
	Description string `json:"description"`
}

func (db *PgxCon) FindGameProfile(id int) (string, error) {
	connCtx, cancel := context.WithTimeout(db.pgConnCtx, waitTimeout)
	var uid string
	defer cancel()
	err := db.pgConn.QueryRow(connCtx,
		"SELECT uid WHERE id=$1 FROM game_profile", id).
		Scan(&uid)

	if err != nil {
		return "", err
	}

	return uid, nil
}

func (db *PgxCon) CreateGameProfile(uid string) (int, error) {
	connCtx, cancel := context.WithTimeout(context.Background(), waitTimeout)
	defer cancel()
	var id int

	_ = db.pgConn.QueryRow(connCtx, "SELECT id from game_profile WHERE uid=$1", uid).Scan(&id)
	if id != 0 {
		return id, nil
	}

	tx, _ := db.pgConn.Begin(connCtx)
	err := tx.QueryRow(connCtx,
		"INSERT INTO game_profile (uid) VALUES ($1) returning id",
		uid).Scan(&id)

	if err != nil {
		tx.Rollback(connCtx)
		return 0, err
	}

	tx.Commit(connCtx)
	return id, nil
}

func (db *PgxCon) UpdateGameProfileContact(contact string, uuid string, id int) error {
	connCtx, cancel := context.WithTimeout(db.pgConnCtx, waitTimeout)
	defer cancel()

	tx, _ := db.pgConn.Begin(connCtx)
	_, err := tx.Exec(connCtx, "UPDATE game_profile SET contact=$1 WHERE id=$2 AND user_id=$3",
		contact,
		uuid,
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

func (db *PgxCon) UpdateGameProfileDescription(description string, uuid string, id int) error {
	connCtx, cancel := context.WithTimeout(db.pgConnCtx, waitTimeout)
	defer cancel()

	tx, _ := db.pgConn.Begin(connCtx)
	_, err := tx.Exec(connCtx, "UPDATE game_profile SET description=$1 WHERE id=$2 AND user_id=$3",
		description,
		uuid,
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
