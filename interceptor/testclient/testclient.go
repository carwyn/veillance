package main

import (
	"fmt"
	"log"

	"golang.org/x/net/websocket"
)

func main() {

	origin := "http://localhost/"
	url := "ws://localhost:8080/entry"
	ws, err := websocket.Dial(url, "", origin)

	if err != nil {
		log.Fatal(err)
	}

	//if _, err := ws.Write([]byte("hello, world!\n")); err != nil {
	//	log.Fatal(err)
	//}

	var msg = make([]byte, 4096)
	var n int

	for {
		if n, err = ws.Read(msg); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Received: %s.\n", msg[:n])
	}

}
