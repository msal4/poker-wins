package main

import (
	"log"
	"os"

	"github.com/msal4/poker-wins"
)

const dbFileName = "game.db.json"

func main() {
	store, cleanDB, err := poker.FileSystemPlayerStoreFromFile(dbFileName)
	if err != nil {
		log.Fatalf("error while initilizing the store, %v", err)
	}
	defer cleanDB()

	cli := poker.NewCLI(store, os.Stdin)

	cli.PlayPoker()
}
