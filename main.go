package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
	server := NewServer()
	tmpl := template.Must(template.ParseFiles("index.html"))
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		server.Add(conn)
	})

	go server.Run()
	http.ListenAndServe(":500", nil)
}
