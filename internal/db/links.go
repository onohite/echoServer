package db

import (
	"context"
	"github.com/jackc/pgx/v4"
	"strconv"
)

type Link struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

func (db *PgxCon) GetAllLinks() (*[]Link, error) {
	var rows pgx.Rows
	var id int
	var url string

	connCtx, cancel := context.WithTimeout(context.Background(), waitTimeout)
	defer cancel()
	rows, err := db.pgConn.Query(connCtx,
		"SELECT id,url FROM app_links")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	links := make([]Link, 0)
	for rows.Next() {
		err = rows.Scan(&id, &url)
		if err != nil {
			return nil, err
		}
		links = append(links, Link{id, url})
	}

	return &links, nil
}

func (db *PgxCon) GetLinkById(id string) (*Link, error) {
	var url string

	connCtx, cancel := context.WithTimeout(db.pgConnCtx, waitTimeout)
	defer cancel()

	err := db.pgConn.QueryRow(connCtx, "SELECT url FROM app_links WHERE id=$1", id).
		Scan(&url)
	if err != nil {
		return nil, err
	}
	s, _ := strconv.Atoi(id)
	link := &Link{
		ID:  s,
		URL: url,
	}
	return link, nil
}

func (db *PgxCon) AddLink(link Link) error {
	connCtx, cancel := context.WithTimeout(db.pgConnCtx, waitTimeout)
	defer cancel()
	tx, _ := db.pgConn.Begin(connCtx)

	_, err := tx.Exec(connCtx, "INSERT INTO app_links (url) values($1)",
		link.URL,
	)

	if err != nil {
		tx.Rollback(connCtx)
		return err
	}

	tx.Commit(connCtx)
	return nil
}
