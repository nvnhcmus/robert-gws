package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"fmt"
	"encoding/json"
	"math/rand"
)

type Message struct {
	Name string
	Body string
	Time int64
}



var ip_addr = flag.String("addr", "192.168.3.140:8080", "http service address")

func init() {
	rand.Seed(time.Now().UnixNano())
}
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func reconnect(type_reconnect int) *websocket.Conn{
	// create ticket time for every 3 seconds
	ticker := time.NewTicker(time.Duration(3) * time.Second)
	var count = 1
	for _ = range ticker.C {
		// Construct the url
		u := url.URL{Scheme: "ws", Host: *ip_addr, Path: "/wsPriceChart"}
		log.Printf("re-connecting to %s, type=%d", u.String(), type_reconnect)

		// Create new websocket
		fmt.Printf("\nRetry Connect : %d times\n", count)
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)


		if err != nil {
			fmt.Printf("Dial failed: %s\n\n", err.Error())
		} else {
			fmt.Println("Re-conected to %s, type=%d, %p", u.String(), type_reconnect)
			return c
		}
		count = count + 1
	}
	return nil
}




func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *ip_addr, Path: "/wsPriceChart"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				c.Close()

				//if c == nil {
					log.Println("Call re-connect from ReadMessage")
					c = reconnect(1)
					log.Printf("recv prt: %p", c)
				//}
			}
			log.Printf("recv: %s, %p", message, c)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			temp := string("goClang44444444444444:") + t.String()
			//log.Println(temp)
			//b := []byte(temp)
			//x := []byte{}

			//x = append(x, b...)
			////x = append(x, message...)
			//k2line = kline{
			//	c:"0.13600000",
			//	t:1528715160000,
			//	v:"101.00000000",
			//	h:"0.13600000",
			//	l:"0.13600000",
			//	o:"0.13600000",
			//}

			genString:=RandStringRunes(1024 * 5)
			m := Message{
				Name: "Alice" + temp,
				Body: genString,
				Time: 1294706395881547000,
			}

			b, err := json.Marshal(m)

			err = c.WriteMessage(websocket.TextMessage, []byte(b))
			if err != nil {
				log.Println("write:", err)
				//c.Close()

				//if c == nil {
				//	log.Println("Call re-connect from WriteMessage")
				//	c = reconnect(2)
				//	log.Printf("write: %p", c)
				//}
			}//else{
			//	log.Printf("write: %p", c)
			//}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
