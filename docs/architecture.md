# Architecture

## Overview

This system treats the **session as the primary object**, while transports are replaceable.

---

## Model

```
        +----------------------+
        |   Session Identity   |
        |   (session_id)       |
        +----------+-----------+
                   |
                   v
        +----------------------+
        |   Continuity Runtime |
        |----------------------|
        | state machine        |
        | decision engine      |
        | authority manager    |
        | trace recorder       |
        +----------+-----------+
                   |
        ---------------------------
        |            |            |
        v            v            v
   +--------+   +--------+   +--------+
   | WiFi   |   |  5G    |   |  LTE   |
   +--------+   +--------+   +--------+

```

---

## Key Idea

```
session
   ↓
runtime (controls continuity)
   ↓
transport (replaceable)
```

NOT:

```
session → tunnel → transport
```

BUT:

```
session → runtime → dynamic transports
```

---

## Components

### Session Identity
- stable across failures
- not tied to transport

---

### Runtime
- owns session lifecycle
- handles failure events
- executes migration
- enforces invariants

---

### Decision Engine
- compares paths
- computes score + confidence
- decides migration

---

### Authority (Epoch)
- ensures single owner
- prevents split-brain
- rejects stale paths

---

### Transport Layer
- WiFi / 5G / LTE (simulated)
- interchangeable
- not trusted for continuity

---

### Trace System
- records every event
- enables replay
- enables debugging

---

## Invariant

```
session continuity must survive transport failure
```

---

## What is different

- no reconnect
- no session reset
- failure is internal state transition
- continuity is enforced, not recovered