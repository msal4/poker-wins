package main

import (
	"log"
	"net/http"
	"os"

	"github.com/msal4/players/server"
)

const dbFileName = "game.db.json"

func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 06666)
	if err != nil {
		log.Fatalf("problem opening %s, %v", dbFileName, err)
	}
	store, err := server.NewFileSystemPlayerStore(db)
	if err != nil {
		log.Fatalf("problem creating file system player store, %v", err)
	}
	handler := server.NewPlayerServer(store)

	log.Fatal(http.ListenAndServe(":5000", handler))
}
