# wscodegolf

| Attempt | Size (in bytes) | Compiler |
| ------- | --------------- | -------- |
| [Using `x/net/websocket`](https://github.com/git-bruh/wscodegolf/tree/main/client/net_websocket) | 4128768 (stripped) | Go |
| [Pure standard library](https://github.com/git-bruh/wscodegolf/tree/main/client/stdlib_only) | 727684 (1814528 without UPX) | Go |
| [Syscalls only](https://github.com/git-bruh/wscodegolf/tree/main/client/syscalls_only) | 360692 (847872 without UPX) | Go |
| [Syscalls only](https://github.com/git-bruh/wscodegolf/tree/main/client/syscalls_only) | 13056 | TinyGo |
| [Syscalls with dummy GC](https://github.com/git-bruh/wscodegolf/tree/main/client/syscalls_no_gc) | 12720 | TinyGo |
| [Syscalls with dummy GC, custom ldflags](https://github.com/git-bruh/wscodegolf/tree/main/client/syscalls_no_gc_ldflags) | 6600 | TinyGo |
| [Syscalls with no GC, custom ldflags, custom entrypoint](https://github.com/git-bruh/wscodegolf/tree/main/client/syscalls_no_gc_custom_entry) | 810 | TinyGo |
| [Syscalls with no GC, custom ldflags, custom entrypoint, 32-bit](https://github.com/git-bruh/wscodegolf/tree/main/client/syscalls_no_gc_custom_entry_32bit) | 538 | TinyGo |
