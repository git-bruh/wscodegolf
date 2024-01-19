package main

import (
	"fmt"
	"log"

	"golang.org/x/net/websocket"
)

func main() {
	url := "ws://localhost:8080/echo"
	ws, err := websocket.Dial(url, "", url)

	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()

	// Write the `hello` message
	if _, err := ws.Write([]byte("hello")); err != nil {
		log.Fatal(err)
	}

	// 512 byte buffer for storing the response
	var response = make([]byte, 512)

	// No. of bytes received
	var received int

	if received, err = ws.Read(response); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Received: %s\n", response[:received])
}
