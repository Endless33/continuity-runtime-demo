# Continuity Runtime Demo

This is a minimal demo of a runtime model where:

failure ≠ connection death  
failure = runtime event

## Idea

Traditional systems:
- connection drops
- reconnect required
- session resets

This approach:
- treats failure as a state transition
- evaluates alternative paths
- transfers authority (epoch-based)
- rejects stale transports
- continues execution

## Demo

Run:

go run ./cmd/demo/main.go

Output:

[EVENT] WiFi failed  
[DECISION] migrate=true  
[AUTHORITY] epoch 2 granted  
[CHECK] stale rejected  
[RESULT] session continues  

## What this shows

This is NOT a retry system.

There is:
- no reconnect
- no session reset
- no interruption

The session continues across transport change.

## Status

Early prototype.  
Focused on modeling runtime behavior before real networking layer.

## Next

- real transport abstraction
- real migration instead of mock
- packet-level continuity

---

Open to feedback.