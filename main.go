package main

import (
	"log"
	"net/http"

	"github.com/msal4/players/server"
)

func main() {
	handler := server.NewPlayerServer(server.NewInMemoryPlayerStore())

	log.Fatal(http.ListenAndServe(":5000", handler))
}
