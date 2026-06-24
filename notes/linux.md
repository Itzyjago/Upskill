# Linux / Unix notes

## Processes
- Every process has a PID and a parent (PPID). `ps aux`, `top`/`htop` to inspect.
- Foreground vs background: `cmd &`, `jobs`, `fg`, `bg`. `nohup`/`disown` to
  survive a shell exit.
- Exit status in `$?`: 0 = success, non-zero = failure.

## Signals
- `SIGTERM` (15) — polite "please stop", catchable. Default for `kill`.
- `SIGKILL` (9) — uncatchable, immediate. Last resort, leaves no cleanup.
- `SIGINT` (2) — Ctrl-C. `SIGHUP` (1) — terminal closed, often "reload config".
- Trap them in scripts: `trap 'cleanup' EXIT INT TERM`.

## Permissions
- `rwx` for user/group/other → octal: `chmod 644` (rw-r--r--), `755` (rwxr-xr-x).
- `chown user:group file`. The `x` bit on a *directory* means "can enter it".
- setuid/setgid/sticky bit: sticky on `/tmp` (1777) stops users deleting each
  other's files.

## File descriptors & redirection
- 0 = stdin, 1 = stdout, 2 = stderr.
- `cmd > out.log 2>&1` — merge stderr into stdout (order matters!).
- Everything is a file: pipes, sockets, devices under `/dev`, process info under
  `/proc/<pid>/`.

## Handy
- `lsof -i :8080` — what's holding a port. `ss -tlnp` — listening sockets.
- `df -h` disk, `du -sh *` per-dir size, `free -h` memory.
