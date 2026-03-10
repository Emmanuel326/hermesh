# Hermesh

> A simple mesh where every node reads and writes to a shared, immutable, append-only nervous system — independently, with no central authority.

Hermesh is a lightweight service mesh coordinator written in Go. Instead of relying on a central server like etcd, Consul, or Zookeeper, Hermesh uses [Hedera Consensus Service (HCS)](https://hedera.com/consensus-service) as its single source of truth — a public, tamper-proof, append-only log that every node in the network reads and writes to directly.

Think of it as Kubernetes service discovery — but the control plane lives on a blockchain nobody owns.

---

## The Problem

Every distributed system needs to answer three questions:

- Who is on the network right now?
- When did they join or leave?
- Can I trust that history was never altered?

Traditional solutions (etcd, Consul, Zookeeper) answer the first two but fail the third — their audit trails live in mutable databases. A compromised node can rewrite history. A failed etcd cluster can take your entire control plane down.

Hermesh answers all three. Permanently.

---

## How It Works
Node A starts          Node B starts          Node C starts
|                      |                      |
└──── announce ────────┴──── announce ─────────┘
|
Hedera HCS Topic
(immutable message bus)
|
┌────────────────┼────────────────┐
|                |                |
Node A reads     Node B reads     Node C reads
full history     full history     full history
|                |                |
Builds peer map  Builds peer map  Builds peer map
1. Every node announces itself to a shared HCS topic on startup
2. Every node subscribes to that topic and replays full history from message #1
3. New nodes get complete network state immediately — no bootstrap node needed
4. Heartbeats flow every 30 seconds — silence means suspected, then dead
5. Graceful shutdown sends a leave message before dying
6. Every event is cryptographically timestamped on Hedera — forever

---

## Why Hedera HCS

| | Traditional (etcd) | Hermesh (HCS) |
|---|---|---|
| Single point of failure | Yes — etcd cluster | No |
| Audit trail mutable | Yes | No — append only |
| Requires infra to run | Yes | No |
| Cost | Servers + ops | Fractions of a cent per message |
| Cryptographic proof | No | Yes |

---

## Architecture

Hermesh is built on the Unix philosophy — each module does one thing and does it well. Modules are deliberately dumb about each other.
config/       — loads and validates configuration
hedera/       — speaks to Hedera, nothing else
node/         — defines what a node IS
announce/     — publishes join, heartbeat, leave to HCS
discover/     — subscribes to HCS, hands messages to handlers
peer/         — thread-safe peer map with garbage collection
cli/          — presents network state to the operator
main.go       — wires everything together
No module knows more than it needs to. Swap Hedera for NATS tomorrow — only touch `hedera/`.

---

## Quick Start

**Prerequisites:**
- Go 1.25+
- A free Hedera testnet account — [portal.hedera.com](https://portal.hedera.com)

**Install:**
```bash
git clone https://github.com/Emmanuel326/hermesh.git
cd hermesh
go mod tidy
Configure:
cp .env.example .env
nano .env
HEDERA_ACCOUNT_ID=0.0.xxxxxxx
HEDERA_PRIVATE_KEY=your_private_key_here
HEDERA_NETWORK=testnet
HERMESH_TOPIC_ID=0.0.xxxxxxx
HERMESH_NODE_NAME=my-node-1
HERMESH_NODE_PORT=8080
Run:
go run ./...
Build (cross-platform with Green Tea GC):
GOEXPERIMENT=greenteagc GOOS=linux GOARCH=amd64 go build -o hermesh-linux-amd64 ./...
GOEXPERIMENT=greenteagc GOOS=darwin GOARCH=arm64 go build -o hermesh-darwin-arm64 ./...
GOEXPERIMENT=greenteagc GOOS=windows GOARCH=amd64 go build -o hermesh-windows-amd64.exe ./...
Demo
██╗  ██╗███████╗██████╗ ███╗   ███╗███████╗███████╗██╗  ██╗
...

  Node     : nairobi-node-1
  ID       : e4047862
  Address  : 10.197.184.229:8080
  Topic    : 0.0.8148188
  Network  : Hedera Testnet

  [15:42:43] Subscribing to network topic...
  [15:42:43] Listening for peers...
  [15:42:45] Announcing node to network...
  [15:42:48] Node announced: [e4047862] nairobi-node-1 @ 10.197.184.229:8080 (alive)
🟢 Node joined:     [e4047862] nairobi-node-1 @ 10.197.184.229:8080 (alive)

  NAME             ID         ADDRESS               STATUS     LAST SEEN
  ----             --         -------               ------     ---------
  nairobi-node-1   e4047862   10.197.184.229:8080   🟢 alive    12s ago

  Total: 1 peers (1 alive)
Verify every event live on HashScan — public, immutable, tamper-proof.
Roadmap
[ ] Node identity signing — cryptographic proof of node authenticity
[ ] Multi-topic sharding — horizontal scaling across HCS topics
[ ] REST API — query peer state programmatically
[ ] TUI dashboard — real-time network visualizer
[ ] Mainnet support
[ ] Trustless Kubernetes control plane — etcd replaced by HCS
Built With
Go — single binary, zero dependencies at runtime
Hedera HCS — immutable message bus
Hiero SDK for Go — Hedera Go SDK
Green Tea GC — Go 1.25+ experimental garbage collector for lower latency
License
MIT
Built in Nairobi, Kenya. From a phone. With a broken laptop. In one evening.
