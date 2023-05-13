package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Server struct {
	Upgrader  websocket.Upgrader
	Clients   map[*websocket.Conn]bool
	Broadcast chan []byte
	Ticker    *time.Ticker
	Data      []VesselData
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
		Broadcast: make(chan []byte),
		Ticker: time.NewTicker(time.Second * 2),
		Data: data,
	}
}

func (s *Server) Start(address string) {
	http.HandleFunc("/ws", s.handleConnections)

	// Start a goroutine that sends data on the Broadcast channel every tick
	go func() {
		index := 0
		for range s.Ticker.C {
			// Send the current data point
			data, err := json.Marshal(s.Data[index])
			if err != nil {
				log.Fatalf("error: %v | input: json.Marshal(%v)", err, s.Data[index])
			}
			s.Broadcast <- data

			// Increment the index, and loop back to the start if necessary
			index++
			if index >= len(s.Data) {
				index = 0
			}
		}
	}()	

	// Start a goroutine that sends broadcast messages to all clients
	go func() {
		for {
			// Grab the next broadcast message
			msg := <-s.Broadcast

			// Send it out to every client
			for client := range s.Clients {
				err := client.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					log.Printf("error: %v | input: client.WriteMessage(%v, %v)", err, websocket.TextMessage, msg)
					client.Close()
					delete(s.Clients, client)
				}
			}
		}
	}()

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
	defer ws.Close()

	s.Clients[ws] = true
}
