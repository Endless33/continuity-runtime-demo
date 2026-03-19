# Protocol Flow

## WiFi Failure → Migration → Continuity

---

## Step-by-step

```
CLIENT (runtime)                NETWORK                SERVER
      |                             |                     |
      |---- data (WiFi) ----------->|                     |
      |                             |                     |
      |        [EVENT] WiFi fails   |                     |
      |                             |                     |
      |---- evaluate paths -------->|                     |
      |                             |                     |
      |<--- decision: migrate ------|                     |
      |                             |                     |
      |---- AUTH_TRANSFER --------->|----> SERVER         |
      |                             |                     |
      |<--- AUTH_GRANTED -----------|<---- SERVER         |
      |                             |                     |
      |---- switch to 5G ---------->|                     |
      |                             |                     |
      |---- data continues -------->|                     |
      |                             |                     |
      |---- stale WiFi rejected ---X                     |
      |                             |                     |
      |        session continues                          |
```

---

## Runtime View

```
ATTACHED
   ↓
[WiFi failure]
   ↓
RECOVERING
   ↓
[authority granted]
   ↓
ATTACHED (new path)
```

---

## Decision Model

```
score(candidate) > score(current)
→ migrate

confidence > threshold
→ safe to migrate
```

---

## Authority Model

```
epoch N   → WiFi owns session
epoch N+1 → 5G owns session

old path cannot send packets anymore
```

---

## Result

```
NO reconnect
NO reset
NO session loss

→ continuity preserved
```

---

## Key Insight

Failure is not:

```
connection death
```

Failure is:

```
runtime event → controlled transition
```