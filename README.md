# =========================
# FILE: LICENSE
# =========================
MIT License

Copyright (c) 2026 Endless33

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction...

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND...

# =========================
# FILE: SECURITY.md
# =========================
# Security Policy

## Status

This project is a **research prototype**.

It is NOT production-ready and should NOT be used as a security product.

## Scope

- Experimental protocol
- Experimental runtime
- Cryptography is NOT finalized
- Transport behavior is evolving

## Reporting Issues

If you discover a vulnerability or protocol flaw:

- Open an issue (preferred for transparency)
- Or contact directly if sensitive

## Notes

- Do not assume confidentiality guarantees
- Do not deploy in real-world critical environments
- Behavior may change without notice

---

# =========================
# FILE: CONTRIBUTING.md
# =========================
# Contributing

This project is a **protocol/runtime experiment**, not a typical application.

## Before contributing

Please understand:

- This is NOT a production VPN
- The focus is **continuity + runtime behavior**
- Design consistency > quick patches

## What is welcome

- protocol analysis
- edge case discovery
- invariant violations
- migration race conditions
- runtime stability improvements

## What is NOT the focus (yet)

- UI
- packaging
- production deployment
- optimizations without architectural reason

## PR guidelines

- keep changes minimal and focused
- explain WHY, not only WHAT
- avoid breaking invariants

---

# =========================
# FILE: ROADMAP.md
# =========================
# Roadmap

## Phase 1 — Continuity demo (current)
- session survives transport failure
- epoch-based authority
- basic migration

## Phase 2 — Runtime stabilization
- hysteresis model
- EWMA smoothing
- decision confidence

## Phase 3 — Transport layer
- real UDP improvements
- QUIC integration

## Phase 4 — Protocol formalization
- packet-level spec
- state machine diagrams
- invariants formalization

## Phase 5 — System-level vision
- transport-independent sessions
- continuity-first networking model

---

# =========================
# FILE: SUPPORT.md
# =========================
# Support the Project

If you find this work interesting or useful:

- GitHub: https://github.com/Endless33/jumping-vpn-preview
- Demo: https://github.com/Endless33/continuity-runtime-demo

Support (optional):

- Ko-fi
- PayPal
- Crypto

(links provided in project discussions / posts)

## What support enables

- deeper protocol research
- better demos
- documentation
- runtime improvements

---

# =========================
# FILE: docs/INVARIANTS.md
# =========================
# System Invariants

These define the behavior of the system.

## Core invariants

- Session identity survives transport death
- Only one authority per epoch
- Authority must move forward (monotonic)
- Stale transports must not revive
- No implicit rollback to worse path
- Continuity > optimality

## Interpretation

The system prioritizes:

- stability over instant reaction
- correctness over speed
- continuity over optimal routing

---

# =========================
# FILE: docs/FLAGSHIP_DEMO.md
# =========================
# Flagship Demo

## Goal

Demonstrate:

```
session survives transport failure without reconnect
```

## Steps

1. Start session over WiFi
2. Inject failure
3. Trigger migration
4. Transfer authority (epoch)
5. Reject stale path
6. Continue data flow

## Expected result

- no reconnect
- no reset
- same session identity

---

# =========================
# ADD TO README (append at bottom)
# =========================

---

## Project structure

```
LICENSE
SECURITY.md
CONTRIBUTING.md
ROADMAP.md
SUPPORT.md
docs/
```

---

## Protocol invariants

See:

```
docs/INVARIANTS.md
```

---

## Flagship demo

See:

```
docs/FLAGSHIP_DEMO.md
```

---

## Support

If you want to support the work:

See:

```
SUPPORT.md
```