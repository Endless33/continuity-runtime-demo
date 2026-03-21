Continuity Runtime Demo

«failure ≠ connection death
failure = runtime event
continuity is enforced, not recovered»

---

⚡ TL;DR

session survives transport death
no reconnect
no reset

---

Core claim

Continuity should be enforced by the runtime,
not reconstructed after failure.

---

Quick start (1 command)

Run the core demo:

go run ./cmd/migration_demo/main.go

Expected:

- WiFi fails
- system detects degradation
- migration happens
- session continues (no reset)

---

Mental model

session != connection

session = identity
transport = attachment
failure = runtime event

This system does not "reconnect".

It rebinds the session to a new transport.

---

What this is

An experimental Go prototype exploring session continuity under transport volatility.

Instead of binding identity to a connection, this project models:

- session as the primary object
- transport as replaceable
- continuity as invariant

---

Architecture

"Continuity Runtime" (docs/architecture.png)

---

Core idea

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

Core invariants

These properties define the system behavior:

- Session identity survives transport death
- Only one authority per epoch
- Authority is monotonic (no rollback)
- Stale transports must not revive
- Continuity > optimality

This is not an implementation detail — this is the contract.

---

Runtime model (decision layer)

The system does NOT react to raw signals directly.

Instead it operates through filtered, time-aware decisions.

---

Signal processing

Raw network signals are noisy.

We convert them into stable signals:

- RTT → EWMA (fast + slow)
- Packet loss → rolling window
- Jitter → smoothed deviation

Goal:

react to trends, not spikes

---

State model

HEALTHY → DEGRADED → FAILED

Transitions are NOT instantaneous.

They require:

- K consecutive bad samples
- time window validation
- confidence threshold

---

Hysteresis (anti-flapping)

We intentionally introduce asymmetry:

- degrade → fast
- recover → slow

Example:

enter DEGRADED: 3 bad samples over 3s
enter FAILED: N missed heartbeats
recover: 10–30s stable window

Goal:

avoid oscillation under unstable conditions

---

Decision engine

Instead of binary logic:

score(path) + confidence → decision

Where:

- score = latency + loss + stability
- confidence = signal consistency over time

Migration condition:

new_path_score - current_score > margin
AND confidence is high

This makes decisions:

- explainable
- testable
- reproducible

---

What is implemented

Runtime

- state machine (ATTACHED → RECOVERING)
- decision engine (score / confidence)
- migration trigger
- authority handoff (epoch model)
- stale transport rejection
- hysteresis + time-window gating
- EWMA-based signal smoothing

---

Protocol

- wire packet format
- versioning
- replay protection
- sequence window validation
- session init / init ack
- authority transfer
- keepalive
- close

---

Reliability

- ACK flow
- retransmission policy
- timeout policy

---

Simulation

- multiple transports (wifi / 5g / lte)
- latency + jitter
- packet loss
- packet duplication
- lossy exchange
- two-node interaction

---

Observability

Designed for decision explainability:

- structured trace
- timeline replay
- invariant checks
- decision logs

Example:

[EVENT] WiFi degraded
[SIGNAL] rtt_ewma=182ms loss=0.12
[DECISION] migrate=true (margin=87.8, confidence=0.94)
[AUTHORITY] epoch 2 granted to 5G
[CHECK] stale WiFi rejected

---

🚀 Demos

Handshake

go run ./cmd/handshake_demo/main.go

---

Migration (recommended)

go run ./cmd/migration_demo/main.go

---

Two-node

go run ./cmd/two_node_demo/main.go

---

Example output

[EVENT] WiFi failed
[DECISION] migrate=true (margin=87.8, confidence=1.00)
[AUTHORITY] epoch 2 granted to 5G
[CHECK] stale WiFi rejected
[RESULT] session continues

---

Key property

NO reconnect
NO session reset
CONTINUITY PRESERVED

---

Why this matters

This is not about building "another VPN".

The question is:

Can session continuity be preserved under failure without reconnect?

If yes → this leads to:

- zero-reset mobile handoff
- transport-independent sessions
- next-gen VPN / overlay models
- runtime-driven networking

---

Why this is hard

Real networks are:

- noisy
- unstable
- inconsistent

Naive systems:

react too fast → flapping
react too slow → long recovery

This project explores:

controlled, explainable decision-making under uncertainty

---

What makes this different

Most systems:

- recover after failure

This model:

- avoids breaking the session in the first place

---

Relation to existing systems

- QUIC (connection migration)
- MPTCP (multipath)
- WireGuard (roaming)

«continuity is a first-class invariant, not a side-effect»

---

Design stance

Most systems optimize for:

- throughput
- latency

This system prioritizes:

- continuity
- stability
- explainability

---

Status

Early prototype.

---

Next steps

- real UDP / QUIC transport
- retransmission improvements
- packet scheduling
- adaptive hysteresis tuning
- formal protocol spec

---

Direction

demo → runtime → protocol → architecture

---

Feedback

Looking for:

- protocol flaws
- edge cases
- invariant violations
- migration race conditions

---

📦 Project structure

LICENSE
SECURITY.md
CONTRIBUTING.md
ROADMAP.md
SUPPORT.md
docs/

---

🧬 Protocol invariants

docs/INVARIANTS.md

---

🔬 Flagship demo

docs/FLAGSHIP_DEMO.md

---

❤️ Support

SUPPORT.md