# Continuity Runtime Demo

> failure ≠ connection death  
> failure = runtime event

---

## Why this exists

Most networking systems treat failure like this:

- path breaks  
- connection dies  
- reconnect  
- session rebuilt  

This leads to:

- packet loss  
- broken ordering  
- state reset  
- visible interruption  

---

## A different model

Instead of:

failure → reconnect → restore

This demo models:

failure → runtime transition → continuity preserved

---

## Core idea

Not connection-centric:

connection → identity

But session-centric:

session → continuity → changing transports

Where:

- session identity is stable  
- transports are replaceable  
- failure is expected  
- continuity is enforced  

---

## What this demo models

This is not just a print demo.

### Transport behavior

- multiple transports (WiFi / 5G / LTE)  
- latency + jitter simulation  
- packet loss  

### Runtime decisions

- scoring-based path selection  
- adaptive multipath (enabled only under degradation)  
- latency racing (fastest path wins)  

### Stream continuity

- packet duplication during migration  
- deduplication  
- reorder buffer  
- partial ordering  

### Reliability layer

- session-level ACK tracking  
- retransmission queue  
- frame assembly  
- forward error correction (FEC simulation)  

### Session model

- explicit session identity  
- epoch-based authority  
- stale transport rejection  

### Observability

- structured trace (JSON events)  
- timeline reconstruction  
- replay system  
- invariant checks (e.g. epoch monotonicity)  

---

## Example flow

WiFi fails  
↓  
runtime detects degradation  
↓  
adaptive overlap starts (WiFi + 5G)  
↓  
latency racing selects fastest path  
↓  
authority transferred (epoch++)  
↓  
stale path rejected  
↓  
packet loss partially recovered (FEC / retransmit)  
↓  
stream continues  

---

## Example output

[EVENT] WiFi failed  
[DECISION] migrate=true (margin=87.8, confidence=1.00)  
[STATE] ATTACHED → RECOVERING  
[MULTIPATH] starting overlap  
[RACE] winner=5g  
[AUTHORITY] epoch 2 granted  
[CHECK] stale wifi rejected  
[FEC] recovered missing packet  
[FRAME] frame assembled  
[STATE] RECOVERING → ATTACHED  
[RESULT] session continues  

---

## What this is NOT

This is not:

- a VPN implementation  
- a retry mechanism  
- a reconnect strategy  

---

## What this could become

If extended beyond simulation:

- zero-loss handoff between networks  
- seamless WiFi ↔ 5G switching  
- transport-independent sessions  
- runtime-level continuity guarantees  
- new class of session-centric protocols  

---

## Research direction

The question:

Can session continuity be enforced as a system invariant even when transports fail?

If yes:

failure becomes just another runtime event  
not a terminal condition  

---

## Status

Early prototype.  
Focused on runtime behavior and continuity modeling.

---

## Next steps

- real transport integration (UDP / QUIC)  
- congestion control  
- real FEC  
- protocol-level design  

---

## Feedback

Looking for:

- critical feedback  
- edge cases  
- architectural concerns  

---

This is an exploration, not a finished system.