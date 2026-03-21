# Hermesh

> Decentralized service mesh coordination over Hedera Consensus Service — no central server, no single point of failure, tamper-proof history forever.

[

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)

](https://golang.org)
[

![Hedera](https://img.shields.io/badge/Hedera-Testnet-8259EF?style=flat)

](https://hedera.com)
[

![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)

](LICENSE)
Node     : nairobi-node-1          Node     : nairobi-node-2
ID       : e4047862                ID       : 659f2432
Topic    : 0.0.8148188             Topic    : 0.0.8148188
Network  : Hedera Testnet          Network  : Hedera Testnet
[15:42:48] Announcing...           [15:43:01] Announcing...
🟢 Node joined: nairobi-node-2      🟢 Node joined: nairobi-node-1
— Zero configuration               — No bootstrap node
— No direct connection             — Pure HCS coordination
---

## The Problem

Every distributed system needs a control plane — something that answers:

- **Who is on the network right now?**
- **When did they join or leave?**
- **Can I trust that history was never tampered with?**

Traditional solutions answer the first two but fail the third:

| | etcd | Consul | Zookeeper | **Hermesh** |
|---|---|---|---|---|
| Single point of failure | Yes | Yes | Yes | **No** |
| Requires infra to operate | Yes | Yes | Yes | **No** |
| Audit trail immutable | No | No | No | **Yes** |
| Cryptographic proof of events | No | No | No | **Yes** |
| Cost | Servers + ops | Servers + ops | Servers + ops | **Fractions of a cent** |

A compromised etcd node can rewrite history. A failed Consul cluster takes your entire control plane down. Hermesh has no cluster to fail — the control plane lives on Hedera, a public network nobody owns.

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
|               |               |
Builds peer map  Builds peer map  Builds peer map
1. Every node **announces itself** to a shared HCS topic on startup
2. Every node **replays full history** from message #1 — new nodes get complete network state instantly, no bootstrap node needed
3. **Heartbeats** flow every 30 seconds — silence triggers suspected → dead transitions
4. **Graceful shutdown** publishes a leave message before dying
5. Every announcement is **signed with an ED25519 keypair** generated fresh on startup
6. Receiving nodes **verify the signature** before trusting any message — tampered or unsigned messages are silently dropped
7. Every verified event is **cryptographically timestamped on Hedera** — permanently, publicly, immutably

---

## Features

- **Zero bootstrap** — no seed node, no leader election needed to join
- **Full history replay** — any node joining sees the entire network history from genesis
- **ED25519 message signing** — every announcement cryptographically authenticated
- **Three-state health model** — `alive` → `suspected` → `dead` with configurable timeouts
- **HTTP API** — query peer state programmatically from any service
- **Dead node pruning** — configurable cleanup of stale entries
- **Graceful shutdown** — leave messages published on SIGINT/SIGTERM
- **Cross-platform binaries** — Linux, macOS, Windows via single Makefile
- **Zero runtime dependencies** — single static binary

---

## HTTP API

Every Hermesh node exposes a local HTTP API, making it a drop-in sidecar for service discovery:

```bash
# Health check
curl http://localhost:9000/health
# {"status":"ok","peers":8,"alive":2}

# Full peer registry
curl http://localhost:9000/peers
# {"peers":[...],"total":8}

# Filter by status
curl http://localhost:9000/peers?status=alive

# Filter by service name
curl http://localhost:9000/peers?service=payments-api
Any service can query Hermesh like a sidecar — no SDK required, just HTTP.
Security
Every node generates a fresh ED25519 keypair on startup. Every message published to HCS is signed with that key. Receiving nodes verify the signature before updating their peer map.
No valid signature  →  message dropped silently
Tampered payload    →  signature mismatch → dropped
Rogue node          →  unsigned messages → ignored by all legitimate nodes
No central certificate authority. No PKI infrastructure. No shared secrets. Just math.
Architecture
Hermesh follows the Unix philosophy — each module does one thing.
config/    — loads and validates environment configuration
hedera/    — speaks to Hedera HCS, nothing else
identity/  — generates ED25519 keypairs, signs and verifies envelopes
node/      — defines what a Node IS (ID, IP, port, status, timestamps)
announce/  — publishes join / heartbeat / leave to HCS
discover/  — subscribes to HCS, verifies signatures, hands messages to handlers
peer/      — thread-safe peer map with GC and pruning
api/       — HTTP server exposing /peers and /health
cli/       — tabwriter peer table, event log, banner
main.go    — wires everything together
Modules are deliberately ignorant of each other. Swap Hedera for NATS tomorrow — only touch hedera/.
Quick Start
Prerequisites:
Go 1.21+
A free Hedera testnet account — portal.hedera.com
Install:
git clone https://github.com/Emmanuel326/hermesh.git
cd hermesh
go mod tidy
Configure:
cp .env.example .env
nano .env
HEDERA_ACCOUNT_ID=0.0.xxxxxxx
HEDERA_PRIVATE_KEY=your_hex_private_key_without_0x_prefix
HEDERA_NETWORK=testnet
HERMESH_TOPIC_ID=0.0.xxxxxxx
HERMESH_NODE_NAME=my-node-1
HERMESH_NODE_PORT=8080
Run:
go build -o bin/hermesh .
./bin/hermesh
With options:
# Custom API port, prune dead nodes after 5 minutes
./bin/hermesh --api-port=9000 --prune=5m
Cross-platform builds:
make build-all
Two-Node Demo
Terminal 1 — Node 1:
./bin/hermesh
Terminal 2 — Node 2 (different account, same topic):
env $(grep -v '^#' .env.node2 | xargs) ./bin/hermesh
Watch Node 1's peer table update live the moment Node 2 announces. Kill Node 2 — watch it go suspected, then dead. Every event is verifiable on HashScan in real time.
Verify on-chain:
https://hashscan.io/testnet/topic/0.0.8148188
Live Testnet
Hermesh has been running against a live Hedera testnet topic since March 2026. Every session is permanently recorded:
Topic: 0.0.8148188
Account: 0.0.8146350
Network: Hedera Testnet
Browse the full immutable event history on HashScan.
Roadmap
[x] Node join / heartbeat / leave lifecycle
[x] ED25519 message signing and verification
[x] Full history replay on node startup
[x] Three-state health model (alive / suspected / dead)
[x] HTTP API — /peers and /health
[x] Dead node pruning with configurable TTL
[x] Cross-platform builds with Green Tea GC
[ ] hermesh query — CLI subcommand for scripting
[ ] Multi-topic sharding — horizontal scaling across HCS topics
[ ] TUI dashboard — real-time network visualizer
[ ] Mainnet support
[ ] Trustless Kubernetes control plane — etcd replaced by HCS
Built With
Go — single binary, zero runtime dependencies
Hedera HCS — immutable message bus
Hiero SDK for Go — Hedera Go SDK
Green Tea GC — Go experimental GC for lower latency
License
MIT
Built in Nairobi, Kenya. From a phone. With a broken laptop. In one evening.
