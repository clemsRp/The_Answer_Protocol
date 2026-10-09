# Server

The **Server** component manages TCP connection lifecycles, network I/O, client rate-limiting, and structured logging for *The Answer Protocol* (TAP). It acts as a lightweight network gateway between connected TCP clients and the underlying core game Engine.

---

## Architectural Overview

The server operates concurrently using Go's lightweight goroutines:
- **Listener Thread**: Listens on incoming TCP sockets and accepts client connections up to a configured max slot limit (`playerSlots`).
- **Client Handler Goroutines (`handleClient`)**: Each active client connection receives dedicated reader and writer goroutines to handle bidirectional TCP stream buffering and non-blocking socket I/O.
- **Rate Limiter (`RateLimiter`)**: Implements a Token Bucket algorithm per client connection to protect the server from spam and denial-of-service attempts.
- **Log Streamer**: Collects structured JSON log entries asynchronously over a dedicated log channel and emits them to configured outputs.

---

## Connection Flow & Exchanger Channel

The server does **not** process game logic directly. Instead, it delegates all business logic to the `Engine` via a thread-safe `pr.Exchanger`:

1. **Client Connects**: Connection accepted -> Rate limiter initialized -> Client ID assigned -> Added to `s.entering`.
2. **Input Stream**: Client sends a line-delimited message -> Rate limiter verifies token balance -> Passed to `Exchanger.ServerInput`.
3. **Engine Response**: Engine processes command -> Returns response via `Exchanger.ServerOutput` -> Server writes JSON payload back to socket.
4. **Client Disconnects**: Socket closed / timeout triggered -> Client ID sent to `Exchanger.LeaveChan` -> `playerSlots` released.

---

## Security & Rate Limiting

- **Token Bucket Algorithm**: Controls message throughput per client socket. Exempt commands (such as real-time position sync requests `NOTIFY_PLAYER_POSITION`) are bypassed to ensure smooth graphical movement.
- **Spam Violations**: Sockets exceeding rate limits receive `WARN_SPAM` (`ERR 902`). Continuous violations trigger automatic socket teardown (`ERR 902 CONNECTION_CLOSED_DUE_TO_SPAM`).
- **Zombie Socket Pruning**: Idle sockets without valid authentication or heartbeat activity within `IdleTimeout` are automatically terminated.

---

## Log Format

Server logs follow structured JSON formatting using Go's `log/slog`:
```json
{
  "timestamp": "2026-09-27T16:00:00Z",
  "level": "INFO",
  "message": "CLIENT_CONNECTED",
  "datas": {
    "client_id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
    "ip": "127.0.0.1:54321"
  }
}
```
