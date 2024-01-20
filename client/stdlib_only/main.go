package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	httpInitMsg := []byte("GET / HTTP/1.1\r\nHost:dyte.io\r\nUpgrade:websocket\r\nConnection:Upgrade\r\nSec-WebSocket-Key:dGhlIHNhbXBsZSBub25jZQ==\r\nSec-WebSocket-Version:13\r\nConnection:Upgrade\r\n\r\n")
	wsPayload := []byte{
		// FIN Bit (Final fragment), OpCode (1 for text payload)
		0b10000001,
		// Mask Bit (Required), followed by 7 bits for length (0b0000101 == 5)
		0b10000101,
		// We don't set the extended payload bits as our payload is only 5 bytes
		// Mask (can be any arbritary 32 bit integer)
		0b00000001,
		0b00000010,
		0b00000011,
		0b00000100,
		// Payload, the string "hello" with each character XOR'd with the
		// corresponding mask bits
		0b01101001, // 'h' ^ 0b00000001
		0b01100111, // 'e' ^ 0b00000010
		0b01101111, // 'l' ^ 0b00000011
		0b01101000, // 'l' ^ 0b00000100
		0b01101110, // 'o' ^ 0b00000001
	}

	// Establish a TCP connection to the server
	conn, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	// Send the initial HTTP message to start talking over the WebSocket protocol
	_, err = conn.Write(httpInitMsg)

	if err != nil {
		log.Fatal(err)
	}

	response := make([]byte, 512)

	// Receive the initial HTTP response
	received, err := conn.Read(response)

	if err != nil {
		log.Fatal(err)
	}

	// Write the websocket frame
	_, err = conn.Write(wsPayload)

	if err != nil {
		log.Fatal(err)
	}

	// Read the reply into the existing buffer
	_, err = conn.Read(response[received:])

	fmt.Println(string(response))
}
