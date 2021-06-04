package main

import (
	"log"
	"net/http"

	"github.com/msal4/poker-wins"
)

const dbFileName = "game.db.json"

func main() {
	store, cleanDB, err := poker.FileSystemPlayerStoreFromFile(dbFileName)
	if err != nil {
		log.Fatalf("error while initilizing the store, %v", err)
	}
	defer cleanDB()

	handler := poker.NewPlayerServer(store)

	log.Fatal(http.ListenAndServe(":5000", handler))
}
