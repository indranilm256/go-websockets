package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type myStruct struct {
	Rollno string `json:"rollno"`
	Words  string `json:"words"`
	Chars  string `json:"chars"`
	Wpmin  string `json:"wpmin"`
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func fnc(val []byte) error {
	// log.Println(val)
	return errors.New(string(val))
}

func main() {

	opt := badger.DefaultOptions("./data")
	db, err := badger.Open(opt)
	check(err)

	defer db.Close()

	// key := func(i int) []byte {
	// 	return []byte(fmt.Sprintf("%d", i))
	// }

	val := func(i int) []byte {
		return []byte(fmt.Sprintf("%0128d", i))
	}

	fmt.Println("Hello There")
	http.HandleFunc("/student", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
		fmt.Fprintf(w, "Holla")
	})
	http.HandleFunc("/student/ws", func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
		}

		log.Println("Student Successfully connected..")

		go func(conn *websocket.Conn) {
			for {
				messageType, p, err := conn.ReadMessage()
				if err != nil {
					log.Println(err)
					return
				}
				temp := string(p)
				//log.Println(temp)

				flag := 0
				for i := 0; i < len(temp)-1; i++ {
					if temp[i] == '_' {
						flag = i
						break
					}
				}

				roll := temp[0:flag]
				count, _ := strconv.Atoi(temp[flag+1:])
				// log.Println(roll)
				// if roll == "roll" {
				// 	println("YES")
				// } else {
				// 	println("No")
				// }
				// log.Println(count)
				//log.Println(temp)
				rollw := roll + "w"
				rollc := roll + "c"
				rollstarttime := roll + "st"
				// log.Println(rollw)

				wb := db.NewWriteBatch()
				defer wb.Cancel()

				txn1 := db.NewTransaction(false)
				defer txn1.Discard()
				entry, err := txn1.Get([]byte(rollw))
				if err != nil {
					check(wb.Set([]byte(rollw), val(1)))
					check(wb.Set([]byte(rollc), val(count)))
					now := time.Now()
					sec := now.Unix()
					check(wb.Set([]byte(rollstarttime), val(int(sec))))
					// check(wb.Set([]byte(rollstarttime), []byte(strconv.Itoa(int(sec)))))
				} else {
					nowords, _ := strconv.Atoi(entry.Value(fnc).Error())
					//log.Println(nowords)
					txn3 := db.NewTransaction(false)
					defer txn3.Discard()
					entry2, err := txn3.Get([]byte(rollc))
					check(err)
					nochars, _ := strconv.Atoi(entry2.Value(fnc).Error())
					check(wb.Set([]byte(rollw), val(nowords+1)))
					check(wb.Set([]byte(rollc), val(nochars+count)))
				}

				check(wb.Flush())

				txn2 := db.NewTransaction(false)
				defer txn2.Discard()
				entry, err = txn2.Get([]byte(rollc))
				check(err)
				valv, _ := strconv.Atoi(entry.Value(fnc).Error())
				fmt.Printf("Read val '%d' using txn.Get\n", valv)

				if err := conn.WriteMessage(messageType, p); err != nil {
					log.Println(err)
					return
				}
			}

		}(conn)
	})

	http.HandleFunc("/teacher", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "teacher.html")
	})
	http.HandleFunc("/teacher/ws", func(w http.ResponseWriter, r *http.Request) {
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
				err = db.View(func(txn *badger.Txn) error {
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
							log.Println(wpmin)
							conn.WriteJSON(myStruct{
								Rollno: roll,
								Words:  strconv.Itoa(val),
								Chars:  strconv.Itoa(nochars),
								Wpmin:  fmt.Sprintf("%f", wpmin),
							})
						}
						i++
					}

					//fmt.Println("Read", i, "Keys")
					return nil

				})

			}
		}(conn)
	})

	// http.HandleFunc("/teacher/ws", )

	// txn1 := db.NewTransaction(true)
	// defer txn1.Discard()

	// check(txn1.Set([]byte(key(300)), []byte("bVal")))

	// fmt.Printf("inserted key '%s' using txn.Set\n", key(300))

	// txn2 := db.NewTransaction(false)
	// entry, err := txn2.Get([]byte(key(300)))
	// check(err)
	// fmt.Printf("Read key '%s' using txn.Get\n", string(entry.Key()))
	// N, M := 20, 0

	// wb := db.NewWriteBatch()
	// defer wb.Cancel()
	// check(wb.Set([]byte("bKey"), []byte("bVal")))

	// for i := 0; i < N; i++ {
	// 	check(wb.Set(key(i), val(i)))
	// }

	// for i := 0; i < M; i++ {
	// 	check(wb.Delete(key(i)))
	// }

	// check(wb.Flush())

	// fmt.Println("Inserted", N, "Deleted", M)

	err = db.View(func(txn *badger.Txn) error {
		iopt := badger.DefaultIteratorOptions
		itr := txn.NewIterator(iopt)
		defer itr.Close()
		i := 0
		for itr.Rewind(); itr.Valid(); itr.Next() {
			log.Println(string(itr.Item().Key()))
			val, _ := strconv.Atoi(itr.Item().Value(fnc).Error())
			log.Println(val)
			i++
		}

		fmt.Println("Read", i, "Keys")
		return nil

	})

	check(err)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
