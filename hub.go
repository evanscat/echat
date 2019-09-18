// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echat

import (
	"github.com/golang/glog"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[string]*Client
	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
	grouper    Grouper
	parser     MessageParser
	logger     Logger
}

func NewHub(option *Option) *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
		grouper:    option.grouper,
		parser:     option.parser,
		logger:     option.logger,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			if client, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.send)
			}
			h.clients[client.ID] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.send)
			}
		case msg := <-h.broadcast:
			if err := h.logger.Log(msg); err != nil {
				glog.Error(err)
			}
			if m, err := h.parser(msg); err != nil {
				glog.Error(err)
			} else {
				h.grouper.ForEach(m.Dest(), msg, h.HandleSend)
			}
		}
	}
}

func (h *Hub) HandleSend(dest string, bt []byte) error {
	if client, ok := h.clients[dest]; ok {
		client.send <- bt
	}
	return nil
}
