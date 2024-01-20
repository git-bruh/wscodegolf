package main

import (
	"syscall"
	"unsafe"
)

var buffer [1024]byte
var used uintptr = 0

// We disable the go GC entirely and provide this stub for handling
// allocations, giving out addresses from a static buffer on the stack
// This saves many bytes over using the "leaking" GC, it is more or less
// used exclusively by the runtime's startup code for tasks like setting up
// the processe's environment variables
// If it crashes, run it with a clean environment (env -i ./main)

//go:linkname alloc runtime.alloc
func alloc(size uintptr, layoutPtr unsafe.Pointer) unsafe.Pointer {
	var ptr = unsafe.Pointer(&buffer[used])

	// Align for x64
	used += ((size + 15) &^ 15)

	return ptr
}

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
	// Connects to an IPv4 server at 127.0.0.1 on port 8080
	sockaddr := []byte{
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
	var response [135]byte

	// Create a IPv4 (AF_INET), TCP (SOCK_STREAM) socket FD
	// __NR_socket, AF_INET, SOCK_STREAM
	var sock, _, _ = syscall.Syscall(syscall.SYS_SOCKET, 0x2, 0x1, 0)

	// Connect to the server using the `sockaddr_in` structure
	// __NR_connect, fd, sockaddr_in, len(sockaddr_in)
	syscall.Syscall6(syscall.SYS_CONNECT, sock, uintptr(unsafe.Pointer(&sockaddr[0])), uintptr(len(sockaddr)), 0, 0, 0)

	// Send the HTTP message over the socket
	// __NR_sendto, fd, buf, len(buf), flags, addr, addr_len
	syscall.Syscall6(syscall.SYS_SENDTO, sock, uintptr(unsafe.Pointer(&httpInitMsg[0])), uintptr(len(httpInitMsg)), 0, 0, 0)

	// Receive the response
	// __NR_recvfrom, fd, buf, len(buf), flags, addr, addr_len
	var n, _, _ = syscall.Syscall6(syscall.SYS_RECVFROM, sock, uintptr(unsafe.Pointer(&response[0])), uintptr(len(response)), 0, 0, 0)

	// Send the WebSocket frame
	// __NR_sendto
	syscall.Syscall6(syscall.SYS_SENDTO, sock, uintptr(unsafe.Pointer(&wsPayload[0])), uintptr(len(wsPayload)), 0, 0, 0)

	// Receive the response
	// __NR_recvfrom
	syscall.Syscall6(syscall.SYS_RECVFROM, sock, uintptr(unsafe.Pointer(&response[n])), uintptr(len(response))-n, 0, 0, 0)

	// Close the socket FD
	// __NR_close
	syscall.Syscall(syscall.SYS_CLOSE, sock, 0, 0)

	// Write the response string to standard output
	// __NR_write, STDOUT_FILENO
	syscall.Syscall(syscall.SYS_WRITE, 1, uintptr(unsafe.Pointer(&response[0])), uintptr(len(response)))
}
