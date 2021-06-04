// Package poker ...
package poker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
)

type Player struct {
	Name string
	Wins int
}

type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() League
}

type PlayerServer struct {
	Store PlayerStore
	http.Handler
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := &PlayerServer{
		Store: store,
	}

	router := http.NewServeMux()

	router.Handle("/league/", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))

	p.Handler = router

	return p
}

const jsonContentType = "application/json"

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(p.Store.GetLeague())
}

func (p *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	switch r.Method {
	case http.MethodPost:
		p.processWin(w, player)
	case http.MethodGet:
		p.showScore(w, player)
	}
}

func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {
	p.Store.RecordWin(player)

	w.WriteHeader(http.StatusAccepted)
}

func (p *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score := p.Store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

type FileSystemPlayerStore struct {
	database *json.Encoder
	league   League
}

func NewFileSystemPlayerStore(db *os.File) (*FileSystemPlayerStore, error) {
	db.Seek(0, 0)
	league, err := NewLeague(db)
	if err != nil {
		return nil, fmt.Errorf("problem loading player from file %s, %v", db.Name(), err)
	}

	return &FileSystemPlayerStore{database: json.NewEncoder(&tape{db}), league: league}, nil
}

func (s *FileSystemPlayerStore) GetLeague() (league League) {
	sort.Slice(s.league, func(i, j int) bool {
		return s.league[i].Wins > s.league[j].Wins
	})
	return s.league
}

func (s *FileSystemPlayerStore) GetPlayerScore(name string) int {
	player := s.league.Find(name)

	if player == nil {
		return 0
	}

	return player.Wins
}

func (s *FileSystemPlayerStore) RecordWin(name string) {
	player := s.league.Find(name)
	if player != nil {
		player.Wins++
	} else {
		player = &Player{Name: name, Wins: 1}
		s.league = append(s.league, *player)
	}

	s.database.Encode(s.league)
}
