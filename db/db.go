package db

import (
	"github.com/jackc/pgx"
	"log"
)

var conn *pgx.ConnPool

func InitDB(config pgx.ConnConfig) {
	var err error
	conn, err = pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     config,
		MaxConnections: 100,
	})

	if err != nil {
		log.Panic(err)
	}
}
