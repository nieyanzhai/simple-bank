package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/nieyanzhai/simple-bank/api"
	db "github.com/nieyanzhai/simple-bank/db/sqlc"
	"github.com/nieyanzhai/simple-bank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		panic(err)
	}
}
