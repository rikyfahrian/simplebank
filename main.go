package main

import (
	"database/sql"
	"sync"
	"techschool/api"
	db "techschool/db/sqlc"

	_ "github.com/lib/pq"
)

func main() {

	pg, err := sql.Open("postgres", "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable")
	if err != nil {
		panic(err)
	}

	store := db.NewStore(pg)
	server := api.NewServer(store)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {

		defer wg.Done()
		server.Start(":8080")
	}()

	wg.Wait()

}
