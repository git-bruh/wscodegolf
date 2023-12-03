package main

import (
	"syscall"
	"unsafe"
)

// tinygo build -no-debug -scheduler=none -gc=none -panic=trap

var buffer [1024]byte
var used uintptr = 0

//go:linkname alloc runtime.alloc
func alloc(size uintptr, layoutPtr unsafe.Pointer) unsafe.Pointer {
	// TODO align
	var ptr = unsafe.Pointer(&buffer[used])
	used += size

	return ptr
}

func main() {
	var str = []byte{'h', 'i', '\n'}
	syscall.Syscall(1, 1, uintptr(unsafe.Pointer(&str[0])), 3)
}
