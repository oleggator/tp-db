package db

import (
	"github.com/jackc/pgx"
	"log"
)

var conn *pgx.Conn

func InitDB(config pgx.ConnConfig) {
	var err error
	conn, err = pgx.Connect(config)

	if err != nil {
		log.Panic(err)
	}
}
