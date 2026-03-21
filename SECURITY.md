# Security Policy

## Status

This project is a **research prototype** and NOT production-ready.

It is designed for experimentation, protocol exploration, and runtime behavior validation.

## Scope

- The protocol is evolving
- Runtime logic is under active development
- Cryptographic components are NOT finalized
- Security guarantees are NOT formally verified

## Known limitations

- No formal security audit
- No hardened transport layer
- No production-grade key management
- No resistance guarantees against advanced adversaries

## Reporting vulnerabilities

If you discover:

- protocol flaws
- security vulnerabilities
- replay / authority issues
- invariant violations

Please:

- open a GitHub issue (preferred)
- or contact privately if the issue is sensitive

## Responsible usage

Do NOT use this project in:

- production environments
- critical infrastructure
- privacy-sensitive systems

This is a research system.

## Philosophy

Security in this project is tied to **protocol invariants**:

- authority must be explicit
- stale paths must be rejected
- session continuity must not break consistency

These properties are still being validated.