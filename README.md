# continuity-runtime-demo

Minimal Go demo showing a different approach to network failure handling.

Instead of reconnecting after failure, this runtime treats failure as a transition inside an active session.

## Behavior

- detect failure
- evaluate alternative path
- transfer authority (epoch-based)
- reject stale transport
- continue session without reconnect

## Run

```bash
go run ./cmd/demo/main.go
