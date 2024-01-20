## Syscalls only, dummy GC, custom ldflags, custom entrypoint

`tinygo build -o main -no-debug -scheduler none -gc none -panic trap -target spec.json`

`strip --strip-section-headers -R .comment -R .note -R .eh_frame main`
