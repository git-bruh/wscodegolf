package main

import (
	"syscall"
	"unsafe"
)

// tinygo build -no-debug -scheduler=none -gc=none -panic=trap -target=spec.json
// strip --strip-all --strip-section-headers -R .comment -R .note -R .eh_frame sys
// $ wc -c sys
//   6592 sys

var buffer [1024]byte
var used uintptr = 0

// We disable the go GC entirely and provide this stub for handling
// allocations, giving out addresses from a static buffer on the stack
// This saves ~152 bytes over using the "leaking" GC, it is more or less
// used exclusively by the runtime's startup code for tasks like setting up
// the processe's environment variables
// If it crashes, run it with a clean environment (env -i ./sys)

//go:linkname alloc runtime.alloc
func alloc(size uintptr, layoutPtr unsafe.Pointer) unsafe.Pointer {
	var ptr = unsafe.Pointer(&buffer[used])

	// Align for x64
	used += ((size + 15) &^ 15)

	return ptr
}

func main() {
	var httpInitMsg = []byte("GET / HTTP/1.1\r\nHost:\r\nUpgrade:websocket\r\nConnection:Upgrade\r\nSec-WebSocket-Key:dGhlIHNhbXBsZSBub25jZQ==\r\nSec-WebSocket-Version:13\r\nConnection:Upgrade\r\n\r\n")
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
	var sockaddr = [16]byte{
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
	var response [135]byte

	// __NR_socket, AF_INET, SOCK_STREAM
	var sock, _, _ = syscall.Syscall(41, 0x2, 0x1, 0)

	// __NR_connect, fd, sockaddr_in, len(sockaddr_in)
	syscall.Syscall6(42, sock, uintptr(unsafe.Pointer(&sockaddr[0])), uintptr(len(sockaddr)), 0, 0, 0)

	// __NR_sendto, fd, buf, len(buf), flags, addr, addr_len
	syscall.Syscall6(44, sock, uintptr(unsafe.Pointer(&httpInitMsg[0])), uintptr(len(httpInitMsg)), 0, 0, 0)

	// __NR_recvfrom, fd, buf, len(buf), flags, addr, addr_len
	var n, _, _ = syscall.Syscall6(45, sock, uintptr(unsafe.Pointer(&response[0])), uintptr(len(response)), 0, 0, 0)

	// __NR_sendto
	syscall.Syscall6(44, sock, uintptr(unsafe.Pointer(&packet[0])), uintptr(len(packet)), 0, 0, 0)

	// __NR_recvfrom
	syscall.Syscall6(45, sock, uintptr(unsafe.Pointer(&response[n])), uintptr(len(response))-n, 0, 0, 0)

	// __NR_close
	syscall.Syscall(3, sock, 0, 0)

	// __NR_write, STDOUT_FILENO
	syscall.Syscall(1, 1, uintptr(unsafe.Pointer(&response[0])), uintptr(len(response)))
}
