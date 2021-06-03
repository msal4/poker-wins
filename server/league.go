package server

import (
	"encoding/json"
	"fmt"
	"io"
)

type League []Player

func (l League) Find(name string) (player *Player) {
	for i, p := range l {
		if p.Name == name {
			player = &l[i]
			break
		}
	}

	return
}

func NewLeague(r io.Reader) (League, error) {
	var league League
	err := json.NewDecoder(r).Decode(&league)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("problem parsing league: %v", err)
	}

	return league, nil
}
