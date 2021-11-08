package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"strconv"
)

type Link struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

type ResponseLink struct {
	ID     int    `json:"id"`
	URL    string `json:"url"`
	Status int    `json:"status"`
}

func (db *PgxCon) GetAllLinks() (*[]ResponseLink, error) {
	var rows pgx.Rows
	var id int
	var url string
	var status int

	connCtx, cancel := context.WithTimeout(context.Background(), waitTimeout)
	defer cancel()
	rows, err := db.pgConn.Query(connCtx,
		"SELECT id,url FROM app_links")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	links := make([]ResponseLink, 0)
	for rows.Next() {
		err = rows.Scan(&id, &url, &status)
		if err != nil {
			return nil, err
		}
		links = append(links, ResponseLink{id, url, status})
	}

	return &links, nil
}

func (db *PgxCon) GetLinkById(id string) (*ResponseLink, error) {
	var url string
	var status int

	connCtx, cancel := context.WithTimeout(db.pgConnCtx, waitTimeout)
	defer cancel()

	err := db.pgConn.QueryRow(connCtx, "SELECT url,status FROM app_links WHERE id=$1", id).
		Scan(&url, &status)
	if err != nil {
		return nil, err
	}
	s, _ := strconv.Atoi(id)
	link := ResponseLink{
		ID:     s,
		URL:    url,
		Status: status,
	}
	return &link, nil
}

func (db *PgxCon) AddLink(link Link) (int, error) {
	var id int
	connCtx, cancel := context.WithTimeout(db.pgConnCtx, waitTimeout)
	defer cancel()
	tx, _ := db.pgConn.Begin(connCtx)

	err := tx.QueryRow(connCtx, "INSERT INTO app_links (url)"+
		" VALUES ($1) returning id",
		link.URL,
	).Scan(&id)

	if err != nil {
		tx.Rollback(connCtx)
		return 0, err
	}

	tx.Commit(connCtx)
	return id, nil
}

func (db *PgxCon) UpdateStatusLink(status int, id string) error {
	connCtx, cancel := context.WithTimeout(db.pgConnCtx, waitTimeout)
	defer cancel()
	tx, _ := db.pgConn.Begin(connCtx)

	cc, err := tx.Exec(connCtx, "UPDATE app_links SET status=$1 WHERE id=$2", status, id)

	if cc.RowsAffected() <= 0 {
		return fmt.Errorf(" не найдено строк по данному id:%d", id)
	}

	if err != nil {
		tx.Rollback(connCtx)
		return err
	}

	tx.Commit(connCtx)
	return nil
}
