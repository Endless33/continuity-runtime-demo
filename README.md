# Continuity Runtime Demo

> failure ≠ connection death  
> failure = runtime event

---

## What this is

An experimental Go prototype exploring **session continuity under transport volatility**.

Instead of binding identity to a connection, this project models:

- session as the primary object
- transport as replaceable
- continuity as invariant

---

## Architecture

![Continuity Runtime](docs/architecture.png)

---

## Core idea

Traditional systems:

- connection drops  
- reconnect required  
- session reset  

This approach:

- failure is a runtime event  
- system evaluates alternatives  
- authority is transferred (epoch-based)  
- stale path is rejected  
- session continues  

---

## What is implemented

### Runtime
- state machine (ATTACHED → RECOVERING)
- decision engine (score / confidence)
- migration trigger
- authority handoff (epoch model)
- stale transport rejection

### Protocol
- wire packet format
- versioning
- replay protection
- sequence window validation
- session init / init ack
- authority transfer
- keepalive
- close

### Reliability
- ACK flow
- retransmission policy
- timeout policy

### Simulation
- multiple transports (wifi / 5g / lte)
- latency + jitter
- packet loss
- packet duplication
- lossy exchange
- two-node interaction

### Observability
- structured trace
- timeline replay
- invariant checks

---

## Demos

### Handshake
```
go run ./cmd/handshake_demo/main.go
```

Shows:
- session init
- init ack
- keepalive
- close

---

### Migration
```
go run ./cmd/migration_demo/main.go
```

Shows:
- data before failure
- WiFi failure event
- migration decision
- authority transfer
- data continues (no reset)

---

### Two-node
```
go run ./cmd/two_node_demo/main.go
```

Shows:
- two nodes exchanging packets
- ACK flow
- lossy network behavior
- migration + invariants

---

## Example output

```
[EVENT] WiFi failed
[DECISION] migrate=true (margin=87.8, confidence=1.00)
[AUTHORITY] epoch 2 granted to 5G
[CHECK] stale WiFi rejected
[RESULT] session continues
```

---

## Key property

```
NO reconnect
NO session reset
CONTINUITY PRESERVED
```

---

## Why this matters

This is not about building "another VPN".

The question is:

**Can session continuity be preserved under failure without reconnect?**

If yes → this leads to:

- zero-reset network handoff
- transport-independent sessions
- next-gen VPN / overlay models
- runtime-driven networking

---

## What makes this different

Most systems:
- recover after failure

This model:
- avoids breaking the session in the first place

It is closer to:
- session migration
- authority transfer
- runtime-controlled networking

Not:
- retry logic
- reconnect loops

---

## Status

Early prototype.

What it is:
- protocol + runtime model
- research-grade implementation
- traceable + testable

What it is not:
- production VPN
- production crypto
- congestion-controlled stack

---

## Next steps

- real UDP / QUIC transport
- retransmission improvements
- packet scheduling
- formal protocol spec
- protocol diagrams (flow, packet-level)

---

## Direction

This repo is evolving:

```
demo → runtime → protocol → architecture
```

---

## Feedback

Looking for:

- protocol flaws
- edge cases
- invariant violations
- migration race conditions
- replay / stale-path issues