# Roadmap

This project evolves in layers:

```
demo → runtime → protocol → architecture
```

---

## Phase 1 — Continuity Demonstration

Goal:

- prove that session survives transport failure

Includes:

- migration demo
- authority transfer
- basic invariants

---

## Phase 2 — Runtime Stabilization

Focus:

- EWMA signal smoothing
- hysteresis (anti-flapping)
- time-window validation
- confidence-based decisions

Goal:

- stable behavior under noisy conditions

---

## Phase 3 — Transport Layer

Focus:

- real UDP improvements
- transport abstraction
- QUIC integration (future)

Goal:

- realistic networking conditions

---

## Phase 4 — Protocol Formalization

Focus:

- explicit state machine
- packet-level specification
- invariant documentation
- failure scenarios

Goal:

- move toward protocol definition

---

## Phase 5 — System Evolution

Focus:

- transport-independent sessions
- continuity-first networking model
- distributed authority model

Goal:

- architecture-level system

---

## Long-term direction

- zero-reset handoff
- mobility-first networking
- resilient session overlays
- runtime-driven networking stacks