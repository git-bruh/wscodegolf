package main

import (
	"syscall"
	"unsafe"
)

var (
	httpInitMsg = []byte("GET / HTTP/1.1\r\nHost: dyte.io\r\nUpgrade:websocket\r\nConnection:Upgrade\r\nSec-WebSocket-Key:dGhlIHNhbXBsZSBub25jZQ==\r\nSec-WebSocket-Version:13\r\nConnection:Upgrade\r\n\r\n")
	wsPayload = []byte{
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
	// Connects to an IPv4 server at 127.0.0.1 on port 8080
	sockaddr = []byte{
		// family - AF_INET (0x2), padded to 16 bits
		0b00000010,
		0b00000000,
		// port - 8080, padded to 16 bits
		0b00011111,
		0b10010000,
		// addr - 127.0.0.1, 32 bits
		// 127 << 0 | 0 << 8 | 0 << 16 | 1 << 24
		0b01111111,
		0b00000000,
		0b00000000,
		0b00000001,
		// 64 bits of padding
		0b00000000, 0b00000000, 0b00000000, 0b00000000,
		0b00000000, 0b00000000, 0b00000000, 0b00000000,
	}
	// The response buffer for receiving server responses
	response [135]byte
)

//export smol_main
//go:nobounds
func main() {
	// Create a IPv4 (AF_INET), TCP (SOCK_STREAM) socket FD
	// __NR_socket, AF_INET, SOCK_STREAM
	var sock, _, _ = syscall.Syscall(359, 0x2, 0x1, 0)

	// Connect to the server using the `sockaddr_in` structure
	// __NR_connect, fd, sockaddr_in, len(sockaddr_in)
	syscall.Syscall6(362, sock, uintptr(unsafe.Pointer(&sockaddr[0])), uintptr(len(sockaddr)), 0, 0, 0)

	// Send the HTTP message over the socket
	// __NR_sendto, fd, buf, len(buf), flags, addr, addr_len
	syscall.Syscall6(369, sock, uintptr(unsafe.Pointer(&httpInitMsg[0])), uintptr(len(httpInitMsg)), 0, 0, 0)

	// Receive the response
	// __NR_recvfrom, fd, buf, len(buf), flags, addr, addr_len
	var n, _, _ = syscall.Syscall6(371, sock, uintptr(unsafe.Pointer(&response[0])), uintptr(len(response)), 0, 0, 0)

	// Send the WebSocket frame
	// __NR_sendto
	syscall.Syscall6(369, sock, uintptr(unsafe.Pointer(&wsPayload[0])), uintptr(len(wsPayload)), 0, 0, 0)

	// Receive the response
	// __NR_recvfrom
	syscall.Syscall6(371, sock, uintptr(unsafe.Pointer(&response[n])), uintptr(len(response))-n, 0, 0, 0)

	// Close the socket FD
	// __NR_close
	syscall.Syscall(6, sock, 0, 0)

	// Write the response string to standard output
	// __NR_write, STDOUT_FILENO
	syscall.Syscall(4, 1, uintptr(unsafe.Pointer(&response[0])), uintptr(len(response)))

	// Cleanly exit the program with status code 0
	// The libc does this for us in the usual flow, that goes like so:
	//   __libc_start_main (libc) -> main (runtime_unix.go) -> main (main.go)
	// But here, the entrypoint is in main.go itself
	// __NR_exit, EXIT_SUCCESS
	syscall.Syscall(1, 0, 0, 0)
}
