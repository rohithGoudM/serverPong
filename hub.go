// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "log"

type Human struct{
    Name string
    address string
    Age int
}

type State struct{
	Player1 []int
	Player2 []int
	Ball int
	BallSpeed int
	DeltaY int
	DeltaX int
	Pause bool
	Player1Score int
	Player2Score int
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan State

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	//json
	state State
}

func (h *Hub) resetState() {
	h.state.Player1[0] = 103
	h.state.Player1[1] = 124
	h.state.Player1[2] = 145
	h.state.Player2[0] = 85
	h.state.Player2[1] = 106
	h.state.Player2[2] = 127
	h.state.Ball = 115
	h.state.BallSpeed = 135
	h.state.DeltaY=-21
	h.state.DeltaX=-1
	h.state.Pause=true
	h.state.Player1Score=0
	h.state.Player2Score=0
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan State),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		state: State{[]int{103, 124, 145},[]int{85, 106, 127},115,135,-21,-1,true,0,0},
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			// log.Println(h.state)
			// log.Println(h.state.player1)
			client.send <- h.state
			log.Println("no of clients:",len(h.clients))
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Println("no of clients:",len(h.clients))
				if len(h.clients)==0{
					h.resetState()
				}
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
