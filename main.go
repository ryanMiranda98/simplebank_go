package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	api "github.com/ryanMiranda98/simplebank/api"
	db "github.com/ryanMiranda98/simplebank/db/sqlc"
	util "github.com/ryanMiranda98/simplebank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	} else {
		log.Println("API server running on", config.ServerAddress)
	}
}
