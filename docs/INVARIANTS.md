# Protocol Invariants

This system is defined by its invariants.

They are not optional — they define correctness.

---

## Core invariants

- Session identity survives transport death
- Authority is epoch-based
- Authority must be monotonic (no rollback)
- Only one authority per epoch (no split-brain)
- Stale transports must not revive
- Migration must not degrade correctness
- Continuity > optimal routing

---

## Why invariants matter

Without invariants:

- behavior becomes unpredictable
- debugging becomes impossible
- system correctness cannot be verified

With invariants:

- decisions are explainable
- behavior is testable
- failures are bounded

---

## Interpretation

This system prioritizes:

- continuity
- stability
- correctness

over:

- immediate reaction
- short-term optimality

---

## Example violation

If a stale path regains authority:

→ invariant broken  
→ session corruption risk  

Such cases must be explicitly prevented.

---

## Goal

Move from:

```
best-effort networking
```

to:

```
invariant-driven networking
```