# continuity-runtime-demo

Minimal Go runtime experiment demonstrating session continuity under transport failure.

Instead of treating connection loss as a terminal event (disconnect → reconnect), this runtime models failure as a transition inside an active session.

---

## Problem

In most systems, when a transport fails (e.g. WiFi drops):

- the connection is considered dead
- a reconnect is required
- session state may reset
- in-flight data can be lost
- ordering guarantees can break

This forces recovery logic into every layer of the stack.

---

## Approach

This runtime treats failure differently:

failure ≠ disconnect  
failure = state transition  

Instead of rebuilding the session, it preserves it.

---

## What happens on failure

When WiFi fails:

1. runtime detects degradation / failure  
2. evaluates an alternative path (e.g. 5G)  
3. decision engine decides whether to migrate  
4. authority is transferred (epoch-based)  
5. old transport is rejected (stale)  
6. session continues on the new path  

---

## Core concepts

### Session continuity

The session identity does not change across transports.

### Transport migration

The underlying path can change without breaking execution.

### Authority (epoch-based)

Each transport operates under an epoch:

- new transport → higher epoch  
- old transport → automatically invalid  

### Stale rejection

Packets from old transports are rejected after migration.

---

## Example output

[EVENT] WiFi failed  
[DECISION] migrate=true (margin=87.8, confidence=1.00, reason=better_path)  
[AUTHORITY] epoch 2 granted to 5G  
[CHECK] stale WiFi rejected  
[RESULT] session continues  

---

## Comparison

Traditional model:

WiFi failed  
→ reconnect required  
→ session reset  
→ possible data loss  

Continuity runtime:

WiFi failed  
→ migration decision  
→ authority transfer  
→ stale path rejected  
→ session continues  

---

## Run

go run ./cmd/demo/main.go

---

## Repository structure

cmd/demo/  
→ minimal runnable example  

internal/runtime/  
→ decision engine  
→ authority (epoch model)  
→ runtime logic  
→ trace system  

---

## Notes

This is a minimal runtime prototype focused on continuity and migration logic.

It is NOT:

- a production VPN  
- a full networking stack  
- a cryptographic implementation  

---

## Why this exists

This project explores a different model:

Can a session survive transport failure without reconnecting?

---

## Status

Experimental / early stage  

---

## License

MIT