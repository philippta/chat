package main

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

// Server represents a websocket chat server
type Server struct {
	clients   map[*websocket.Conn]*sync.Mutex
	broadcast chan []byte
	add       chan *websocket.Conn
	del       chan *websocket.Conn
}

// NewServer creates a new server
func NewServer() *Server {
	return &Server{
		broadcast: make(chan []byte, 256),
		add:       make(chan *websocket.Conn, 256),
		del:       make(chan *websocket.Conn, 256),
		clients:   make(map[*websocket.Conn]*sync.Mutex),
	}
}

// Run starts the server
func (s *Server) Run() {
	for {
		select {
		case msg := <-s.broadcast:
			for conn, m := range s.clients {
				m.Lock()
				conn.WriteMessage(websocket.TextMessage, msg)
				m.Unlock()
			}
		case conn := <-s.add:
			s.clients[conn] = new(sync.Mutex)
			go s.listen(conn)
		case conn := <-s.del:
			if _, ok := s.clients[conn]; ok {
				delete(s.clients, conn)
			}
		}
	}
}

// Add adds a new client to the server
func (s *Server) Add(conn *websocket.Conn) {
	s.add <- conn
}

// listen starts listening for websocket messages and sends them to the broadcast channel
func (s *Server) listen(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.del <- conn
			return
		}
		fmt.Println("got msg:", string(msg))
		s.broadcast <- msg
	}
}
