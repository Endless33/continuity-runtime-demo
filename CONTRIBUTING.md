# Contributing

This project explores a **continuity-first networking model**.

It is not just code — it is a protocol + runtime + invariant system.

## Core principles

- Continuity > convenience
- Invariants > shortcuts
- Explainability > hidden behavior
- Determinism > heuristics

## What is valuable

We are especially interested in:

- protocol-level feedback
- invariant violations
- race conditions
- migration edge cases
- instability / flapping scenarios
- replay or authority bugs
- decision engine flaws

## Before contributing

Ask yourself:

- does this preserve invariants?
- does this improve explainability?
- does this introduce hidden state?

## Contribution style

- keep changes minimal and focused
- prefer explicit state over implicit behavior
- document WHY, not just WHAT
- avoid “magic fixes”

## Areas to contribute

- runtime decision engine (EWMA, hysteresis, scoring)
- protocol state machine
- transport abstraction
- observability and trace clarity
- simulation scenarios
- documentation

## Communication

- open issues
- propose changes
- challenge assumptions

Strong critique is welcome.