package db

import (
	"context"
	"errors"
	"fmt"
	"log"
)

const (
	lol = iota + 1
	dota
	csgo
	apex
)

type Games struct {
	Response []Game `json:"games"`
}

type Game struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Ranks struct {
	Response []Rank `json:"ranks"`
}

type Rank struct {
	ID   int    `json:"id"`
	Rank string `json:"rank"`
}

func (db *PgxCon) GetAllGames() (*Games, error) {
	connCtx, cancel := context.WithTimeout(context.Background(), waitTimeout)
	defer cancel()
	rows, err := db.pgConn.Query(connCtx,
		"SELECT id,name FROM games")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	games := Games{Response: make([]Game, 0, 0)}
	for rows.Next() {
		var game Game
		err = rows.Scan(&game.ID, &game.Name)
		if err != nil {
			log.Print("error scan game")
			continue
		}
		games.Response = append(games.Response, game)
	}

	if len(games.Response) > 1 {
		return &games, nil
	}

	return nil, errors.New("доступных игр не найдено")
}

func (db *PgxCon) GetGameRanks(gameID int) (*Ranks, error) {
	connCtx, cancel := context.WithTimeout(context.Background(), waitTimeout)
	defer cancel()

	var table string
	switch gameID {
	case lol:
		table = "lol_rank"
	case dota:
		table = "dota_rank"
	case csgo:
		table = "csgo_rank"
	case apex:
		table = "apex_rank"
	}
	rows, err := db.pgConn.Query(connCtx, fmt.Sprintf("SELECT id,rank FROM %s", table))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	ranks := Ranks{Response: make([]Rank, 0, 0)}
	for rows.Next() {
		var rank Rank
		err = rows.Scan(&rank.ID, &rank.Rank)
		if err != nil {
			log.Print("error scan game")
			continue
		}
		log.Printf("found game %v", rank)
		ranks.Response = append(ranks.Response, rank)
	}
	if len(ranks.Response) > 1 {
		return &ranks, nil
	}

	return nil, errors.New("доступных рангов не найдено")
}
