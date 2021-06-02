package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type stubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   []Player
}

func (store *stubPlayerStore) GetPlayerScore(name string) int {
	return store.scores[name]
}

func (store *stubPlayerStore) RecordWin(name string) {
	store.winCalls = append(store.winCalls, name)
}

func (store *stubPlayerStore) GetLeague() []Player {
	return store.league
}

func TestGETPlayers(t *testing.T) {
	store := &stubPlayerStore{
		scores: map[string]int{"Pepper": 20, "Floyd": 10},
	}

	srv := NewPlayerServer(store)
	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		srv.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		srv.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "10")
	})

	t.Run("returns 404 on missing players", func(t *testing.T) {
		request := newGetScoreRequest("Apollo")
		response := httptest.NewRecorder()

		srv.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func TestStoreWins(t *testing.T) {
	store := &stubPlayerStore{
		scores: map[string]int{},
	}

	srv := NewPlayerServer(store)

	t.Run("it returns accepted on POST", func(t *testing.T) {
		const player = "Pepper"
		req := newPostWinRequest(player)
		res := httptest.NewRecorder()

		srv.ServeHTTP(res, req)

		assertStatus(t, res.Code, http.StatusAccepted)

		if len(store.winCalls) != 1 {
			t.Fatalf("got %d to RecordWin, want %d", len(store.winCalls), 1)
		}

		if store.winCalls[0] != player {
			t.Errorf("did not store correct winner got %s, want %s", store.winCalls[0], player)
		}
	})
}

func newPostWinRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func TestLeague(t *testing.T) {
	t.Run("it returns 200 on /league", func(t *testing.T) {
		wantedLeague := []Player{
			{"Cleo", 2},
			{"Chris", 5},
			{"Thiest", 9},
		}
		store := &stubPlayerStore{league: wantedLeague}
		srv := NewPlayerServer(store)

		req := newLeagueRequest()
		res := httptest.NewRecorder()

		srv.ServeHTTP(res, req)

		got := getLeagueFromResponse(t, res.Body)

		assertLeague(t, got, wantedLeague)

		assertStatus(t, res.Code, http.StatusOK)

		assertContentType(t, res, jsonContentType)
	})
}

func newLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league/", nil)
	return req
}

func getLeagueFromResponse(t testing.TB, body io.Reader) (league []Player) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&league)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, %v", body, err)
	}

	return
}

func assertLeague(t testing.TB, got, want []Player) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, wanted %+v", got, want)
	}
}

func assertContentType(t testing.TB, res *httptest.ResponseRecorder, want string) {
	t.Helper()
	if res.Result().Header.Get("content-type") != want {
		t.Fatalf("response did not have content-type of %s, got %v", want, res.Result().Header)
	}
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got status %d, want %d", got, want)
	}
}

func assertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got response %q, want %q", got, want)
	}
}
