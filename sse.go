package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"time"
)

type Broker struct {
	clients        map[chan string]bool
	newClients     chan chan string
	defunctClients chan chan string
	messages       chan string
}

func (b *Broker) Start() {

	go func() {
		for {
			select {

			case s := <-b.newClients:
				b.clients[s] = true
				log.Println("Added new client")

			case s := <-b.defunctClients:
				delete(b.clients, s)
				log.Println("Removed client")

			case msg := <-b.messages:
				for s, _ := range b.clients {
					s <- msg
				}
				log.Printf("Broadcast message to %d clients", len(b.clients))
			}
		}
	}()
}

func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	messageChan := make(chan string)
	b.newClients <- messageChan
	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		b.defunctClients <- messageChan
		log.Println("HTTP connection just closed.")
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {

		msg := <-messageChan

		fmt.Fprintf(w, "data: Message: %s\n\n", msg)
		f.Flush()
	}

	log.Println("Finished HTTP request at ", r.URL.Path)
}
func main() {

	b := &Broker{
		make(map[chan string]bool),
		make(chan (chan string)),
		make(chan (chan string)),
		make(chan string),
	}
	b.Start()
	http.Handle("/events/", b)

	go func() {
		var lm35val string
		// open db connection
		db, err := sql.Open("mysql", "root:adm1n01@/test")
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()

		for {

			// select row
			rows, err := db.Query("select suhu from datasuhu order by waktu desc limit 1")
			if err != nil {
				panic(err.Error())
			}
			defer rows.Close()

			for rows.Next() {
				err = rows.Scan(&lm35val)
				if err != nil {
					panic(err.Error())
				}
			}

			b.messages <- fmt.Sprintln(lm35val)
			time.Sleep(5 * 1e9)
		}
	}()

	http.Handle("/", http.FileServer(http.Dir("www")))
	http.ListenAndServe(":8000", nil)
}
