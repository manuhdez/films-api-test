package main

import (
	"log"

	"github.com/manuhdez/films-api-test/database"
	"github.com/manuhdez/films-api-test/internal/http/server"
)

func main() {
	db, err := database.NewPostgresConnection()
	if err != nil {
		panic(err)
	}

	srv := server.New(db)
	log.Fatal(srv.Start())
}
