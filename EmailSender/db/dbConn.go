package db

import (
	"context"
	"emailsender/config"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

const waitTimeout = 10 * time.Second

type PgxCon struct {
	pgConn    *pgxpool.Pool
	pgConnCtx context.Context
}

func NewPgxCon(ctx context.Context, config *config.Config) (*PgxCon, error) {
	// Контекст ограниченный по времени ожидания
	instance := &PgxCon{}

	instance.pgConnCtx = ctx

	var err error
	var count = 0
	for {
		if count < 6 {
			count++
		}
		err = instance.reconnect(config.PostgresAdress)
		if err != nil {
			log.Printf("connection was lost. Error: %s. Wait %d sec.", err, count*5)
		} else {
			break
		}
		log.Println("Try to reconnect...")
		time.Sleep(time.Duration(count*5) * time.Second)
	}
	return instance, nil
}

func (db *PgxCon) Close() error {
	db.pgConn.Close()
	return nil
}

func (db *PgxCon) reconnect(address string) error {
	connCtx, cancel := context.WithTimeout(db.pgConnCtx, waitTimeout)
	defer cancel()

	conn, err := pgxpool.Connect(connCtx, address)
	if err != nil {
		return fmt.Errorf("unable to connection to database: %v", err)
	}

	db.pgConn = conn
	return nil
}
