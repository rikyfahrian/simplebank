package main

import (
	"database/sql"
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

	server.Start(":8080")

}
