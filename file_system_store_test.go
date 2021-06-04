package poker

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("league from a reader", func(t *testing.T) {
		database, removeDatabase := createTempFile(t, `[{"Name": "Pepper", "Wins": 10},{"Name": "Cleo", "Wins": 22}]`)
		defer removeDatabase()
		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		got := store.GetLeague()
		want := []Player{{Name: "Cleo", Wins: 22}, {Name: "Pepper", Wins: 10}}

		assertLeague(t, got, want)
	})

	t.Run("get player score", func(t *testing.T) {
		database, removeDatabase := createTempFile(t, `[{"Name": "Pepper", "Wins": 10},{"Name": "Cleo", "Wins": 22}]`)
		defer removeDatabase()
		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		got := store.GetPlayerScore("Cleo")
		want := 22

		assertScoreEquals(t, got, want)
	})

	t.Run("record win for existing player", func(t *testing.T) {
		const player = "Cleo"

		database, removeDatabase := createTempFile(t, `[{"Name": "Pepper", "Wins": 10},{"Name": "Cleo", "Wins": 22}]`)
		defer removeDatabase()

		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		store.RecordWin(player)

		got := store.GetPlayerScore(player)
		want := 23

		assertScoreEquals(t, got, want)
	})

	t.Run("record win for new player", func(t *testing.T) {
		const player = "Cleo"

		database, removeDatabase := createTempFile(t, `[{"Name": "Pepper", "Wins": 10}]`)
		defer removeDatabase()

		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		store.RecordWin(player)

		got := store.GetPlayerScore(player)
		want := 1

		assertScoreEquals(t, got, want)
	})

	t.Run("league sorted", func(t *testing.T) {
		database, removeDatabase := createTempFile(t, `[{"Name": "Pepper", "Wins": 10},{"Name": "Cleo", "Wins": 22}]`)
		defer removeDatabase()

		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		got := store.GetLeague()
		want := []Player{{"Cleo", 22}, {"Pepper", 10}}

		assertLeague(t, got, want)

		got = store.GetLeague()
		assertLeague(t, got, want)
	})

	t.Run("works with an empty file", func(t *testing.T) {
		database, removeDatabase := createTempFile(t, ``)
		defer removeDatabase()

		_, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)
	})
}

func assertScoreEquals(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}

func createTempFile(t testing.TB, initialData string) (tmpFile *os.File, removeFile func()) {
	tmpFile, err := ioutil.TempFile("", "db")

	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}

	tmpFile.Write([]byte(initialData))

	removeFile = func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}

	return
}

func assertNoError(t testing.TB, err error) {
	if err != nil {
		t.Fatalf("didnt expect an error but got one, %v", err)
	}
}
