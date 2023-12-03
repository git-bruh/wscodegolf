package main

import (
	"syscall"
	"unsafe"
)

// tinygo build -no-debug -scheduler=none -gc=none -panic=trap -target=spec.json
// strip --strip-all --strip-section-headers -R .comment -R .note -R .eh_frame sys

var buffer [2048]byte
var used uintptr = 0

//go:linkname alloc runtime.alloc
func alloc(size uintptr, layoutPtr unsafe.Pointer) unsafe.Pointer {
	// TODO align
	var ptr = unsafe.Pointer(&buffer[used])
	used += size

	return ptr
}

func main() {
	var httpInitMsg = []byte("GET /echo HTTP/1.1\r\nHost: localhost.com:8080\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==\r\nSec-WebSocket-Version: 13\r\nConnection: keep-alive, Upgrade\r\nSec-Fetch-Mode: websocket\r\n\r\n")
	var packet = []byte{
		0b10000001, // FIN, RSV1, RSV2, RSV3, OpCode
		0b10000101, // Mask Bit (Compulsary for client to set) + Payload
		// NOTE: We don't need to set extended payload bits if our
		// msg is less than 126 length
		0b00000001,
		0b00000010,
		0b00000011,
		0b00000100, // Mask
		0b01101001,
		0b01100111,
		0b01101111,
		0b01101000,
		0b01101110, // Payload
	}
	var response [135]byte

	var sa = syscall.SockaddrInet4{
		Port: 8080,
		Addr: [4]byte{127, 0, 0, 1},
	}

	var sock, _ = syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	syscall.Connect(sock, &sa)

	syscall.Sendto(sock, httpInitMsg, 0, nil)
	var n, _, _ = syscall.Recvfrom(sock, response[:], 0)
	syscall.Sendto(sock, packet, 0, nil)
	syscall.Recvfrom(sock, response[n:], 0)

	syscall.Close(sock)
	syscall.Write(1, response[:])
}
