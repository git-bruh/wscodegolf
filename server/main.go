package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()

		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("recv: %s", message)

		if string(message) != "hello" {
			log.Fatalf("incorrect msg sent:", string(message))
		}

		err = c.WriteMessage(mt, []byte("dyte"))

		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", echo)

	fmt.Println("starting server at:", *addr)
	http.ListenAndServe(*addr, nil)
}
