package db

import "github.com/dgraph-io/badger"

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var DB *badger.DB
