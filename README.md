# Hermesh

Trustless control plane for distributed systems — powered by Hedera Consensus Service.

No clusters. No bootstrap nodes. No single point of failure.
Just an immutable, cryptographically verifiable coordination layer.

---

## Why This Matters

Modern distributed systems assume one thing:

> *"You must trust your control plane."*

Hermesh removes that assumption.

Instead of relying on a cluster you operate and trust, Hermesh uses a public, immutable log as the coordination layer — making every event verifiable, tamper-proof, and globally consistent.

---

## The Problem

Every distributed system needs a control plane — something that answers:

- Who is on the network right now?
- When did they join or leave?
- Can I trust that history was never altered?

Traditional systems solve the first two — but rely on infrastructure you must operate and trust.

| | etcd | Consul | Zookeeper | **Hermesh** |
|---|---|---|---|---|
| Requires cluster management | Yes | Yes | Yes | **No** |
| Operational complexity | High | High | High | **Minimal** |
| Immutable audit trail | No | No | No | **Yes** |
| Cryptographic verification | No | No | No | **Yes** |
| Bootstrap required | Yes | Yes | Yes | **No** |

Hermesh has no cluster to fail — the control plane lives on a public consensus network.

---

## How It Works
Node A starts        Node B starts        Node C starts
|                    |                    |
└──── announce ───────┴──── announce ──────┘
|
┌─ Hedera HCS Topic ─┐
│  immutable log     │
│  append-only       │
│  cryptographically │
│  timestamped       │
└────────────────────┘
|
┌───────────────┼───────────────┐
|               |               |
Node A reads    Node B reads    Node C reads
full history    full history    full history
**Lifecycle**

1. Node announces itself on startup
2. Nodes replay full history from genesis
3. Heartbeats maintain liveness
4. Silence → "suspected" → "dead"
5. Graceful shutdown publishes leave event

**Security Model**

- Every node generates an ED25519 keypair
- Every message is signed
- Every message is verified before trust
- All events are immutably timestamped

No PKI. No shared secrets. No central authority.

---

## Features

- **Zero bootstrap** — no seed nodes required
- **Full history replay** — instant state reconstruction
- **Cryptographic message signing** (ED25519)
- **Three-state health model** — `alive` → `suspected` → `dead`
- **HTTP API** for service discovery
- **Dead node pruning** (configurable TTL)
- **Graceful shutdown** handling
- **Cross-platform static binaries**
- **Zero runtime dependencies**

---

## When NOT to Use Hermesh

Hermesh is not a replacement for everything.

**Do NOT use it for:**
- High-frequency service discovery
- Low-latency coordination (<100ms requirements)
- Private systems with strict data locality constraints

**Hermesh is designed for:**
- Trustless or multi-party environments
- Audit-critical systems
- Coordination where history integrity matters more than latency

---

## HTTP API

Hermesh runs as a sidecar exposing a simple HTTP interface:

```bash
# Health check
curl http://localhost:9000/health
# {"status":"ok","peers":8,"alive":2}

# Full peer registry
curl http://localhost:9000/peers

# Filter by status
curl http://localhost:9000/peers?status=alive

# Filter by service
curl http://localhost:9000/peers?service=payments-api
No SDK required. Just HTTP.
Architecture
config/    — environment configuration
hedera/    — HCS interaction layer
identity/  — key generation & signing
node/      — node model
announce/  — publish lifecycle events
discover/  — consume & verify events
peer/      — thread-safe peer registry
api/       — HTTP interface
cli/       — terminal UI
main.go    — composition root
Each module is isolated and replaceable.
Quick Start
Prerequisites:
Go 1.21+
Hedera testnet account — portal.hedera.com
Install:
git clone https://github.com/Emmanuel326/hermesh.git
cd hermesh
go mod tidy
Configure:
cp .env.example .env
nano .env
HEDERA_ACCOUNT_ID=0.0.xxxxxxx
HEDERA_PRIVATE_KEY=your_private_key
HEDERA_NETWORK=testnet
HERMESH_TOPIC_ID=0.0.xxxxxxx
HERMESH_NODE_NAME=my-node
HERMESH_NODE_PORT=8080
Run:
go build -o bin/hermesh .
./bin/hermesh
Demo
Start Node A
Start Node B → instantly appears in Node A's peer table
Kill Node B → transitions to suspected then dead
Open HashScan → every event immutably on-chain
Restart any node → full state rebuilt from log alone
"Even if all nodes disappear, the network state can be reconstructed from public history."
Live Testnet
Topic: 0.0.8148188
Network: Hedera Testnet
HashScan: hashscan.io/testnet/topic/0.0.8148188
All events are permanently recorded and publicly verifiable.
Roadmap
[x] Node lifecycle (join / heartbeat / leave)
[x] ED25519 signing & verification
[x] Full history replay
[x] Health state transitions
[x] HTTP API
[x] Dead node pruning
[x] Cross-platform builds
[ ] CLI query tooling
[ ] Multi-topic sharding
[ ] TUI dashboard
[ ] Mainnet support
[ ] Kubernetes control plane integration
Built With
Go
Hedera Consensus Service
Hiero SDK for Go
License
MIT
Built in Nairobi, Kenya — in one evening, on limited hardware.
