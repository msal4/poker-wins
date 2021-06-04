package poker

import (
	"fmt"
	"os"
)

func FileSystemPlayerStoreFromFile(path string) (*FileSystemPlayerStore, func(), error) {
	db, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, nil, fmt.Errorf("could not open db file %s, %v", path, err)
	}

	closeFunc := func() {
		db.Close()
	}

	store, err := NewFileSystemPlayerStore(db)
	if err != nil {
		closeFunc()
		return nil, nil, fmt.Errorf("problem creating file system store, %v", err)
	}

	return store, closeFunc, nil
}
