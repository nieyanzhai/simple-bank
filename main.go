package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/nieyanzhai/simple-bank/api"
	db "github.com/nieyanzhai/simple-bank/db/sqlc"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:password@localhost:5433/simple_bank?sslmode=disable"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal(err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(":8080")
	if err != nil {
		panic(err)
	}
}
