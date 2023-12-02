package main

import (
	"runtime"
	"syscall"
	"unsafe"
)

// tinygo build -no-debug -scheduler=none -gc=custom

// void *mmap(void *addr, size_t length, int prot, int flags, int fd, off_t offset);
//
//export mmap
//func libc_mmap(addr unsafe.Pointer, length uintptr, prot, flags, fd int32, offset uintptr) unsafe.Pointer

// Syscall is a variadic function that can take upto 6 arguments other than the
// syscall number, so we declare it accordingly. The libc implementation anyways
// pops off all 6 arguments off the stack regardless of the syscall in question

// long syscall(long, long, long, long, long, long, long)
//
//export syscall
//func syscall(n int64, arg1 int64, arg2 int64, arg3 int64, arg4 int64, arg5 int64, arg6 int64) int64

//go:linkname initHeap runtime.initHeap
func initHeap() {
}

//go:linkname alloc runtime.alloc
func alloc(size uintptr, layoutPtr unsafe.Pointer) unsafe.Pointer {
	var ptr, _, err = syscall.Syscall6(9, 0x0, size, 3, 34, 0xFFFFFFFF, 0)

	if err != 0 {
		return unsafe.Pointer(nil)
	}

	return unsafe.Pointer(uintptr(ptr))
}

//go:linkname free runtime.free
func free(ptr unsafe.Pointer) {
}

//go:linkname markRoots runtime.markRoots
func markRoots(start, end uintptr) {
}

//go:linkname GC runtime.GC
func GC() {
}

//go:linkname SetFinalizer runtime.SetFinalizer
func SetFinalizer(obj interface{}, finalizer interface{}) {
}

//go:linkname ReadMemStats runtime.ReadMemStats
func ReadMemStats(ms *runtime.MemStats) {
}

func main() {
	var str = []byte{'h', 'i', '\n'}
	syscall.Syscall(1, 1, uintptr(unsafe.Pointer(&str[0])), 3)
}
