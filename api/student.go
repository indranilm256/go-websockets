package api

import (
	"errors"
	"fmt"
	"go-websockets/db"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/gorilla/websocket"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func fnc(val []byte) error {
	// log.Println(val)
	return errors.New(string(val))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Open() {
	opt := badger.DefaultOptions("./data")
	var err error
	db.DB, err = badger.Open(opt)
	check(err)
}

func ServeStudent(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
	fmt.Fprintf(w, "Holla")
}

func SocketStudent(w http.ResponseWriter, r *http.Request) {

	val := func(i int) []byte {
		return []byte(fmt.Sprintf("%0128d", i))
	}

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
			rollw := roll + "w"
			rollc := roll + "c"
			rollstarttime := roll + "st"

			wb := db.DB.NewWriteBatch()
			defer wb.Cancel()

			txn1 := db.DB.NewTransaction(false)
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
				txn3 := db.DB.NewTransaction(false)
				defer txn3.Discard()
				entry2, err := txn3.Get([]byte(rollc))
				check(err)
				nochars, _ := strconv.Atoi(entry2.Value(fnc).Error())
				check(wb.Set([]byte(rollw), val(nowords+1)))
				check(wb.Set([]byte(rollc), val(nochars+count)))
			}

			check(wb.Flush())

			txn2 := db.DB.NewTransaction(false)
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
}
