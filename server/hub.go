// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"github.com/ivpusic/grpool"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.

//type RegInformation struct {
//	conn *websocket.Conn
//	hashCode string
//}
//numCPUs := runtime.NumCPU()

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}


func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client,),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}



func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Println("Add new client into list")

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Println("Remove clients from list")

				// remove element from map (tt)
			}

		case message := <-h.broadcast:
			//pool := tunny.NewFunc(4, func(payload interface{}) interface{} {
			//	var result []byte
			//
			//	// TODO: Something CPU heavy with payload
			//	for client := range h.clients {
			//		select {
			//	case client.send <- message:
			//	default:
			//		log.Println("Broadcast close connection")
			//		close(client.send)
			//		client.conn.Close()
			//		delete(h.clients, client)
			//		//rm_counter=rm_counter+1
			//	}
			//		//counter = counter + 1
			//	}
			//
			//	return result
			//})
			//
			//defer pool.Close()
			//pool.Process(4)

			//===============================================//
			// number of workers, and size of job queue
			pool := grpool.NewPool(10, 50)

			// release resources used by pool
			defer pool.Release()

			// submit one or more jobs to pool
			var counter int32
			//var rm_counter int32
			counter=0
			//rm_counter=0
			// submit one or more jobs to pool
			// how many jobs we should wait

			len_clients := len(h.clients)
			pool.WaitCount(len_clients)

			for client := range h.clients{

				pool.JobQueue <- func() {
					select {
					case client.send <- message:
						counter = counter + 1
						log.Printf("Broadcast %d users --> %d", len(h.clients), counter)
					default:
						log.Println("Broadcast client.send <- message failed")
						if _, ok := h.clients[client]; ok {
							delete(h.clients, client)
							close(client.send)
							log.Println("Broadcast -->Remove clients from list")

							// remove element from map (tt)
							//rm_counter=rm_counter+1
							//continue
						}
					}
				}
			}
			// WaitAll is an alternative to Results() where you
			// may want/need to wait until all work has been
			// processed, but don't need to check results.
			// eg. individual units of work may handle their own
			// errors, logging...
			pool.WaitAll()

			// dummy wait until jobs are finished
			//time.Sleep(1 * time.Second)
			//var counter int32
			//var rm_counter int32
			//counter=0
			//rm_counter=0
			////len_clients := len(h.clients)
			//
			////defer TimeTaken(time.Now(), "LongRunningFunction")
			//for client := range h.clients {
			//	//client.send <- message
			//	select {
			//	case client.send <- message:
			//		counter = counter + 1
			//	default:
			//		log.Println("Broadcast client.send <- message failed")
			//		if _, ok := h.clients[client]; ok {
			//			delete(h.clients, client)
			//			close(client.send)
			//			log.Println("Broadcast -->Remove clients from list")
			//
			//			// remove element from map (tt)
			//			rm_counter=rm_counter+1
			//			//continue
			//		}
			//	}
			//	//default:
			//	//	log.Println("Broadcast client.send <- message failed")
			//	//	close(client.send)
			//	//	client.conn.Close()
			//	//	delete(h.clients, client)
			//	//	rm_counter=rm_counter+1
			//	//	log.Println("Broadcast client.send <- close client")
			//	//}
			//	//scounter = counter + 1
			//
			//	//if (counter == int32(len_clients/2)){
			//	//	log.Println("Broadcast client.send <- force close socket")
			//	//	close(client.send)
			//	//	client.conn.Close()
			//	//	delete(h.clients, client)
			//	//	rm_counter=rm_counter+1
			//	//}
			//}

			// exit the loop broadcast
			//log.Printf("Broadcast %d users --> %d", len(h.clients), counter)
		}
	}
}
