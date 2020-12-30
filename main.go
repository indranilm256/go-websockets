package main

import (
	"fmt"
	"go-websockets/api"
	"go-websockets/db"
	"log"
	"net/http"

	"github.com/dgraph-io/badger"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	api.Open()
}

func main() {

	// key := func(i int) []byte {
	// 	return []byte(fmt.Sprintf("%d", i))
	// }

	fmt.Println("Hello There")
	http.HandleFunc("/student", api.ServeStudent)
	http.HandleFunc("/student/ws", api.SocketStudent)

	http.HandleFunc("/teacher", api.ServeTeacher)
	http.HandleFunc("/teacher/ws", api.SocketTeacher)

	var err = db.DB.Update(func(txn *badger.Txn) error {
		iopt := badger.DefaultIteratorOptions
		itr := txn.NewIterator(iopt)
		defer itr.Close()
		i := 0
		for itr.Rewind(); itr.Valid(); itr.Next() {
			if err := txn.Delete(itr.Item().Key()); err != nil {
				log.Println(err)
			}
			i++
		}

		fmt.Println("Deleted", i, "Keys")
		return nil

	})

	check(err)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
