package main

import (
	"database/sql"
	"sync"
	"techschool/api"
	db "techschool/db/sqlc"
	"techschool/util"

	_ "github.com/lib/pq"
)

func main() {

	config, err := util.LoadConfig(".env")
	if err != nil {
		panic(err)
	}

	pg, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		panic(err)
	}

	store := db.NewStore(pg)
	server := api.NewServer(store, config)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {

		defer wg.Done()
		server.Start(config.ServerAddress)
	}()

	wg.Wait()

}
