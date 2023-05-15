package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Server struct {
	Upgrader websocket.Upgrader
	Clients  map[*websocket.Conn]bool
	Mu       sync.Mutex // Protects the Clients map
	Data     []VesselData
}

func NewServer(filename string) *Server {
	data, err := ReadCSV(filename)
	if err != nil {
		log.Fatalf("error: %v | input: ReadCSV(%s)", err, filename)
	}

	return &Server{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		Clients: make(map[*websocket.Conn]bool),
		Data:    data,
	}
}

func (s *Server) Start(address string) {
	http.HandleFunc("/ws", s.handleConnections)

	go s.sendUpdatesToClients()

	log.Println("Server starting on", address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatalf("error: %v | input: http.ListenAndServe(%s, %v)", err, address, nil)
	}
}

func (s *Server) handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("error: %v | input: Upgrader.Upgrade(%v, %v, %v)", err, w, r, nil)
	}

	// Add the client to the Clients map.
	s.Mu.Lock()
	s.Clients[ws] = true
	s.Mu.Unlock()

	go s.handleClient(ws)
}

func (s *Server) handleClient(ws *websocket.Conn) {
	defer ws.Close()

	for {
		// Check for a close message from the client.
		_, _, err := ws.ReadMessage()
		if err != nil {
			// If there's an error (like a closed connection), remove the client.
			s.Mu.Lock()
			delete(s.Clients, ws)
			s.Mu.Unlock()
			break
		}
	}
}

func (s *Server) sendUpdatesToClients() {
	for i := 0; i < len(s.Data); i++ {
		dataPoint := s.Data[i]
		dataPointJSON, _ := json.Marshal(dataPoint)

		s.Mu.Lock()
		for client := range s.Clients {
			err := client.WriteMessage(websocket.TextMessage, dataPointJSON)
			if err != nil {
				log.Printf("error: %v | input: client.WriteMessage(1, %+v)", err, dataPointJSON)
				client.Close()
				delete(s.Clients, client)
				continue
			}
		}
		s.Mu.Unlock()
		time.Sleep(time.Second * 1)
	}
}
