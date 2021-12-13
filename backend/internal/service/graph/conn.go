package graph

import (
	"backend/internal/config"
	"context"
	"fmt"
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"google.golang.org/grpc"
	"log"
	"time"
)

const waitTimeout = 10 * time.Second

type GraphConn struct {
	client *dgo.Dgraph
	conn   *grpc.ClientConn
	ctx    context.Context
}

func NewGraphConn(ctx context.Context, config *config.Config) (*GraphConn, error) {
	instance := GraphConn{}
	instance.ctx = ctx
	var err error
	var count = 0
	for {
		if count < 6 {
			count++
		}
		err = instance.reconnect(config.GraphAdress)
		if err != nil {
			log.Printf("connection was lost. Error: %s. Wait %d sec.", err, count*5)
		} else {
			break
		}
		log.Println("Try to reconnect...")
		time.Sleep(time.Duration(count*5) * time.Second)
	}
	return &instance, nil
}

func (g *GraphConn) Close() error {
	return g.conn.Close()
}

func (g *GraphConn) reconnect(address string) error {
	connCtx, cancel := context.WithTimeout(g.ctx, waitTimeout)
	defer cancel()
	conn, err := grpc.DialContext(connCtx, address, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("unable to connection to graphDB: %v", err)
	}
	dgraphClient := dgo.NewDgraphClient(api.NewDgraphClient(conn))
	g.client = dgraphClient
	g.conn = conn
	return nil
}
