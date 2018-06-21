// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"log"
	"net/http"
	"time"
	"github.com/gorilla/websocket"
	"github.com/manucorporat/try"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	//pongWait = 60 * time.Second
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 10240
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  10240 * 2,
	WriteBufferSize: 10240 * 2,
	EnableCompression: false, // Compression EXPERIMENTAL
}

type RegInformation struct {
	conn *websocket.Conn
	hashCode string
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// hashCode
	//hashCode string
}

func (c *Client) call_broadcast_msg(x []byte) {
	defer TimeTaken(time.Now(), "call_broadcast_msg")

	//start := time.Now()
	try.This(func() {
		c.hub.broadcast <- x
	}).Catch(func(e try.E){
		log.Println("Broadcast message fail", e)
	})
	//time.Sleep(2 * time.Second)
	//elapsed := time.Since(start)
	//log.Printf("call_broadcast_msg took %s", elapsed)
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	//defer func() {
	//	c.hub.unregister <- c
	//	c.conn.Close()
	//	log.Println("c.hub.unregister <- c")
	//}()
	c.conn.SetReadLimit(maxMessageSize)
	//c.conn.SetReadDeadline(time.Now().Add(pongWait))
	//c.conn.SetPongHandler(func(string) error {
	//	c.conn.SetReadDeadline(time.Now().Add(pongWait));
	//	log.Println("Wait for pong")
	//	return nil
	//	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived) {
				log.Printf("error: %v", err)

				// remove client from list here
				c.hub.unregister <- c
				c.conn.Close()
				log.Println("IsUnexpectedCloseError")
			}

			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived) {
				log.Printf("error: %v", err)

				// remove client from list here
				c.hub.unregister <- c
				c.conn.Close()
				log.Println("IsCloseError")
			}

			// IsCloseError
			break
		}


		//if ce, ok := err.(*websocket.CloseError); ok {
		//	switch ce.Code {
		//	case websocket.CloseNormalClosure,
		//		websocket.CloseGoingAway,
		//		websocket.CloseNoStatusReceived:
		//		log.Printf("error: %v", err)
		//		break
		//	}
		//}

		temp := string(">>>echo:")
		//log.Println(temp)
		b := []byte(temp)
		x := []byte{}

		x = append(x, b...)
		x = append(x, message...)


		x = bytes.TrimSpace(bytes.Replace(x, newline, space, -1))
		c.call_broadcast_msg(x)
	}
}

func (c *Client) sendBroadCast(resp []byte ) {
	for client := range c.hub.clients {
		//select {
		if client != nil {
			client.send <- resp
		}
		//default:
		//	log.Println("Broadcast close connection")
		//	close(client.send)
		//	delete(c.hub.clients, client)
		//}
	}
}

func TimeTaken(t time.Time, name string) {
	elapsed := time.Since(t)
	log.Printf("TIME: %s took %s\n", name, elapsed)
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	//ticker := time.NewTicker(pingPeriod)
	//defer func() {
	//	ticker.Stop()
	//	c.conn.Close()
	//	log.Println("c.conn.Close()")
	//}()
	for {
		select {
		case message, ok := <-c.send:
			//c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println("The hub closed the channel.")
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		//case <-ticker.C:
		//	//c.conn.SetWriteDeadline(time.Now().Add(writeWait))
		//	if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
		//		log.Println("case <-ticker.C: quit game!!!")
		//		return
		//	}else{
		//		log.Println("case <-ticker.C, send ping OK")
		//	}

		}

	}
}


//const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
//const (
//	letterIdxBits = 6                    // 6 bits to represent a letter index
//	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
//	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
//)
//
//func RandStringBytesMaskImpr(n int) string {
//	b := make([]byte, n)
//	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
//	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
//		if remain == 0 {
//			cache, remain = rand.Int63(), letterIdxMax
//		}
//		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
//			b[i] = letterBytes[idx]
//			i--
//		}
//		cache >>= letterIdxBits
//		remain--
//	}
//
//	return string(b)
//}


// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// init Upgrade
	// A server application calls the Upgrader.Upgrade method from an HTTP request handler to get a *Conn

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}


	// Use conn to send and receive messages.
	log.Println("serveWs handles websocket requests from the peer")
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
