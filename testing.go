package poker

import "testing"

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   League
}

func (store *StubPlayerStore) GetPlayerScore(name string) int {
	return store.scores[name]
}

func (store *StubPlayerStore) RecordWin(name string) {
	store.winCalls = append(store.winCalls, name)
}

func (store *StubPlayerStore) GetLeague() League {
	return store.league
}

func AssertPlayerWin(t testing.TB, playerStore *StubPlayerStore, winner string) {
	t.Helper()

	if len(playerStore.winCalls) != 1 {
		t.Fatalf("got %d calls to RecordWin want %d", len(playerStore.winCalls), 1)
	}

	got := playerStore.winCalls[0]

	if got != winner {
		t.Fatalf("did not record the correct winner, got %q want %q", got, winner)
	}
}
