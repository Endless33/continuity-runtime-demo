# Flagship Continuity Demo

## Goal

Demonstrate:

```
session survives transport failure without reconnect
```

---

## Scenario

1. Session is established over WiFi
2. WiFi becomes unstable
3. Runtime detects degradation
4. Decision engine evaluates alternatives
5. New transport selected (e.g. 5G)
6. Authority transferred (epoch++)
7. Old transport becomes stale
8. Session continues

---

## Key properties

- same session identity
- no reconnect
- no reset
- no state rebuild

---

## What to observe

- state transition (ATTACHED → RECOVERING)
- decision logs (score + confidence)
- authority transfer
- stale path rejection

---

## Expected output

```
[EVENT] WiFi degraded
[DECISION] migrate=true
[AUTHORITY] epoch 2 granted to 5G
[CHECK] stale WiFi rejected
[RESULT] session continues
```

---

## Why this matters

This demo proves:

- continuity is preserved
- failure is handled as runtime event
- system does not rely on reconnect semantics

---

## Interpretation

This is not failover.

This is **continuity without session break**.