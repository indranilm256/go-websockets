package api

import (
	"fmt"
	"go-websockets/db"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/gorilla/websocket"
)

type myStruct struct {
	Rollno string `json:"rollno"`
	Words  string `json:"words"`
	Chars  string `json:"chars"`
	Wpmin  string `json:"wpmin"`
}

func ServeTeacher(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/teacher.html")
}

func SocketTeacher(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Teacher Successfully connected..")

	go func(conn *websocket.Conn) {
		ch := time.Tick(time.Second)

		for range ch {
			// txn := db.NewTransaction(true)
			err = db.DB.View(func(txn *badger.Txn) error {
				iopt := badger.DefaultIteratorOptions
				itr := txn.NewIterator(iopt)
				defer itr.Close()
				i := 0
				for itr.Rewind(); itr.Valid(); itr.Next() {
					key := string(itr.Item().Key())
					//log.Println(key)
					val, _ := strconv.Atoi(itr.Item().Value(fnc).Error())

					if key[len(key)-1] == 'w' {
						roll := key[0 : len(key)-1]
						rollc := roll + "c"
						rollstarttime := roll + "st"
						entry2, err := txn.Get([]byte(rollc))
						if err != nil {
							log.Println(err)
						}
						nochars, _ := strconv.Atoi(entry2.Value(fnc).Error())
						//log.Println(val)
						entry3, err := txn.Get([]byte(rollstarttime))
						if err != nil {
							log.Println(err)
						}
						starttime, _ := strconv.Atoi(entry3.Value(fnc).Error())
						timeNow := time.Now()
						secNow := timeNow.Unix()
						wpmin := float64(val*60) / float64(int(secNow)-starttime)
						conn.WriteJSON(myStruct{
							Rollno: roll,
							Words:  strconv.Itoa(val),
							Chars:  strconv.Itoa(nochars),
							Wpmin:  fmt.Sprintf("%f", wpmin),
						})
					}
					i++
				}
				return nil

			})

		}
	}(conn)
}
