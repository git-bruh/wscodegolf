## Syscalls only

`go build -trimpath -ldflags '-s -w' -o main`

`tinygo build -o main -no-debug -scheduler none -gc leaking -panic trap`
